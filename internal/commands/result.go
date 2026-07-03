package commands

// Result mirrors Double Commander TCommandFuncResult.
type Result int

const (
	ResultSuccess Result = iota
	ResultDisabled
	ResultNotFound
)
