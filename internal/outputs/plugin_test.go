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
			name: "stdout success",
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
			name: "http success",
			args: args{
				block: config.PluginBlock{
					Type: "http",
					Config: map[string]interface{}{
						"url":       "https://example.com/api",
						"method":    "PUT",
						"env_token": "TOKEN1234",
					},
				},
			},
			want: func(w OutputWriter) error {
				var writer = w.(*httpOutput)
				assert.NotEmpty(t, writer.opts.URL, "URL was empty")
				assert.NotEmpty(t, writer.opts.Method, "Method was empty")
				assert.NotEmpty(t, writer.opts.EnvToken, "EnvToken was empty")
				assert.NotNil(t, writer.client, "Client was nil")
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
			err = tt.want(got)
			assert.NoError(t, err, "FromPlugin(%v, %v)", tt.args.block, stdout)
			if err != nil || got == nil {
				return
			}
			_ = got.WriteOutput(map[string]interface{}{"test": "test"})
		})
	}
}
