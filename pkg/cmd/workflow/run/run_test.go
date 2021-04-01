package run

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdRun(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		tty      bool
		wants    RunOptions
		wantsErr bool
		stdin    string
	}{
		// TODO blank nontty
		// TODO blank tty
		// TODO workflow arg
		// TODO ref
		// TODO json as --json
		// TODO extra args
		// TODO both json on stdin and extra args
		// TODO both json arg and extra args
		{
			name:  "json on STDIN",
			cli:   "workflow.yml",
			stdin: `{"cool":"yeah"}`,
			wants: RunOptions{
				JSON:     `{"cool":"yeah"}`,
				Selector: "workflow.yml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			io, stdin, _, _ := iostreams.Test()
			if tt.stdin == "" {
				io.SetStdinTTY(tt.tty)
			} else {
				stdin.WriteString(tt.stdin)
				io.SetStdoutTTY(tt.tty)
			}

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *RunOptions
			cmd := NewCmdRun(f, func(opts *RunOptions) error {
				gotOpts = opts
				return nil
			})
			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(ioutil.Discard)
			cmd.SetErr(ioutil.Discard)

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tt.wants.Selector, gotOpts.Selector)
			assert.Equal(t, tt.wants.Prompt, gotOpts.Prompt)
			assert.Equal(t, tt.wants.JSON, gotOpts.JSON)
			assert.Equal(t, tt.wants.Ref, gotOpts.Ref)
			assert.ElementsMatch(t, tt.wants.InputArgs, gotOpts.InputArgs)
		})
	}
}

// TODO execution tests
