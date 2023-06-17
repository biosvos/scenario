package scenario

type Step interface {
	Title() string
	Input() any
	DefinitionOfDone() []string
	Run() (*Artifacts, error)
}
