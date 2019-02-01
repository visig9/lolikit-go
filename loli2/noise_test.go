package loli2

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/visig9/lolikit-go/loli2/simplerepo"
)

func TestIsNoise(t *testing.T) {
	sr := simplerepo.New()
	defer sr.Close()

	for _, path := range sr.XPaths {
		assert.Equal(
			t, true, IsNoise(path),
			"IsNoise() on noise path",
		)
	}

	var notNoises []string
	notNoises = append(notNoises, sr.SNPaths...)
	notNoises = append(notNoises, sr.CNPaths...)
	notNoises = append(notNoises, sr.DPaths...)

	for _, path := range notNoises {
		assert.Equal(
			t, false, IsNoise(path),
			"IsNoise() on not-noise path",
		)
	}
}
