package v3

import (
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
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
func (r *Renderer) RenderNode(ignored io.Writer, node ast.Node, entering bool) (w ast.WalkStatus) {
	switch node := node.(type) {
	case *ast.Text:
		w = r.text(node)
	case *ast.Code:
		w = r.code(node)
	case *ast.TableBody:
		w = r.tableBody(node, entering)
	}
	return w
}

// ParsedEndpoints returns all the endpoint data parsed out
func (r *Renderer) ParsedEndpoints() []*service.Endpoint {
	return r.endpoints
}

func (r *Renderer) text(text *ast.Text) ast.WalkStatus {
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
	return ast.GoToNext
}

func (r *Renderer) code(node *ast.Code) ast.WalkStatus {
	_, parentIsPara := node.Parent.(*ast.Paragraph)
	if parentIsPara && r.state.Current() == RendererStateStartEndpoint {
		s := strings.Split(string(node.Literal), " ")
		httpMethod := s[0]
		route := s[1]
		r.currentEndpoint = service.NewEndpoint(r.ResourceName, httpMethod, route)
		r.endpoints = append(r.endpoints, r.currentEndpoint)
	}
	return ast.GoToNext
}

func (r *Renderer) tableBody(tableBody *ast.TableBody, entering bool) ast.WalkStatus {
	if r.state.Current() == RendererStateRequiredParameters {
		pr := NewParamRenderer(r.currentEndpoint)
		markdown.Render(tableBody, pr)
		return ast.SkipChildren
	}
	return ast.GoToNext
}
