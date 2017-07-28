package uaago_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	. "github.com/cf-platform-eng/uaago/uaago"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("uaa_token_fetcher", func() {
	var testServer *httptest.Server
	var capturedRequest *http.Request

	BeforeEach(func() {
		testServer = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			capturedRequest = request
			jsonData := []byte(`{
				"access_token":"some-token",
				"token_type":"bearer",
				"expires_in":599,
				"scope":"cloud_controller.write doppler.firehose",
				"jti":"28edda5c-4e37-4a63-9ba3-b32f48530a51"
			}`)
			writer.Write(jsonData)
		}))
	})

	AfterEach(func() {
		testServer.Close()
	})

	It("requests token using credentials", func() {
		fetcher := NewUAATokenFetcher(testServer.URL, "user", "pass", true)
		fetcher.RefreshAuthToken()

		authValue := capturedRequest.Header.Get("Authorization")
		expectedAuthValue := "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass"))

		//Client struct doesn't implement interface, avhe to basically re-test
		Expect(authValue).To(Equal(expectedAuthValue))
	})

	It("returns token", func() {
		fetcher := NewUAATokenFetcher(testServer.URL, "user", "password", true)
		token, err := fetcher.RefreshAuthToken()

		//Client struct doesn't implement interface, avhe to basically re-test
		Expect(err).To(BeNil())
		Expect(token).To(Equal("bearer some-token"))
	})
})
