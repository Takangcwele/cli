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
		errMsg   string
		stdin    string
	}{
		{
			name:     "blank nontty",
			wantsErr: true,
			errMsg:   "workflow ID, name, or filename required when not running interactively",
		},
		{
			name: "blank tty",
			tty:  true,
			wants: RunOptions{
				Prompt: true,
			},
		},
		{
			name: "ref flag",
			tty:  true,
			cli:  "--ref 12345abc",
			wants: RunOptions{
				Prompt: true,
				Ref:    "12345abc",
			},
		},
		{
			name: "extra args",
			tty:  true,
			cli:  `workflow.yml -- --cool=nah --foo bar`,
			wants: RunOptions{
				InputArgs: []string{"--cool=nah", "--foo", "bar"},
				Selector:  "workflow.yml",
			},
		},
		{
			name:     "both json on STDIN and json arg",
			cli:      `workflow.yml --json '{"cool":"yeah"}'`,
			stdin:    `{"cool":"yeah"}`,
			wantsErr: true,
			errMsg:   "JSON can only be passed on one of STDIN or --json at a time",
		},
		{
			name:     "both json on STDIN and extra args",
			cli:      `workflow.yml -- --cool=nah`,
			stdin:    `{"cool":"yeah"}`,
			errMsg:   "only one of JSON or input arguments can be passed at a time",
			wantsErr: true,
		},
		{
			name:     "both json arg and extra args",
			tty:      true,
			cli:      `workflow.yml --json '{"cool":"yeah"}' -- --cool=nah`,
			errMsg:   "only one of JSON or input arguments can be passed at a time",
			wantsErr: true,
		},
		{
			name: "json via argument",
			cli:  `workflow.yml --json '{"cool":"yeah"}'`,
			tty:  true,
			wants: RunOptions{
				JSON:     `{"cool":"yeah"}`,
				Selector: "workflow.yml",
			},
		},
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
			}
			io.SetStdoutTTY(tt.tty)

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
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
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
