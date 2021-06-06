package hue

type Light struct {
	State    State  `json:"state"`
	Name     string `json:"name"`
	UniqueID string `json:"uniqueid"`
}

type State struct {
	On         bool      `json:"on,omitempty"`
	Hue        int       `json:"hue,omitempty"`
	Saturation int       `json:"sat,omitempty"`
	Brightness int       `json:"bri,omitempty"`
	Effect     string    `json:"effect,omitempty"`
	XY         []float64 `json:"xy,omitempty"`
	Ct         int       `json:"ct,omitempty"`
	Alert      string    `json:"alert,omitempty"`
	ColorMode  string    `json:"colormode,omitempty"`
	Mode       string    `json:"mode,omitempty"`
	Reachable  bool      `json:"reachable,omitempty"`
}

type Config struct {
	Endpoint string
	Username string
}
