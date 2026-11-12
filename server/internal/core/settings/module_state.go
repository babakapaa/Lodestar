package settings

type ModuleState struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Optional bool   `json:"optional"`
	Enabled  bool   `json:"enabled"`
}
