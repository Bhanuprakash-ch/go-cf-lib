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

var _ = Describe("cf services", func() {

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

	Describe("create service instance method", func() {
		req := types.CfServiceInstanceCreateRequest{}
		res := types.CfServiceInstanceCreateResponse{
			Meta: types.CfMeta{
				GUID: "guid",
			},
		}

		Context("when response code is CREATED", func() {

			resp := responderGenerator(201, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_instances?accepts_incomplete=false", resp)

				results, err := sut.CreateServiceInstance(&req)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when response code is ACCEPTED", func() {

			resp := responderGenerator(202, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_instances?accepts_incomplete=false", resp)

				results, err := sut.CreateServiceInstance(&req)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_instances?accepts_incomplete=false", negativeResponder)

				results, err := sut.CreateServiceInstance(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_instances?accepts_incomplete=false", requestFail)

				results, err := sut.CreateServiceInstance(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("create service binding method", func() {
		req := types.CfServiceBindingCreateRequest{}
		res := types.CfServiceBindingCreateResponse{
			Meta: types.CfMeta{
				GUID: "guid",
			},
		}

		Context("when response code is CREATED", func() {
			resp := responderGenerator(201, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_bindings", resp)

				results, err := sut.CreateServiceBinding(&req)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_bindings", negativeResponder)

				results, err := sut.CreateServiceBinding(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_bindings", requestFail)

				results, err := sut.CreateServiceBinding(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("get service bindings method", func() {
		res := types.CfBindingsResources{
			TotalResults: 1,
			Resources: []types.CfBindingResource{
				types.CfBindingResource{
					Meta: types.CfMeta{
						GUID: "guid1",
					},
				},
			},
		}

		Context("when response code is CREATED", func() {
			resp := responderGenerator(200, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("GET", "/v2/service_instances/guid/service_bindings", resp)

				results, err := sut.GetServiceBindings("guid")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/service_instances/guid/service_bindings", negativeResponder)

				results, err := sut.GetServiceBindings("guid")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/service_instances/guid/service_bindings", requestFail)

				results, err := sut.GetServiceBindings("guid")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("delete service instance method", func() {

		Context("when successfully deleted", func() {
			resp := responderGenerator(204, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/service_instances/guid", resp)

				err := sut.DeleteServiceInstance("guid")

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/service_instances/guid", negativeResponder)

				err := sut.DeleteServiceInstance("guid")

				Expect(err).NotTo(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/service_instances/guid", requestFail)

				err := sut.DeleteServiceInstance("guid")

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("get service of name method", func() {

		Context("when service is found", func() {
			res := types.CfBindingsResources{
				TotalResults: 1,
				Resources: []types.CfBindingResource{
					types.CfBindingResource{
						Meta: types.CfMeta{
							GUID: "guid1",
						},
					},
				},
			}
			resp := responderGenerator(200, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("GET", "/v2/services?q=label:name", resp)

				results, err := sut.GetServiceOfName("name")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal("guid1"))
			})
		})
		Context("when no services are found", func() {
			res := types.CfBindingsResources{
				TotalResults: 0,
				Resources:    []types.CfBindingResource{},
			}
			resp := responderGenerator(200, res)

			It("should not return error but return nil as response", func() {
				httpmock.RegisterResponder("GET", "/v2/services?q=label:name", resp)

				results, err := sut.GetServiceOfName("name")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).To(BeNil())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/services?q=label:name", negativeResponder)

				results, err := sut.GetServiceOfName("name")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/services?q=label:name", requestFail)

				results, err := sut.GetServiceOfName("name")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("purge service method", func() {
		res := types.CfServicePlansResources{
			TotalResults: 1,
			Resources: []types.CfServicePlanResource{
				types.CfServicePlanResource{
					Meta: types.CfMeta{
						GUID: "guid",
					},
				},
			},
		}
		resp := responderGenerator(200, res)
		resp2 := responderGenerator(204, nil)

		Context("when all is ok", func() {
			It("should not return error", func() {
				httpmock.RegisterResponder("GET", "/planUrl", resp)
				httpmock.RegisterResponder("DELETE", "/v2/service_plans/guid", resp2)
				httpmock.RegisterResponder("DELETE", "/v2/services/serviceID", resp2)

				err := sut.PurgeService("serviceID", "serviceName", "/planUrl")

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when service instance is not found", func() {
			resp3 := responderGenerator(404, nil)
			It("should not return error as it is already deleted", func() {
				httpmock.RegisterResponder("GET", "/planUrl", resp3)
				httpmock.RegisterResponder("DELETE", "/v2/service_plans/guid", resp2)
				httpmock.RegisterResponder("DELETE", "/v2/services/serviceID", resp2)

				err := sut.PurgeService("serviceID", "serviceName", "/planUrl")

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when service instance request is rejected", func() {
			resp4 := responderGenerator(400, nil)
			It("should not return error as it is already deleted", func() {
				httpmock.RegisterResponder("GET", "/planUrl", resp4)
				httpmock.RegisterResponder("DELETE", "/v2/service_plans/guid", resp2)
				httpmock.RegisterResponder("DELETE", "/v2/services/serviceID", resp2)

				err := sut.PurgeService("serviceID", "serviceName", "/planUrl")

				Expect(err).NotTo(BeNil())
			})
		})
		Context("when getting plan fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/planUrl", requestFail)

				err := sut.PurgeService("serviceID", "serviceName", "/planUrl")

				Expect(err).NotTo(BeNil())
			})
		})
		Context("when deleting service plan fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/planUrl", resp)
				httpmock.RegisterResponder("DELETE", "/v2/service_plans/guid", requestFail)

				err := sut.PurgeService("serviceID", "serviceName", "/planUrl")

				Expect(err).NotTo(BeNil())
			})
		})
		Context("when deleting service instance fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/planUrl", resp)
				httpmock.RegisterResponder("DELETE", "/v2/service_plans/guid", resp2)
				httpmock.RegisterResponder("DELETE", "/v2/services/serviceID", requestFail)

				err := sut.PurgeService("serviceID", "serviceName", "/planUrl")

				Expect(err).NotTo(BeNil())
			})
		})

	})

})
