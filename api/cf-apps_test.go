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
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/trustedanalytics/go-cf-lib/types"
	"net/http"
)

var _ = Describe("Cf apps", func() {

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

	Describe("create app method", func() {
		appToReturn := types.CfAppResource{}
		appToReturn.Meta = types.CfMeta{GUID: "super_fake_guid"}
		positiveCreateAppResponder := responderGenerator(201, appToReturn)

		Context("with correct data passed", func() {
			It("should respond with guid", func() {
				httpmock.RegisterResponder("POST", "/v2/apps", positiveCreateAppResponder)

				result, err := sut.CreateApp(types.CfApp{Name: "appName"})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Meta.GUID).To(Equal("super_fake_guid"))
			})
		})
		Context("when CF responds with different status code", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps", negativeResponder)

				result, err := sut.CreateApp(types.CfApp{Name: "appName"})

				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps", requestFail)

				result, err := sut.CreateApp(types.CfApp{Name: "appName"})

				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("get app summary method", func() {
		Context("when entity exists", func() {
			appSummary := types.CfAppSummary{GUID: "guid"}
			resp := responderGenerator(200, appSummary)

			It("should retrieve entity", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid/summary", resp)

				result, err := sut.GetAppSummary("guid")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.GUID).To(Equal("guid"))
			})
		})
		Context("when entity was not found", func() {
			It("should return error", func() {
				resp := responderGenerator(404, nil)
				httpmock.RegisterResponder("GET", "/v2/apps/guid/summary", resp)

				result, err := sut.GetAppSummary("guid")

				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			})
		})
		Context("when response body is wrong", func() {
			It("should return error", func() {
				resp := responderGenerator(200, "{\"mishmash\":\"\"")
				httpmock.RegisterResponder("GET", "/v2/apps/guid/summary", resp)

				result, err := sut.GetAppSummary("guid")

				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			})
		})
		Context("when http request fail", func() {
			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid/summary", requestFail)

				result, err := sut.GetAppSummary("guid")

				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("assert app has routes method", func() {
		Context("when app summary does not have routes", func() {
			appSummaryWithoutRoutes := types.CfAppSummary{GUID: "guid", Routes: []types.CfAppSummaryRoute{}}

			It("should return error", func() {
				err := sut.AssertAppHasRoutes(&appSummaryWithoutRoutes)

				Expect(err).To(HaveOccurred())
			})
		})
		Context("when app summary has routes", func() {
			appSummaryWithRoute := types.CfAppSummary{
				GUID: "guid",
				Routes: []types.CfAppSummaryRoute{
					types.CfAppSummaryRoute{
						GUID: "guid",
						Host: "hostname",
						Domain: types.CfDomain{
							GUID: "guid2",
							Name: "example.com",
						},
					},
				},
			}

			It("should not return error", func() {

				err := sut.AssertAppHasRoutes(&appSummaryWithRoute)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("delete app", func() {
		Context("when successfully deleted", func() {
			resp := responderGenerator(204, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/guid", resp)

				err := sut.DeleteApp("guid")

				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when entity not found", func() {
			resp := responderGenerator(404, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/guid", resp)

				err := sut.DeleteApp("guid")

				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when internal server error occur", func() {
			resp := responderGenerator(500, nil)

			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/guid", resp)

				err := sut.DeleteApp("guid")

				Expect(err).To(HaveOccurred())
			})
		})
		Context("when request fails", func() {
			resp := responderFailGenerator(nil)

			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/guid", resp)

				err := sut.DeleteApp("guid")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("get app bindings", func() {
		Context("when found", func() {
			resources := types.CfBindingsResources{
				TotalResults: 1,
				Resources:    []types.CfBindingResource{},
			}
			resp := responderGenerator(200, resources)

			It("should return results", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid/service_bindings", resp)

				bindings, err := sut.GetAppBindigs("guid")

				Expect(err).NotTo(HaveOccurred())
				Expect(bindings).NotTo(BeNil())
				Expect(bindings.TotalResults).To(Equal(resources.TotalResults))
			})
		})
		Context("when entity not found", func() {
			resp := responderGenerator(404, nil)

			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid/service_bindings", resp)

				bindings, err := sut.GetAppBindigs("guid")

				Expect(err).To(HaveOccurred())
				Expect(bindings).To(BeNil())
			})
		})
		Context("when internal server error occur", func() {
			resp := responderGenerator(500, nil)

			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid/service_bindings", resp)

				bindings, err := sut.GetAppBindigs("guid")

				Expect(err).To(HaveOccurred())
				Expect(bindings).To(BeNil())
			})
		})
		Context("when request fails", func() {
			resp := responderFailGenerator(nil)

			It("should return error", func() {
				httpmock.RegisterResponder("GET", "/v2/apps/guid/service_bindings", resp)

				bindings, err := sut.GetAppBindigs("guid")

				Expect(err).To(HaveOccurred())
				Expect(bindings).To(BeNil())
			})
		})
	})

	Describe("delete binding", func() {
		binding := types.CfBindingResource{
			Meta: types.CfMeta{
				GUID: "guid",
			},
			Entity: types.CfBinding{
				AppGUID: "app_guid",
			},
		}
		Context("when successfully deleted", func() {
			resp := responderGenerator(204, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/service_bindings/guid", resp)

				err := sut.DeleteBinding(binding)

				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when entity not found", func() {
			resp := responderGenerator(404, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/service_bindings/guid", resp)

				err := sut.DeleteBinding(binding)

				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when internal server error occur", func() {
			resp := responderGenerator(500, nil)

			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/service_bindings/guid", resp)

				err := sut.DeleteBinding(binding)

				Expect(err).To(HaveOccurred())
			})
		})
		Context("when request fails", func() {
			resp := responderFailGenerator(nil)

			It("should return error", func() {
				httpmock.RegisterResponder("DELETE", "/v2/apps/app_guid/service_bindings/guid", resp)

				err := sut.DeleteBinding(binding)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("copy bits", func() {
		Context("when successfully copied", func() {
			jobResponse := types.CfJobResponse{
				Meta: types.CfMeta{
					URL: "/v2/jobs/guid",
				},
				Entity: types.CfJob{
					Status: "finished",
				},
			}
			resp := responderGenerator(201, jobResponse)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid/copy_bits", resp)
				errorCh := make(chan error)

				go sut.CopyBits("source", "guid", errorCh)

				err := <-errorCh
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when CF responds with wrong status code", func() {
			resp := responderGenerator(500, nil)

			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid/copy_bits", resp)
				errorCh := make(chan error)

				go sut.CopyBits("source", "guid", errorCh)

				err := <-errorCh
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when job failed", func() {
			jobResponse := types.CfJobResponse{
				Meta: types.CfMeta{
					URL: "/v2/jobs/guid",
				},
				Entity: types.CfJob{
					Status: "failed",
				},
			}
			resp := responderGenerator(201, jobResponse)

			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid/copy_bits", resp)
				httpmock.RegisterResponder("GET", "/v2/jobs/guid", resp)
				errorCh := make(chan error)

				go sut.CopyBits("source", "guid", errorCh)

				err := <-errorCh
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when job queued", func() {
			jobResponse := types.CfJobResponse{
				Meta: types.CfMeta{
					URL: "/v2/jobs/guid",
				},
				Entity: types.CfJob{
					Status: "queued",
				},
			}

			resp := responderGenerator(201, jobResponse)

			Context("when status changed to finished", func() {
				finished := jobResponse
				finished.Entity.Status = "finished"
				resp2 := responderGenerator(200, finished)

				It("should not return error", func() {

					httpmock.RegisterResponder("POST", "/v2/apps/guid/copy_bits", resp)
					httpmock.RegisterResponder("GET", "/v2/jobs/guid", resp2)
					errorCh := make(chan error)

					go sut.CopyBits("source", "guid", errorCh)

					err := <-errorCh
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			Context("when status changed to failed", func() {
				finished := jobResponse
				finished.Entity.Status = "failed"
				resp2 := responderGenerator(200, finished)

				It("should not return error", func() {

					httpmock.RegisterResponder("POST", "/v2/apps/guid/copy_bits", resp)
					httpmock.RegisterResponder("GET", "/v2/jobs/guid", resp2)
					errorCh := make(chan error)

					go sut.CopyBits("source", "guid", errorCh)

					err := <-errorCh
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("when request for status failed", func() {
				resp2 := responderFailGenerator(nil)

				It("should not return error", func() {

					httpmock.RegisterResponder("POST", "/v2/apps/guid/copy_bits", resp)
					httpmock.RegisterResponder("GET", "/v2/jobs/guid", resp2)
					errorCh := make(chan error)

					go sut.CopyBits("source", "guid", errorCh)

					err := <-errorCh
					Expect(err).Should(HaveOccurred())
				})
			})

		})
		Context("when request fails", func() {
			resp := responderFailGenerator(nil)

			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid/copy_bits", resp)
				errorCh := make(chan error)

				go sut.CopyBits("source", "guid", errorCh)

				err := <-errorCh
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("restage app", func() {
		Context("when successfully deleted", func() {
			resp := responderGenerator(201, nil)

			It("should not return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid/restage", resp)

				err := sut.RestageApp("guid")

				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when internal server error occur", func() {
			resp := responderGenerator(500, nil)

			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid/restage", resp)

				err := sut.RestageApp("guid")

				Expect(err).To(HaveOccurred())
			})
		})
		Context("when request fails", func() {
			resp := responderFailGenerator(nil)

			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid/restage", resp)

				err := sut.RestageApp("guid")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("update app", func() {
		executeTestCase := func(responder httpmock.Responder) error {
			guid, _ := uuid.NewV4()
			app := types.CfAppResource{}
			app.Meta.GUID = guid.String()
			updateURL := fmt.Sprintf("/v2/apps/%v", app.Meta.GUID)

			httpmock.RegisterResponder("PUT", updateURL, responder)
			sut := CfAPI{Client: http.DefaultClient}
			return sut.UpdateApp(&app)
		}

		Context("in positive scenario", func() {
			positiveStatusCreatedResponder := responderGenerator(201, nil)

			It("should not return error", func() {
				err := executeTestCase(positiveStatusCreatedResponder)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("in negative scenario", func() {
			It("should propagate error", func() {
				err := executeTestCase(negativeResponder)
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("when request fails", func() {
			resp := responderFailGenerator(nil)

			It("should return error", func() {
				httpmock.RegisterResponder("POST", "/v2/apps/guid", resp)
				app := types.CfAppResource{}
				app.Meta.GUID = "guid"

				err := sut.UpdateApp(&app)

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("start app", func() {
		app := types.CfAppResource{
			Meta: types.CfMeta{
				GUID: "guid",
			},
			Entity: types.CfApp{
				State: "STOPPED",
			},
		}
		decodedInstances := make(map[string]types.CfAppInstance)
		decodedInstances["guid"] = types.CfAppInstance{
			State: "RUNNING",
		}
		responder := responderGenerator(201, app)

		Context("when successfully started", func() {
			resp := responderGenerator(200, decodedInstances)

			It("should not return error", func() {
				httpmock.RegisterResponder("PUT", "/v2/apps/guid", responder)
				httpmock.RegisterResponder("GET", "/v2/apps/guid/instances", resp)

				err := sut.StartApp(&app)

				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("when request fails", func() {
			resp := responderFailGenerator(nil)

			It("should return error", func() {
				httpmock.RegisterResponder("PUT", "/v2/apps/guid", responder)
				httpmock.RegisterResponder("GET", "/v2/apps/guid/instances", resp)

				err := sut.StartApp(&app)

				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func responderGenerator(code int, v interface{}) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		resp, _ := httpmock.NewJsonResponse(code, v)
		return resp, nil
	}
}

func responderFailGenerator(err error) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		if err == nil {
			err = errors.New("Request failed")
		}
		return nil, err
	}
}
