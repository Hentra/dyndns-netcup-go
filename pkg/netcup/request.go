package netcup

// Request represents a request to the netcup api.
type Request struct {
	Action string `json:"action"`
	Param  Params `json:"param"`
}

// NewRequest returns a new Request struct by a action and parameters.
func NewRequest(action string, params *Params) *Request {
	return &Request{
		Action: action,
		Param:  *params,
	}
}

// Params represent parameters to be send to the netcup api.
type Params map[string]interface{}

// NewParams returns a new Params struct.
func NewParams() Params {
	return make(map[string]interface{})
}

// AddParam adds a parameter with a specified key and value
func (p Params) AddParam(key string, value interface{}) {
	p[key] = value
}
