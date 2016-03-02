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

var _ = Describe("Cf delete", func() {

	var (
		negativeResponder httpmock.Responder
		sut               CfAPI
	)

	BeforeEach(func() {
		httpmock.Activate()

		negativeResponder = responderGenerator(400, nil)
		sut = CfAPI{Client: http.DefaultClient}
	})

	AfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	Describe("delete service instance if unbound method", func() {

		serviceBindingsRes := types.CfBindingsResources{
			TotalResults: 0,
			Resources:    []types.CfBindingResource{},
		}
		serviceBindingsResp := responderGenerator(200, serviceBindingsRes)

		deleteResp := responderGenerator(204, nil)

		comp := types.Component{
			GUID:         "guid",
			Name:         "Name",
			DependencyOf: []string{},
		}

		Context("without errors", func() {
			It("should send nil to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/service_instances/guid/service_bindings", serviceBindingsResp)
				httpmock.RegisterResponder("DELETE", "/v2/service_instances/guid", deleteResp)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteServiceInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when service is bound to other application", func() {
			serviceBindingsRes := types.CfBindingsResources{
				TotalResults: 1,
				Resources: []types.CfBindingResource{
					types.CfBindingResource{
						Meta: types.CfMeta{
							GUID: "guid1",
						},
					},
				},
			}
			serviceBindingsResp := responderGenerator(200, serviceBindingsRes)

			It("should not delete but finish ok", func() {
				httpmock.RegisterResponder("GET", "/v2/service_instances/guid/service_bindings", serviceBindingsResp)

				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteServiceInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when deletion fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/service_instances/guid/service_bindings", serviceBindingsResp)
				httpmock.RegisterResponder("DELETE", "/v2/service_instances/guid", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteServiceInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when getting service bindings fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/service_instances/guid/service_bindings", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteServiceInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("delete user provided service instance if unbound method", func() {

		serviceBindingsRes := types.CfBindingsResources{
			TotalResults: 0,
			Resources:    []types.CfBindingResource{},
		}
		serviceBindingsResp := responderGenerator(200, serviceBindingsRes)

		deleteResp := responderGenerator(204, nil)

		comp := types.Component{
			GUID:         "guid",
			Name:         "Name",
			DependencyOf: []string{},
		}

		Context("without errors", func() {
			It("should send nil to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid/service_bindings", serviceBindingsResp)
				httpmock.RegisterResponder("DELETE", "/v2/user_provided_service_instances/guid", deleteResp)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteUPSInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when service is bound to other application", func() {
			serviceBindingsRes := types.CfBindingsResources{
				TotalResults: 1,
				Resources: []types.CfBindingResource{
					types.CfBindingResource{
						Meta: types.CfMeta{
							GUID: "guid1",
						},
					},
				},
			}
			serviceBindingsResp := responderGenerator(200, serviceBindingsRes)

			It("should not delete but finish ok", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid/service_bindings", serviceBindingsResp)

				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteUPSInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when deletion fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid/service_bindings", serviceBindingsResp)
				httpmock.RegisterResponder("DELETE", "/v2/user_provided_service_instances/guid", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteUPSInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when getting service bindings fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid/service_bindings", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteUPSInstIfUnbound(comp, errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("delete routes method", func() {
		appSummary := types.CfAppSummary{
			GUID: "app_guid",
			Routes: []types.CfAppSummaryRoute{
				types.CfAppSummaryRoute{
					GUID: "route1_guid",
				},
				types.CfAppSummaryRoute{
					GUID: "route2_guid",
				},
			},
		}
		getAppSummaryResp := responderGenerator(200, appSummary)
		unassociateRouteResp := responderGenerator(200, nil)
		deleteResp := responderGenerator(204, nil)

		Context("without errors", func() {
			It("should send nil to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/routes/route1_guid", unassociateRouteResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/routes/route2_guid", unassociateRouteResp)
				httpmock.RegisterResponder("DELETE", "/v2/routes/route1_guid", deleteResp)
				httpmock.RegisterResponder("DELETE", "/v2/routes/route2_guid", deleteResp)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteRoutes("app_guid", errorCh, &wg)

				err := <-errorCh
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when application is not found", func() {
			notFoundResp := responderGenerator(404, nil)

			It("should send nil to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", notFoundResp)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteRoutes("app_guid", errorCh, &wg)

				err := <-errorCh
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when deletion fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/routes/route1_guid", unassociateRouteResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/routes/route2_guid", unassociateRouteResp)
				httpmock.RegisterResponder("DELETE", "/v2/routes/route1_guid", deleteResp)
				httpmock.RegisterResponder("DELETE", "/v2/routes/route2_guid", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteRoutes("app_guid", errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when unassociating route fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/routes/route1_guid", unassociateRouteResp)
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/routes/route2_guid", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)

				go sut.DeleteRoutes("app_guid", errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
