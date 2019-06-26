package v3

//markdown.Render(doc, v3ApiRenderer)

import (
	"io"

	"github.com/gomarkdown/markdown/ast"

	"github.com/cloudfoundry-community/capiclientgen/pkg/service"
)

// ParamRendererStates represents the current state of parsing
type ParamRendererStates int

const (
	// ParamRendererStateDefault is the default starting state
	ParamRendererStateDefault ParamRendererStates = iota

	// ParamRendererStateRequiredParametersName - inside required params name cell
	ParamRendererStateRequiredParametersName

	// ParamRendererStateRequiredParametersType - inside required params type cell
	ParamRendererStateRequiredParametersType

	// ParamRendererStateRequiredParametersDescription - inside required params description cell
	ParamRendererStateRequiredParametersDescription
)

func (r ParamRendererStates) String() string {
	switch r {
	case ParamRendererStateDefault:
		return "ParamRendererStateDefault"
	case ParamRendererStateRequiredParametersName:
		return "ParamRendererStateRequiredParametersName"
	case ParamRendererStateRequiredParametersType:
		return "ParamRendererStateRequiredParametersType"
	case ParamRendererStateRequiredParametersDescription:
		return "ParamRendererStateRequiredParametersDescription"
	}

	return "Unknown"
}

// ParamRenderer renders required parameters for a single endpoint
type ParamRenderer struct {
	state    ParamRendererStates
	endpoint *service.Endpoint
}

// NewParamRenderer creates a new v3 API required param renderer
func NewParamRenderer(endpoint *service.Endpoint) *ParamRenderer {
	return &ParamRenderer{
		endpoint: endpoint,
		state:    ParamRendererStateDefault,
	}
}

// RenderHeader is a no-op
func (r *ParamRenderer) RenderHeader(w io.Writer, ast ast.Node) {

}

// RenderFooter is a no-op
func (r *ParamRenderer) RenderFooter(w io.Writer, ast ast.Node) {

}

// RenderNode renders a markdown node to HTML
func (r *ParamRenderer) RenderNode(ignored io.Writer, node ast.Node, entering bool) (w ast.WalkStatus) {
	switch node := node.(type) {
	case *ast.Text:
		w = r.text(node)
	case *ast.TableRow:
		w = r.tableRow(node)
	case *ast.TableHeader:
		w = r.tableHeader(node)
	case *ast.TableCell:
		w = r.tableCell(node, entering)
	}
	return w
}

func (r *ParamRenderer) text(text *ast.Text) ast.WalkStatus {
	if r.state == ParamRendererStateRequiredParametersName {
		v := string(text.Literal)
		if v != "" {
			p := service.Parameter{
				Name: v,
			}
			r.endpoint.BodyParameters = append(r.endpoint.BodyParameters, p)
		}
	}
	if r.state == ParamRendererStateRequiredParametersType {
		v := string(text.Literal)
		if v != "" {
			r.endpoint.BodyParameters[len(r.endpoint.BodyParameters)-1].Type = toOpenAPIDataType(v)
		}
	}
	if r.state == ParamRendererStateRequiredParametersDescription {
		v := string(text.Literal)
		if v != "" {
			r.endpoint.BodyParameters[len(r.endpoint.BodyParameters)-1].Description = v
		}
	}
	return ast.GoToNext
}

func (r *ParamRenderer) tableHeader(tableHeader *ast.TableHeader) ast.WalkStatus {
	return ast.SkipChildren
}

func (r *ParamRenderer) tableRow(tableCell *ast.TableRow) ast.WalkStatus {
	r.state = ParamRendererStateDefault
	return ast.GoToNext
}

func (r *ParamRenderer) tableCell(tableCell *ast.TableCell, entering bool) ast.WalkStatus {
	if entering {
		switch r.state {
		case ParamRendererStateDefault:
			r.state = ParamRendererStateRequiredParametersName
		case ParamRendererStateRequiredParametersName:
			r.state = ParamRendererStateRequiredParametersType
		case ParamRendererStateRequiredParametersType:
			r.state = ParamRendererStateRequiredParametersDescription
		}
	}
	return ast.GoToNext
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
