package orchestrator

type InputReader interface {
	ReadInput() (map[string]interface{}, error)
}

type FakeReader struct {
	Data map[string]interface{}
}

func (f *FakeReader) ReadInput() (map[string]interface{}, error) {
	return f.Data, nil
}
