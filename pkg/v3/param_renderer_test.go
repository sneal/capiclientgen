package v3_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"

	"github.com/cloudfoundry-community/capiclientgen/pkg/service"
	v3 "github.com/cloudfoundry-community/capiclientgen/pkg/v3"
)

var _ = Describe("Param Renderer", func() {
	Describe("rendering required parameters", func() {
		Context("with name and space params", func() {

			var endpoint *service.Endpoint

			BeforeEach(func() {
				t := []byte(`
Name | Type | Description
---- | ---- | -----------
**name** | _string_ | Name of the app.
**space** | [_to-one relationship_](#to-one-relationships) | A relationship to a space.
`)
				table := markdown.Parse(t, parser.NewWithExtensions(parser.CommonExtensions))
				endpoint = service.NewEndpoint("App", "POST", "/v3/apps")
				markdown.Render(table, v3.NewParamRenderer(endpoint))
			})

			It("contains all params", func() {
				Expect(endpoint.BodyParameters).To(HaveLen(2))
			})
			It("first param is named name", func() {
				Expect(endpoint.BodyParameters[0].Name).To(Equal("name"))
			})
			It("first param has type string", func() {
				Expect(endpoint.BodyParameters[0].Type).To(Equal("string"))
			})
			It("first param has description", func() {
				Expect(endpoint.BodyParameters[0].Description).To(Equal("Name of the app."))
			})
			It("second param is named space", func() {
				Expect(endpoint.BodyParameters[1].Name).To(Equal("space"))
			})
			It("second param has type object", func() {
				Expect(endpoint.BodyParameters[1].Type).To(Equal("object"))
			})
			It("second param has description", func() {
				Expect(endpoint.BodyParameters[1].Description).To(Equal("A relationship to a space."))
			})
		})
	})
})
