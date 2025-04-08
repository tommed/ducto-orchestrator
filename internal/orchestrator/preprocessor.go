package orchestrator

import "context"

// Preprocessor modifies or enriches input before transformation.
type Preprocessor interface {
	Process(ctx context.Context, input map[string]interface{}) error
}
