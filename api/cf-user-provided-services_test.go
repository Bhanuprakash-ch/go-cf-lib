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

var _ = Describe("cf ups", func() {

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

	Describe("create user provided service instance method", func() {
		req := types.CfUserProvidedService{}
		res := types.CfUserProvidedServiceResource{
			Meta: types.CfMeta{
				GUID: "guid",
			},
		}

		Context("when response code is CREATED", func() {

			resp := responderGenerator(201, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/user_provided_service_instances", resp)

				results, err := sut.CreateUserProvidedServiceInstance(&req)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when response code is ACCEPTED", func() {

			resp := responderGenerator(202, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/user_provided_service_instances", resp)

				results, err := sut.CreateUserProvidedServiceInstance(&req)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/user_provided_service_instances", negativeResponder)

				results, err := sut.CreateUserProvidedServiceInstance(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/user_provided_service_instances", requestFail)

				results, err := sut.CreateUserProvidedServiceInstance(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("get user provided service method", func() {
		res := types.CfUserProvidedServiceResource{
			Meta: types.CfMeta{
				GUID: "guid",
			},
		}

		Context("when found", func() {
			resp := responderGenerator(200, res)

			It("should not return error", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid", resp)

				results, err := sut.GetUserProvidedService("guid")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid", negativeResponder)

				results, err := sut.GetUserProvidedService("guid")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid", requestFail)

				results, err := sut.GetUserProvidedService("guid")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("create user provided service binding method", func() {
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

				results, err := sut.CreateUserProvidedServiceBinding(&req)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.Meta.GUID).To(Equal(res.Meta.GUID))
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_bindings", negativeResponder)

				results, err := sut.CreateUserProvidedServiceBinding(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_bindings", requestFail)

				results, err := sut.CreateUserProvidedServiceBinding(&req)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})

	Describe("delete user provided service instance method", func() {

		Context("when successfully deleted", func() {
			resp := responderGenerator(204, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/user_provided_service_instances/guid", resp)

				err := sut.DeleteUserProvidedServiceInstance("guid")

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/user_provided_service_instances/guid", negativeResponder)

				err := sut.DeleteUserProvidedServiceInstance("guid")

				Expect(err).NotTo(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/user_provided_service_instances/guid", requestFail)

				err := sut.DeleteUserProvidedServiceInstance("guid")

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("get user provided service bindings method", func() {
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
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid/service_bindings", resp)

				results, err := sut.GetUserProvidedServiceBindings("guid")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid/service_bindings", negativeResponder)

				results, err := sut.GetUserProvidedServiceBindings("guid")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/user_provided_service_instances/guid/service_bindings", requestFail)

				results, err := sut.GetUserProvidedServiceBindings("guid")

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})
})
