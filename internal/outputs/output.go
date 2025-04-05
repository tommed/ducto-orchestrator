package outputs

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

func NewFailingWriter(err error) OutputWriter {
	return &failingWriter{err: err}
}

type failingWriter struct {
	err error
}

func (f *failingWriter) WriteOutput(data map[string]interface{}) error {
	return f.err
}
