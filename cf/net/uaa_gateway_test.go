package net_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry/cli/cf/configuration/coreconfig"
	"github.com/cloudfoundry/cli/cf/errors"
	. "github.com/cloudfoundry/cli/cf/net"
	"github.com/cloudfoundry/cli/cf/terminal/terminalfakes"
	"github.com/cloudfoundry/cli/cf/trace/tracefakes"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var failingUAARequest = func(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
	jsonResponse := `{ "error": "foo", "error_description": "The foo is wrong..." }`
	fmt.Fprintln(writer, jsonResponse)
}

var _ = Describe("UAA Gateway", func() {
	var gateway Gateway
	var config coreconfig.Reader

	BeforeEach(func() {
		config = testconfig.NewRepository()
		gateway = NewUAAGateway(config, new(terminalfakes.FakeUI), new(tracefakes.FakePrinter))
	})

	It("parses error responses", func() {
		ts := httptest.NewTLSServer(http.HandlerFunc(failingUAARequest))
		defer ts.Close()
		gateway.SetTrustedCerts(ts.TLS.Certificates)

		request, apiErr := gateway.NewRequest("GET", ts.URL, "TOKEN", nil)
		_, apiErr = gateway.PerformRequest(request)

		Expect(apiErr).NotTo(BeNil())
		Expect(apiErr.Error()).To(ContainSubstring("The foo is wrong"))
		Expect(apiErr.(errors.HTTPError).ErrorCode()).To(ContainSubstring("foo"))
	})
})
