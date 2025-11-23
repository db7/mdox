package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefDumpForFiles(t *testing.T) {
	reg := newRegistry("src")
	file := CompoundDef{
		Id:           "file_1",
		Kind:         "file",
		CompoundName: "foo.h",
		Location: Location{
			File: "include/foo.h",
		},
	}
	reg.Register(file)

	ctx := DumpContext{
		Reg:  reg,
		Path: "include/foo.h",
	}

	ref := newRef("Foo Header", "file_1")
	w := new(Writer)
	err := ref.Dump(ctx, w)
	if assert.NoError(t, err) {
		assert.Equal(t, "[Foo Header](foo.h.md)", strings.Join(w.data, ""))
	}
}

func TestRefDumpForMembersAnchors(t *testing.T) {
	reg := newRegistry("src")
	compound := CompoundDef{
		Id:           "file_1",
		Kind:         "file",
		CompoundName: "foo.h",
		Location: Location{
			File: "include/foo.h",
		},
		SectionDef: []SectionDef{
			{
				MemberDef: []MemberWrapper{
					{MemberDef: MemberDef{
						Id:       "macro_1",
						Kind:     "define",
						Name:     "FOO_VALUE",
						Location: Location{File: "include/foo.h"},
					}},
					{MemberDef: MemberDef{
						Id:       "func_1",
						Kind:     "function",
						Name:     "foo_init",
						Location: Location{File: "include/foo.h"},
					}},
				},
			},
		},
	}
	reg.Register(compound)

	ctx := DumpContext{
		Reg:  reg,
		Path: "include/bar.h",
	}

	w := new(Writer)
	if assert.NoError(t, newRef("FOO_VALUE", "macro_1").Dump(ctx, w)) {
		assert.Equal(t, "[FOO_VALUE](foo.h.md#macro-foo_value)", strings.Join(w.data, ""))
	}

	w = new(Writer)
	if assert.NoError(t, newRef("foo_init", "func_1").Dump(ctx, w)) {
		assert.Equal(t, "[foo_init](foo.h.md#function-foo_init)", strings.Join(w.data, ""))
	}
}

func TestRefDumpFallsBackToRawRef(t *testing.T) {
	reg := newRegistry("src")
	ctx := DumpContext{
		Reg:  reg,
		Path: "include/foo.h",
	}

	ref := newRef("Unknown", "missing_id")
	w := new(Writer)
	err := ref.Dump(ctx, w)
	if assert.NoError(t, err) {
		assert.Equal(t, "[Unknown](missing_id)", strings.Join(w.data, ""))
	}
}
