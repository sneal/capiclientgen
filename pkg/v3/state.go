package v3

// RendererStates represents the current state of parsing
type RendererStates int

const (
	// RendererStateDefault is the default starting state
	RendererStateDefault RendererStates = iota

	// RendererStateStartEndpoint - found new endpoint
	RendererStateStartEndpoint

	// RendererStateRequiredParameters - inside the required parameters section
	RendererStateRequiredParameters

	// RendererStatePermittedRoles - inside the permitted roles section
	RendererStatePermittedRoles
)

func (r RendererStates) String() string {
	switch r {
	case RendererStateStartEndpoint:
		return "RendererStateStartEndpoint"
	case RendererStateRequiredParameters:
		return "RendererStateRequiredParameters"
	case RendererStatePermittedRoles:
		return "RendererStatePermittedRoles"
	}

	return "Unknown"
}

// RenderStates keeps track of the current renderer state
type RenderStates struct {
	state RendererStates
}

// NewRenderState creates a new render state
func NewRenderState() *RenderStates {
	return &RenderStates{
		state: RendererStateDefault,
	}
}

// SetRendererState sets the renderer to the specified state
func (rs *RenderStates) SetRendererState(state RendererStates) {
	//fmt.Println(fmt.Sprintf("Render state: %s", state))
	rs.state = state
}

// Current renderer state
func (rs *RenderStates) Current() RendererStates {
	return rs.state
}

// Reset back to the default state
func (rs *RenderStates) Reset() {
	rs.state = RendererStateDefault
}
