package sources

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStdinEventSource_Start(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name         string
		args         args
		wantErrStart error
		wantErrClose error
	}{
		{
			name: "simplest",
			args: args{
				fileName: "../../testdata/cases/simplest.input.json",
			},
		},
		{
			name: "invalid json",
			args: args{
				fileName: "../../testdata/invalid.json",
			},
			wantErrStart: errors.New("failed to decode stdin input: invalid character '}' looking for beginning of value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.ReadFile(tt.args.fileName)
			require.NoError(t, err)
			source := NewStdinEventSource(bytes.NewReader(input))
			_, err = source.Start(context.Background())
			if tt.wantErrStart != nil {
				assert.Equal(t, tt.wantErrStart.Error(), err.Error())
			}

			err = source.Close()
			if tt.wantErrClose != nil {
				assert.Equal(t, tt.wantErrClose.Error(), err.Error())
			}
		})
	}
}

func TestStdinEventSource_Close_Failing(t *testing.T) {
	source := NewStdinEventSource(&failingCloser{
		child: bytes.NewReader([]byte("")),
		err:   errors.New("intentional error"),
	})
	err := source.Close()
	assert.Error(t, err)
	assert.Equal(t, "intentional error", err.Error())
}

type failingCloser struct {
	child io.Reader
	err   error
}

func (f *failingCloser) Read(p []byte) (n int, err error) {
	return f.child.Read(p)
}

func (f *failingCloser) Close() error {
	return f.err
}
