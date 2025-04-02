package orchestrator

type InputReader interface {
	ReadInput() (map[string]interface{}, error)
}
