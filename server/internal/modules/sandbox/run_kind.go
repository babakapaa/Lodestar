package sandbox

const (
	RunJupyter = "jupyter"
	RunVSCode  = "vscode"
	RunJob     = "job"
)

func interactive(kind string) bool {
	return kind == RunJupyter || kind == RunVSCode
}
