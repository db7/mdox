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

func Test_getCommonPrefix(t *testing.T) {
	testCases := []struct {
		name   string
		first  string
		second string
		want   string
	}{
		{
			name:   "noOverlap",
			first:  "alpha/beta",
			second: "beta/alpha",
			want:   "",
		},
		{
			name:   "sharedDirectory",
			first:  "foo/bar/baz.h",
			second: "foo/bar/qux.h",
			want:   "foo/bar/",
		},
		{
			name:   "sharedPrefixOnly",
			first:  "foo/bar-alpha/qux.h",
			second: "foo/bar-beta/baz.h",
			want:   "foo/",
		},
		{
			name:   "identicalPath",
			first:  "foo/bar",
			second: "foo/bar",
			want:   "foo/bar",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := getCommonPrefix(tc.first, tc.second)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_getRelativePath(t *testing.T) {
	testCases := []struct {
		name string
		from string
		to   string
		want string
	}{
		{
			name: "sameDirectory",
			from: "src/foo/bar.h",
			to:   "src/foo/baz.h",
			want: "baz.h",
		},
		{
			name: "goUpOneDirectory",
			from: "src/foo/bar/baz.h",
			to:   "src/foo/qux/quux.h",
			want: "../qux/quux.h",
		},
		{
			name: "completelyDifferentRoots",
			from: "alpha/beta/gamma.h",
			to:   "theta/lambda/kappa.h",
			want: "../../theta/lambda/kappa.h",
		},
		{
			name: "directoriesWithCommonPrefix",
			from: "foo/bar-alpha/a.h",
			to:   "foo/bar-beta/b.h",
			want: "../bar-beta/b.h",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := getRelativePath(tc.from, tc.to)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_shouldIgnore(t *testing.T) {
	testCases := []struct {
		name string
		path string
		want bool
	}{
		{"docHeader", "/proj/include/doc.h", true},
		{"internalDir", "/proj/internal/foo.h", true},
		{"commonDir", "/proj/common/foo.h", true},
		{"utilsDir", "/proj/foo/utils/bar.h", true},
		{"cdsDir", "/proj/foo/bar/cds/baz.h", true},
		{"regularFile", "/proj/src/foo.h", false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := shouldIgnore(tc.path)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_shouldIgnoreFile(t *testing.T) {
	testCases := []struct {
		name string
		path string
		want bool
	}{
		{"vtypesHeader", "include/vtypes.h", true},
		{"vstdintHeader", "include/vstdint.h", true},
		{"docHeader", "doc/doc.h", true},
		{"regularHeader", "include/foo.h", false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := shouldIgnoreFile(tc.path)
			assert.Equal(t, tc.want, got)
		})
	}
}
