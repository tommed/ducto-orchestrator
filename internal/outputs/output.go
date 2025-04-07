package outputs

import "context"

type OutputWriter interface {
	WriteOutput(context.Context, map[string]interface{}) error
}
