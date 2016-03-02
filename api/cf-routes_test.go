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
)

var _ = Describe("cf routes", func() {

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

	Describe("create route method", func() {
		req := types.CfCreateRouteRequest{}
		res := types.CfRouteResource{
			Meta: types.CfMeta{
				GUID: "guid",
			},
		}

		Context("with correct data passed", func() {

			resp := responderGenerator(201, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/routes", resp)

				results, err := sut.CreateRoute(&req)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/routes", negativeResponder)

				results, err := sut.CreateRoute(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/routes", requestFail)

				results, err := sut.CreateRoute(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("associate route method", func() {

		Context("when successful", func() {

			resp := responderGenerator(200, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("PUT", "/v2/apps/guid1/routes/guid2", resp)

				err := sut.AssociateRoute("guid1", "guid2")

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid1/routes/guid2", requestFail)

				err := sut.AssociateRoute("guid1", "guid2")

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("unassociate route method", func() {

		Context("when successful", func() {

			resp := responderGenerator(200, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/guid1/routes/guid2", resp)

				err := sut.UnassociateRoute("guid1", "guid2")

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/guid1/routes/guid2", requestFail)

				err := sut.UnassociateRoute("guid1", "guid2")

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("get app routes method", func() {
		res := types.CfRoutesResponse{
			Count: 1,
			Pages: 1,
			Resources: []types.CfRouteResource{
				types.CfRouteResource{},
			},
		}

		Context("when successful", func() {

			resp := responderGenerator(200, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid1/routes", resp)

				results, err := sut.GetAppRoutes("guid1")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Count).To(Equal(1))
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid1/routes", requestFail)

				results, err := sut.GetAppRoutes("guid1")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("get space route for hostname method", func() {
		res := types.CfRoutesResponse{
			Count: 1,
			Pages: 1,
			Resources: []types.CfRouteResource{
				types.CfRouteResource{},
			},
		}

		Context("when successful", func() {
			resp := responderGenerator(200, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("GET", "/v2/spaces/guid1/routes?q=host:hostname", resp)

				results, err := sut.GetSpaceRoutesForHostname("guid1", "hostname")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Count).To(Equal(1))
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/spaces/guid1/routes?q=host:hostname", requestFail)

				results, err := sut.GetSpaceRoutesForHostname("guid1", "hostname")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("get apps from route method", func() {
		res := types.CfRoutesResponse{
			Count: 1,
			Pages: 1,
			Resources: []types.CfRouteResource{
				types.CfRouteResource{},
			},
		}

		Context("when successful", func() {
			resp := responderGenerator(200, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("GET", "/v2/routes/guid1/apps", resp)

				results, err := sut.GetAppsFromRoute("guid1")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Count).To(Equal(1))
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/routes/guid1/apps", requestFail)

				results, err := sut.GetAppsFromRoute("guid1")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("delete route method", func() {

		Context("when successful", func() {
			resp := responderGenerator(204, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/routes/guid1", resp)

				err := sut.DeleteRoute("guid1")

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/routes/guid1", requestFail)

				err := sut.DeleteRoute("guid1")

				Expect(err).NotTo(BeNil())
			})
		})
	})

})
