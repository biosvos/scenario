package scenario

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var _ Step = &GameStep{}

type GameInput struct {
	Title string `json:"title"`
}

type GameStep struct {
	input GameInput
}

func (g *GameStep) Run() (*Artifacts, error) {
	artifacts, _ := NewArtifacts(map[string]any{
		"time": 1,
	})
	return artifacts, nil
}

func (g *GameStep) Title() string {
	return "game"
}

func (g *GameStep) DefinitionOfDone() []string {
	return []string{"time"}
}

func (g *GameStep) Input() any {
	return &g.input
}

var _ Step = &NextStep{}

type NextInput struct {
	Time int `json:"time"`
}

type NextStep struct {
	input NextInput
}

func (n *NextStep) DefinitionOfDone() []string {
	return []string{"last"}
}

func (n *NextStep) Input() any {
	return &n.input
}

func (n *NextStep) Run() (*Artifacts, error) {
	artifacts, _ := NewArtifacts(map[string]any{
		"last": 1,
	})
	return artifacts, nil
}

func (n *NextStep) Title() string {
	return "next"
}

func TestNewScenario(t *testing.T) {
	ret, err := NewScenario([]Step{
		&GameStep{},
	}, WithArtifacts(map[string]any{}), WithArtifact("a", 2))
	require.NotNil(t, ret)
	require.NoError(t, err)
}

func TestScenario_Run(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		scenario, _ := NewScenario([]Step{
			&GameStep{},
		},
			WithArtifact("title", "game step"),
		)

		err := scenario.Run(os.Stdout)

		require.NoError(t, err)
	})
	t.Run("failed", func(t *testing.T) {
		scenario, _ := NewScenario([]Step{
			&GameStep{},
		},
			WithArtifact("time", 1),
		)

		err := scenario.Run(nil)

		require.Error(t, err)
	})
	t.Run("", func(t *testing.T) {
		scenario, _ := NewScenario([]Step{
			&GameStep{},
			&NextStep{},
		},
			WithArtifact("title", "game step"),
		)

		err := scenario.Run(os.Stdout)

		require.NoError(t, err)
	})
}
