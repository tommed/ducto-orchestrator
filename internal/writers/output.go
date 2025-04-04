package writers

type OutputWriter interface {
	WriteOutput(map[string]interface{}) error
}

type FakeWriter struct {
	Written map[string]interface{}
}

func (f *FakeWriter) WriteOutput(data map[string]interface{}) error {
	f.Written = data
	return nil
}
