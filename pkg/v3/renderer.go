package v3

import (
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"

	"github.com/cloudfoundry-community/capiclientgen/pkg/service"
)

// Renderer parses CAPI V3 markdown into a structured doc
type Renderer struct {
	ResourceName    string
	endpoints       []*service.Endpoint
	currentEndpoint *service.Endpoint
	state           *RenderStates
}

// RendererOptions to create a new Renderer
type RendererOptions struct {
	ResourceName string
}

// NewRenderer creates a new v3 API renderer
func NewRenderer(opts RendererOptions) *Renderer {
	return &Renderer{
		ResourceName: opts.ResourceName,
		endpoints:    []*service.Endpoint{},
		state:        NewRenderState(),
	}
}

// RenderHeader is a no-op
func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {

}

// RenderFooter is a no-op
func (r *Renderer) RenderFooter(w io.Writer, ast ast.Node) {

}

// RenderNode renders a markdown node to HTML
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Text:
		r.text(node)
	case *ast.Code:
		r.code(node)
	case *ast.TableCell:
		r.tableCell(node, entering)
	}
	return ast.GoToNext
}

// ParsedEndpoints returns all the endpoint data parsed out
func (r *Renderer) ParsedEndpoints() []*service.Endpoint {
	return r.endpoints
}

func (r *Renderer) text(text *ast.Text) {
	_, parentIsHeader := text.Parent.(*ast.Heading)
	if parentIsHeader {
		switch string(text.Literal) {
		case "Definition":
			r.state.SetRendererState(RendererStateStartEndpoint)
		case "Required Parameters":
			r.state.SetRendererState(RendererStateRequiredParameters)
		case "Permitted Roles":
			r.state.SetRendererState(RendererStatePermittedRoles)
		}
	}
	if r.state.Current() == RendererStateRequiredParametersName {
		v := string(text.Literal)
		if v != "" {
			p := service.Parameter{
				Name: v,
			}
			r.currentEndpoint.BodyParameters = append(r.currentEndpoint.BodyParameters, p)
		}
	}
	if r.state.Current() == RendererStateRequiredParametersType {
		v := string(text.Literal)
		if v != "" {
			r.currentEndpoint.BodyParameters[len(r.currentEndpoint.BodyParameters)-1].Type = toOpenAPIDataType(v)
		}
	}
	if r.state.Current() == RendererStateRequiredParametersDescription {
		v := string(text.Literal)
		if v != "" {
			r.currentEndpoint.BodyParameters[len(r.currentEndpoint.BodyParameters)-1].Description = v
		}
	}
}

func (r *Renderer) code(node *ast.Code) {
	_, parentIsPara := node.Parent.(*ast.Paragraph)
	if parentIsPara && r.state.Current() == RendererStateStartEndpoint {
		s := strings.Split(string(node.Literal), " ")
		httpMethod := s[0]
		route := s[1]
		r.currentEndpoint = service.NewEndpoint(r.ResourceName, httpMethod, route)
		r.endpoints = append(r.endpoints, r.currentEndpoint)
	}
}

func (r *Renderer) tableCell(tableCell *ast.TableCell, entering bool) {
	if !entering || !isTableBodyCell(tableCell) {
		return
	}

	switch r.state.Current() {
	case RendererStateRequiredParameters:
		r.state.SetRendererState(RendererStateRequiredParametersName)
	case RendererStateRequiredParametersName:
		r.state.SetRendererState(RendererStateRequiredParametersType)
	case RendererStateRequiredParametersType:
		r.state.SetRendererState(RendererStateRequiredParametersDescription)
	}
}

// Normalize data types to https://swagger.io/docs/specification/data-models/data-types/
func toOpenAPIDataType(v string) string {
	switch v {
	case "obect":
		fallthrough
	case "to-one relationship":
		return "object"
	case "array":
		fallthrough
	case "to-many relationship":
		return "array"
	case "string":
		return "string"
	case "number":
		return "number"
	case "integer":
		return "integer"
	case "boolean":
		return "boolean"
	default:
		return "object"
	}
}
