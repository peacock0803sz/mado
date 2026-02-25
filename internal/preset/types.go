// Package preset implements YAML preset management for window layouts.
package preset

// Preset is a named window layout definition loaded from the config file.
type Preset struct {
	Name        string `json:"name"        yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Rules       []Rule `json:"rules"       yaml:"rules"`
}

// Rule is a single window operation instruction within a preset.
type Rule struct {
	App      string `json:"app"                yaml:"app"`
	Title    string `json:"title,omitempty"     yaml:"title"`
	Screen   string `json:"screen,omitempty"    yaml:"screen"`
	Position []int  `json:"position,omitempty"  yaml:"position"`
	Size     []int  `json:"size,omitempty"      yaml:"size"`
}
