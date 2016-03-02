/**
 * Copyright (c) 2015 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/trustedanalytics/go-cf-lib/types"
	"net/http"
	"sync"
)

var _ = Describe("Cf bindings", func() {

	var (
		negativeResponder httpmock.Responder
		requestFail       httpmock.Responder
		sut               CfAPI
	)

	BeforeEach(func() {
		httpmock.Activate()

		negativeResponder = responderGenerator(400, nil)
		requestFail = responderFailGenerator(nil)
		sut = CfAPI{Client: http.DefaultClient}
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	Describe("bind service method", func() {
		res := types.CfServiceBindingCreateResponse{
			Meta: types.CfMeta{
				GUID: "guid",
			},
		}
		resp := responderGenerator(201, res)

		Context("without errors", func() {
			It("should send nil to channel", func() {
				httpmock.RegisterResponder("POST", "/v2/service_bindings", resp)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error)

				go sut.BindService("app_guid", "service_guid", errorCh, &wg)

				err := <-errorCh

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("with errors", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("POST", "/v2/service_bindings", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error)

				go sut.BindService("app_guid", "service_guid", errorCh, &wg)

				err := <-errorCh

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("unbind app services method", func() {
		resources := types.CfBindingsResources{
			TotalResults: 1,
			Resources: []types.CfBindingResource{
				types.CfBindingResource{
					Meta: types.CfMeta{
						GUID: "guid",
					},
					Entity: types.CfBinding{
						AppGUID: "app_guid",
					},
				},
			},
		}
		getAppBindingsResp := responderGenerator(200, resources)
		deleteBindingResp := responderGenerator(204, nil)

		Context("without errors", func() {
			It("should send nil to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/service_bindings", getAppBindingsResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/service_bindings/guid", deleteBindingResp)

				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error)

				go sut.UnbindAppServices("app_guid", errorCh, &wg)

				err := <-errorCh

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when getting app bindings fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/service_bindings", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error)

				go sut.UnbindAppServices("app_guid", errorCh, &wg)

				err := <-errorCh

				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when deletion fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/service_bindings", getAppBindingsResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/service_bindings/guid", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error)

				go sut.UnbindAppServices("app_guid", errorCh, &wg)

				err := <-errorCh

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
