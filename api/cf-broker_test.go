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

var _ = Describe("cf broker", func() {

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

	Describe("register broker method", func() {
		name := "name"
		url := "url"
		user := "user"
		pass := "password"

		Context("with correct data passed", func() {

			resp := responderGenerator(201, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_brokers", resp)

				err := sut.RegisterBroker(name, url, user, pass)

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_brokers", negativeResponder)

				err := sut.RegisterBroker(name, url, user, pass)

				Expect(err).NotTo(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/service_brokers", requestFail)

				err := sut.RegisterBroker(name, url, user, pass)

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("update broker method", func() {
		guid := "guid"
		url := "url"
		user := "user"
		pass := "password"

		Context("with correct data passed", func() {

			resp := responderGenerator(200, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("PUT", "/v2/service_brokers/guid", resp)

				err := sut.UpdateBroker(guid, url, user, pass)

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("PUT", "/v2/service_brokers/guid", negativeResponder)

				err := sut.UpdateBroker(guid, url, user, pass)

				Expect(err).NotTo(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("PUT", "/v2/service_brokers/guid", requestFail)

				err := sut.UpdateBroker(guid, url, user, pass)

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("get brokers method", func() {
		name := "name"
		brokers := types.CfServiceBrokerResources{
			TotalResults: 1,
			Resources: []types.CfServiceBrokerResource{
				types.CfServiceBrokerResource{},
			},
		}

		Context("when succeeded", func() {
			resp := responderGenerator(200, brokers)

			It("should return resources matching", func() {
				httpmock.RegisterResponder("GET", "/v2/service_brokers?q=name:name", resp)

				results, err := sut.GetBrokers(name)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(results).NotTo(BeNil())
				Expect(results.TotalResults).To(Equal(brokers.TotalResults))
			})
		})
		Context("when CF responds with not expected body", func() {
			resp := responderGenerator(200, "{\"syntax|':\"true\"")

			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/service_brokers?q=name:name", resp)

				results, err := sut.GetBrokers(name)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/service_brokers?q=name:name", negativeResponder)

				results, err := sut.GetBrokers(name)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/service_brokers?q=name:name", requestFail)

				results, err := sut.GetBrokers(name)

				Expect(err).NotTo(BeNil())
				Expect(results).To(BeNil())
			})
		})
	})
})
