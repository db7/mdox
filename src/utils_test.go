package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getDepth(t *testing.T) {
	testCases := []struct {
		str_path string
		r        int
	}{
		{"vsync/", 1},
		{"vsync/queue/", 2},
		{"vsync/queue", 2},
	}

	for _, tc := range testCases {
		v := getDepth(tc.str_path)
		assert.Equal(t, v, tc.r)
	}
}
