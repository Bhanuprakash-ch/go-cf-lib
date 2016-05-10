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

var _ = Describe("Cf create", func() {

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

	Describe("create service clone method", func() {
		appSummary := types.CfAppSummary{
			GUID: "app_guid",
			Services: []types.CfAppSummaryService{
				types.CfAppSummaryService{
					GUID: "svc_guid",
					Name: "service_name",
					Plan: types.CfAppSummaryServicePlan{
						GUID: "plan_guid",
						Name: "Simple",
					},
				},
			},
		}
		params := make(map[string]interface{})
		getAppSummaryResp := responderGenerator(200, appSummary)

		createSvcInstanceRes := types.CfServiceInstanceCreateResponse{
			Meta: types.CfMeta{
				GUID: "svc_clone_guid",
			},
		}
		createSvcInstanceResp := responderGenerator(201, createSvcInstanceRes)

		comp := types.Component{
			GUID: "svc_guid",
			Name: "Name",
			DependencyOf: []string{
				"app_guid",
			},
		}
		suffix := "abcde"

		Context("without errors", func() {
			It("should send nil to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("POST", "/v2/service_instances?accepts_incomplete=false", createSvcInstanceResp)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)
				resultsCh := make(chan types.ComponentClone, 1)

				go sut.CreateServiceClone("space_guid", params, comp, suffix, resultsCh, errorCh, &wg)

				err := <-errorCh
				Expect(err).ShouldNot(HaveOccurred())

				result := <-resultsCh
				Expect(result.CloneGUID).To(Equal(createSvcInstanceRes.Meta.GUID))
				Expect(result.Component.GUID).To(Equal(comp.GUID))
			})
		})
		Context("when passed not bound component", func() {
			notBound := types.Component{
				GUID:         "svc_guid",
				Name:         "Name",
				DependencyOf: []string{},
			}

			It("should send error to channel", func() {
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)
				resultsCh := make(chan types.ComponentClone, 1)

				go sut.CreateServiceClone("space_guid", params, notBound, suffix, resultsCh, errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when getting app summary fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)
				resultsCh := make(chan types.ComponentClone, 1)

				go sut.CreateServiceClone("space_guid", params, comp, suffix, resultsCh, errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when creation of service fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("POST", "/v2/service_instances?accepts_incomplete=false", negativeResponder)
				wg := sync.WaitGroup{}
				wg.Add(1)
				errorCh := make(chan error, 1)
				resultsCh := make(chan types.ComponentClone, 1)

				go sut.CreateServiceClone("space_guid", params, comp, suffix, resultsCh, errorCh, &wg)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("create application clone method", func() {
		appSummary := types.CfAppSummary{
			GUID: "guid",
			Routes: []types.CfAppSummaryRoute{
				types.CfAppSummaryRoute{
					GUID: "route_guid",
				},
			},
		}
		getAppSummaryResp := responderGenerator(200, appSummary)

		appToReturn := types.CfAppResource{
			Meta: types.CfMeta{
				GUID: "app_clone_guid",
			},
			Entity: types.CfApp{
				Name: "name",
			},
		}
		positiveCreateAppResponder := responderGenerator(201, appToReturn)

		routeRes := types.CfRouteResource{
			Meta: types.CfMeta{
				GUID: "route_guid",
			},
		}
		createRouteResp := responderGenerator(201, routeRes)

		associateRouteResp := responderGenerator(200, nil)

		params := make(map[string]string)
		params["name"] = "name"
		params["key"] = "value"

		Context("without errors", func() {
			It("should send nil to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("POST", "/v2/apps", positiveCreateAppResponder)
				httpmock.RegisterResponder("POST", "/v2/routes", createRouteResp)
				httpmock.RegisterResponder("PUT", "/v2/apps/app_clone_guid/routes/route_guid", associateRouteResp)

				result, err := sut.CreateApplicationClone("app_guid", "space_guid", params)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Meta.GUID).To(Equal("app_clone_guid"))
			})
		})
		Context("when getting application summary fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", negativeResponder)

				_, err := sut.CreateApplicationClone("app_guid", "space_guid", params)

				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when reference app has no route associated", func() {
			appSummaryWithoutRoutes := types.CfAppSummary{
				GUID:   "guid",
				Routes: []types.CfAppSummaryRoute{},
			}
			getAppSummaryResp := responderGenerator(200, appSummaryWithoutRoutes)

			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)

				_, err := sut.CreateApplicationClone("app_guid", "space_guid", params)

				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when app creation fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("POST", "/v2/apps", negativeResponder)

				_, err := sut.CreateApplicationClone("app_guid", "space_guid", params)

				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when route creation fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("POST", "/v2/apps", positiveCreateAppResponder)
				httpmock.RegisterResponder("POST", "/v2/routes", negativeResponder)

				_, err := sut.CreateApplicationClone("app_guid", "space_guid", params)

				Expect(err).Should(HaveOccurred())
			})
		})
		Context("when route association fails", func() {
			It("should send error to channel", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/app_guid/summary", getAppSummaryResp)
				httpmock.RegisterResponder("POST", "/v2/apps", positiveCreateAppResponder)
				httpmock.RegisterResponder("POST", "/v2/routes", createRouteResp)
				httpmock.RegisterResponder("PUT", "/v2/apps/app_guid/routes/route_guid", negativeResponder)

				_, err := sut.CreateApplicationClone("app_guid", "space_guid", params)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
