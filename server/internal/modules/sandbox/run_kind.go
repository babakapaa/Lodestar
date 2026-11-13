package sandbox

const (
	RunJupyter = "jupyter"
	RunVSCode  = "vscode"
	RunJob     = "job"
	RunDemo    = "demo"
)

func commanded(kind string) bool {
	return kind == RunJob || kind == RunDemo
}

func interactive(kind string) bool {
	return kind == RunJupyter || kind == RunVSCode
}
