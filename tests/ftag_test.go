package tests

import (
	"errors"
	"testing"

	"github.com/Southclaws/fault/ftag"
	"github.com/stretchr/testify/assert"
)

func TestWrapWithKind(t *testing.T) {
	err := ftag.Wrap(errors.New("a problem"), ftag.NotFound)
	out := ftag.Get(err)

	assert.Equal(t, ftag.NotFound, out)
}

func TestWrapWithKindChanging(t *testing.T) {
	err := ftag.Wrap(errors.New("a problem"), ftag.Internal)
	err = ftag.Wrap(err, ftag.Internal)
	err = ftag.Wrap(err, ftag.Internal)
	err = ftag.Wrap(err, ftag.InvalidArgument)
	err = ftag.Wrap(err, ftag.InvalidArgument)
	err = ftag.Wrap(err, ftag.NotFound)
	out := ftag.Get(err)

	assert.Equal(t, ftag.NotFound, out, "Should always pick the most recent kind from an error chain.")
}

func TestMultipleWrappedKind(t *testing.T) {
	err := ftag.Wrap(errors.New("a problem"), ftag.Internal)
	err = ftag.Wrap(err, ftag.InvalidArgument)
	err = ftag.Wrap(err, ftag.NotFound)
	out := ftag.GetAll(err)

	assert.Equal(t, []ftag.Kind{ftag.NotFound, ftag.InvalidArgument, ftag.Internal}, out)
}

func TestIsKind(t *testing.T) {
	err := ftag.Wrap(errors.New("a problem"), ftag.Internal)
	err = ftag.Wrap(err, ftag.InvalidArgument)

	assert.True(t, ftag.Is(err, ftag.InvalidArgument), "Should return true for InvalidArgument kind")
	assert.False(t, ftag.Is(err, ftag.NotFound), "Should return false for NotFound kind")
	assert.True(t, ftag.Is(err, ftag.Internal), "Should return True for Internal kind")
}
