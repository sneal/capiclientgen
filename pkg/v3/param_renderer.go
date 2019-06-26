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
	// ParamRendererStateUnknown is the default starting state
	ParamRendererStateUnknown ParamRendererStates = iota

	// ParamRendererStateName - params name cell
	ParamRendererStateName

	// ParamRendererStateType - params type cell
	ParamRendererStateType

	// ParamRendererStateDescription - description cell
	ParamRendererStateDescription

	// ParamRendererStateDefault - param default value cell
	ParamRendererStateDefault
)

func (r ParamRendererStates) String() string {
	switch r {
	case ParamRendererStateName:
		return "ParamRendererStateName"
	case ParamRendererStateType:
		return "ParamRendererStateType"
	case ParamRendererStateDescription:
		return "ParamRendererStateDescription"
	case ParamRendererStateDefault:
		return "ParamRendererStateDefault"
	}
	return "ParamRendererStateUnknown"
}

// ParamRenderer renders required parameters for a single endpoint
type ParamRenderer struct {
	state        ParamRendererStates
	endpoint     *service.Endpoint
	currentParam *service.Parameter
	required     bool
}

// NewParamRenderer creates a new v3 API required param renderer
func NewParamRenderer(endpoint *service.Endpoint, required bool) *ParamRenderer {
	return &ParamRenderer{
		endpoint: endpoint,
		state:    ParamRendererStateUnknown,
		required: required,
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
		w = r.tableRow(node, entering)
	case *ast.TableHeader:
		w = r.tableHeader(node)
	case *ast.TableCell:
		w = r.tableCell(node, entering)
	}
	return w
}

func (r *ParamRenderer) text(text *ast.Text) ast.WalkStatus {
	if r.state == ParamRendererStateName {
		v := string(text.Literal)
		if v != "" {
			r.currentParam.Name = r.currentParam.Name + v
		}
	}
	if r.state == ParamRendererStateType {
		v := string(text.Literal)
		if v != "" && r.currentParam.Type == "" {
			r.currentParam.Type = toOpenAPIDataType(v)
		}
	}
	if r.state == ParamRendererStateDescription {
		v := string(text.Literal)
		if v != "" {
			r.currentParam.Description = r.currentParam.Description + v
		}
	}
	if r.state == ParamRendererStateDefault {
		v := string(text.Literal)
		if v != "" {
			r.currentParam.Default = r.currentParam.Default + v
		}
	}
	return ast.GoToNext
}

func (r *ParamRenderer) tableHeader(tableHeader *ast.TableHeader) ast.WalkStatus {
	return ast.SkipChildren
}

func (r *ParamRenderer) tableRow(tableCell *ast.TableRow, entering bool) ast.WalkStatus {
	if entering {
		r.state = ParamRendererStateUnknown
		r.currentParam = &service.Parameter{
			Required: r.required,
		}
		r.endpoint.BodyParameters = append(r.endpoint.BodyParameters, r.currentParam)
	}
	return ast.GoToNext
}

func (r *ParamRenderer) tableCell(tableCell *ast.TableCell, entering bool) ast.WalkStatus {
	if entering {
		switch r.state {
		case ParamRendererStateUnknown:
			r.state = ParamRendererStateName
		case ParamRendererStateName:
			r.state = ParamRendererStateType
		case ParamRendererStateType:
			r.state = ParamRendererStateDescription
		case ParamRendererStateDescription:
			r.state = ParamRendererStateDefault
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
