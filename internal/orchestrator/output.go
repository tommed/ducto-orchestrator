package orchestrator

type OutputWriter interface {
	WriteOutput(map[string]interface{}) error
}
