package github

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-github/v26/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChecksService struct {
	mock.Mock
}

func (m *mockChecksService) ListCheckRunsForRef(context context.Context, owner string, repo string, ref string, opt *github.ListCheckRunsOptions) (*github.ListCheckRunsResults, *github.Response, error) {
	args := m.Called()
	results := args.Get(0)

	if results != nil {
		return args.Get(0).(*github.ListCheckRunsResults), nil, nil
	}

	return nil, nil, args.Error(2)

}

func TestGetTravisCIBuildURLForRef_Success(t *testing.T) {
	ctx := context.Background()
	g := New(ctx, "Wattpad", "fake-access-token")

	m := new(mockChecksService)
	m.On("ListCheckRunsForRef").Return(&github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{
			{
				DetailsURL: github.String("https://travis-ci.com/Wattpad/highlander/builds/115827260"),
				App: &github.App{
					Name: github.String("Travis CI"),
				},
			},
		},
	}, nil, nil)
	g.Checks = m

	url, _ := g.GetTravisCIBuildURLForRef(ctx, "highlander", "master")

	assert.Equal(t, "https://travis-ci.com/Wattpad/highlander/builds/115827260", url)
}

func TestGetTravisCIBuildURLForRef_NoTravisCI(t *testing.T) {
	ctx := context.Background()
	g := New(ctx, "Wattpad", "fake-access-token")

	m := new(mockChecksService)
	m.On("ListCheckRunsForRef").Return(&github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{
			{
				DetailsURL: github.String("https://codeship.com/build/12345678"),
				App: &github.App{
					Name: github.String("Codeship"),
				},
			},
		},
	}, nil, nil)
	g.Checks = m

	_, err := g.GetTravisCIBuildURLForRef(ctx, "highlander", "master")

	assert.EqualError(t, errTravisCIBuildNotFound, err.Error())
}

func TestGetTravisCIBuildURLForRef_Empty(t *testing.T) {
	ctx := context.Background()
	g := New(ctx, "Wattpad", "fake-access-token")

	m := new(mockChecksService)
	m.On("ListCheckRunsForRef").Return(&github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{},
	}, nil, nil)
	g.Checks = m

	_, err := g.GetTravisCIBuildURLForRef(ctx, "highlander", "master")

	assert.EqualError(t, errTravisCIBuildNotFound, err.Error())
}

func TestGetTravisCIBuildURLForRef_Error(t *testing.T) {
	ctx := context.Background()
	g := New(ctx, "Wattpad", "fake-access-token")

	m := new(mockChecksService)
	m.On("ListCheckRunsForRef").Return(nil, nil, errors.New("fake"))
	g.Checks = m

	_, err := g.GetTravisCIBuildURLForRef(ctx, "highlander", "master")

	assert.Error(t, err)
}
