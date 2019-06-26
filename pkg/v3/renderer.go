package v3

import (
	"fmt"
	"io"
	"strings"
	"github.com/gomarkdown/markdown/ast"

	"github.com/cloudfoundry-community/capiclientgen/pkg/service"
)


// RendererState represents the current state of parsing
type RendererStates int

const (
	RendererStateDefault RendererStates = iota
	RendererStateStartEndpoint
	RendererStateRequiredParameters
	RendererStatePermittedRoles
	RendererStateFinishEndpoint
)

type Renderer struct {
	ResourceName string
	endpoints []*service.Endpoint
	currentEndpoint *service.Endpoint
	state RendererStates
}

type RendererOptions struct {
	ResourceName string
}

func NewRenderer(opts RendererOptions) *Renderer {
	return &Renderer{
		ResourceName: opts.ResourceName,
		endpoints: []*service.Endpoint{},
		state: RendererStateDefault,
	}
}

func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {

}

func (r *Renderer) RenderFooter(w io.Writer, ast ast.Node) {
	
}

// RenderNode renders a markdown node to HTML
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Text:
		r.text(w, node)
	case *ast.Code:
		r.code(w, node)
	}
	return ast.GoToNext
}

func (r *Renderer) text(w io.Writer, text *ast.Text) {
		_, parentIsHeader := text.Parent.(*ast.Heading)
		if parentIsHeader {
			switch string(text.Literal) {
			case "Definition":
				r.state = RendererStateStartEndpoint
			case "Required Parameters":
				r.state = RendererStateRequiredParameters
			case "Permitted Roles":
				r.state = RendererStatePermittedRoles
		}
	}
}

func (r *Renderer) code(w io.Writer, node *ast.Code) {
	_, parentIsPara := node.Parent.(*ast.Paragraph)
	if parentIsPara && r.state == RendererStateStartEndpoint {
		s := strings.Split(string(node.Literal), " ")
		httpMethod := s[0]
		route := s[1]
		r.currentEndpoint = service.NewEndpoint(r.ResourceName, httpMethod, route)
		r.endpoints = append(r.endpoints, r.currentEndpoint)
		fmt.Println(fmt.Sprintf("Endpoint: %v+", r.currentEndpoint))
	}
}
