package scenario

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewArtifacts(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ret, err := NewArtifacts(map[string]any{
			"a": "hello",
			"b": 1,
		})
		require.NoError(t, err)
		require.NotNil(t, ret)
	})
	t.Run("same key", func(t *testing.T) {
		ret, err := NewArtifacts(map[string]any{
			"a": "hello",
		}, map[string]any{
			"a": 1,
		})
		require.Error(t, err)
		require.Nil(t, ret)
	})
}

func TestArtifacts_Add(t *testing.T) {
	t.Run("", func(t *testing.T) {
		artifacts, _ := NewArtifacts()
		err := artifacts.Add("A", 1)
		require.NoError(t, err)
	})
	t.Run("zero value", func(t *testing.T) {
		artifacts, _ := NewArtifacts()
		err := artifacts.Add("a", 0)
		require.Error(t, err)
	})
	t.Run("false is not zero value", func(t *testing.T) {
		artifacts, _ := NewArtifacts()
		err := artifacts.Add("a", false)
		require.NoError(t, err)
	})
	t.Run("empty key", func(t *testing.T) {
		artifacts, _ := NewArtifacts()
		err := artifacts.Add("", 1)
		require.Error(t, err)
	})
	t.Run("same key", func(t *testing.T) {
		artifacts, _ := NewArtifacts()

		_ = artifacts.Add("A", 1)
		err := artifacts.Add("A", 1)

		require.Error(t, err)
	})
	t.Run("same key", func(t *testing.T) {
		artifacts, _ := NewArtifacts(map[string]any{
			"A": 1,
		})

		err := artifacts.Add("A", 1)

		require.Error(t, err)
	})
}

func TestArtifacts_Fill(t *testing.T) {
	type input struct {
		A int
	}
	t.Run("", func(t *testing.T) {
		artifacts, _ := NewArtifacts(map[string]any{
			"a": 1,
		})
		in := input{}

		err := artifacts.Fill(&in)

		require.NoError(t, err)
		require.NotZero(t, in.A)
	})
	t.Run("no artifacts", func(t *testing.T) {
		artifacts, _ := NewArtifacts()
		in := input{}

		err := artifacts.Fill(&in)

		require.Error(t, err)
		require.Zero(t, in.A)
	})
	t.Run("no artifacts", func(t *testing.T) {
		artifacts, _ := NewArtifacts()

		err := artifacts.Fill(nil)

		require.Error(t, err)
	})
	t.Run("not struct", func(t *testing.T) {
		artifacts, _ := NewArtifacts()
		var a int

		err := artifacts.Fill(&a)

		require.Error(t, err)
	})
}

func TestArtifacts_IsExists(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		artifacts, _ := NewArtifacts(map[string]any{
			"a": 1,
		})
		ret := artifacts.IsExists("a")
		require.True(t, ret)
	})
	t.Run("failed", func(t *testing.T) {
		artifacts, _ := NewArtifacts()
		ret := artifacts.IsExists("a")
		require.False(t, ret)
	})
}

func TestArtifacts_Merge(t *testing.T) {
	a, _ := NewArtifacts(map[string]any{
		"A": 1,
	})
	b, _ := NewArtifacts(map[string]any{
		"B": 2,
	})

	err := a.Merge(b)

	require.NoError(t, err)
	require.Len(t, a.artifacts, 2)
}
