package outputs

type OutputWriter interface {
	WriteOutput(map[string]interface{}) error
}
