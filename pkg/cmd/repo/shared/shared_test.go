package shared

import (
	"errors"
	"net/http"
	"testing"

	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/config"
	"github.com/cli/cli/pkg/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewRepo(t *testing.T) {
	type input struct {
		s      string
		config func() (config.Config, error)
		client *api.Client
	}

	tests := []struct {
		name      string
		input     input
		wantName  string
		wantOwner string
		wantHost  string
		wantErr   bool
	}{
		{
			name: "config returns error",
			input: input{s: "REPO",
				config: func() (config.Config, error) { return nil, errors.New("error") },
				client: successClient()},
			wantErr: true,
		},
		{
			name: "client returns error",
			input: input{s: "REPO",
				config: blankConfigFunc(),
				client: errorClient()},
			wantErr: true,
		},
		{
			name: "config is nil",
			input: input{s: "REPO",
				config: nil,
				client: successClient()},
			wantName:  "REPO",
			wantOwner: "OWNER",
			wantHost:  "github.com",
		},
		{
			name: "client is nil",
			input: input{s: "REPO",
				config: blankConfigFunc(),
				client: nil,
			},
			wantName: "REPO",
			wantHost: "github.com",
		},
		{
			name: "REPO returns proper values",
			input: input{s: "REPO",
				config: blankConfigFunc(),
				client: successClient()},
			wantName:  "REPO",
			wantOwner: "OWNER",
			wantHost:  "github.com",
		},
		{
			name: "SOMEONE/REPO returns proper values",
			input: input{s: "SOMEONE/REPO",
				config: blankConfigFunc(),
				client: successClient()},
			wantName:  "REPO",
			wantOwner: "SOMEONE",
			wantHost:  "github.com",
		},
		{
			name: "HOST/SOMEONE/REPO returns proper values",
			input: input{s: "HOST/SOMEONE/REPO",
				config: blankConfigFunc(),
				client: successClient()},
			wantName:  "REPO",
			wantOwner: "SOMEONE",
			wantHost:  "host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRepo(tt.input.s, tt.input.config, tt.input.client)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantName, r.RepoName())
			assert.Equal(t, tt.wantOwner, r.RepoOwner())
			assert.Equal(t, tt.wantHost, r.RepoHost())
		})
	}
}

func blankConfigFunc() func() (config.Config, error) {
	return func() (config.Config, error) {
		return config.NewBlankConfig(), nil
	}
}

func errorClient() *api.Client {
	reg := &httpmock.Registry{}
	reg.Register(
		httpmock.GraphQL(`query UserCurrent`),
		httpmock.StatusStringResponse(404, "not found"),
	)
	httpClient := &http.Client{Transport: reg}
	apiClient := api.NewClientFromHTTP(httpClient)
	return apiClient
}

func successClient() *api.Client {
	reg := &httpmock.Registry{}
	reg.Register(
		httpmock.GraphQL(`query UserCurrent`),
		httpmock.StringResponse(`{"data":{"viewer":{"login":"OWNER"}}}`),
	)
	httpClient := &http.Client{Transport: reg}
	apiClient := api.NewClientFromHTTP(httpClient)
	return apiClient
}
