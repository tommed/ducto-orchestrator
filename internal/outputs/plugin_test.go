package outputs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-orchestrator/internal/config"
	"testing"
)

func TestFromPlugin(t *testing.T) {
	type args struct {
		block config.PluginBlock
	}
	tests := []struct {
		name    string
		args    args
		want    func(OutputWriter) error
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				block: config.PluginBlock{
					Type: "stdout",
					Config: map[string]interface{}{
						"pretty": true,
					},
				},
			},
			want: func(w OutputWriter) error {
				if w.(*stdoutWriter).opts.Pretty != true {
					return errors.New("pretty not applied")
				}
				return nil
			},
			wantErr: assert.NoError,
		},
		{
			name: "bad decoding",
			args: args{
				block: config.PluginBlock{
					Type: "stdout",
					Config: map[string]interface{}{
						"pretty": []string{"type", "mismatch"},
					},
				},
			},
			want: func(w OutputWriter) error {
				return nil
			},
			wantErr: assert.Error,
		},
		{
			name: "unsupported type",
			args: args{
				block: config.PluginBlock{
					Type: "invalid-type",
				},
			},
			want: func(w OutputWriter) error {
				return nil
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			got, err := FromPlugin(tt.args.block, stdout)
			if !tt.wantErr(t, err, fmt.Sprintf("FromPlugin(%v, %v)", tt.args.block, stdout)) {
				return
			}
			assert.NoError(t, tt.want(got), "FromPlugin(%v, %v)", tt.args.block, stdout)
		})
	}
}
