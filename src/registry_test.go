package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryRegisterAndGetFileEntry(t *testing.T) {
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

	entry := reg.get(DumpContext{Reg: reg}, "file_1")
	if !assert.NotNil(t, entry) {
		return
	}
	assert.Equal(t, KindFile, entry.Kind)
	assert.Equal(t, "include/foo.h", entry.Location)
}

func TestRegistryRegistersMembers(t *testing.T) {
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
						Name:     "FOO",
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

	ctx := DumpContext{Reg: reg}

	macro := reg.get(ctx, "macro_1")
	if !assert.NotNil(t, macro) {
		return
	}
	assert.Equal(t, KindMacro, macro.Kind)
	assert.Equal(t, "include/foo.h", macro.Location)

	function := reg.get(ctx, "func_1")
	if !assert.NotNil(t, function) {
		return
	}
	assert.Equal(t, KindFunc, function.Kind)
	assert.Equal(t, "include/foo.h", function.Location)
}

func TestRegistrySearchAndGroups(t *testing.T) {
	reg := newRegistry("src")
	file := CompoundDef{
		Id:           "file_1",
		Kind:         "file",
		CompoundName: "foo.h",
		Location: Location{
			File: "src/foo.h",
		},
	}
	reg.Register(file)

	dir := CompoundDef{
		Id:           "dir_1",
		Kind:         "dir",
		CompoundName: "src/core",
		Location: Location{
			File: "src/core",
		},
	}
	reg.Register(dir)

	group := CompoundDef{
		Id:           "group_1",
		Kind:         "group",
		CompoundName: "core",
		Location: Location{
			File: "unused",
		},
		InnerFile: []InnerFile{
			{RefID: "file_1"},
		},
	}
	reg.Register(group)

	path, err := reg.search("dir_1")
	if assert.NoError(t, err) {
		assert.Equal(t, "core/README.md", path)
	}

	groups := reg.getGroupsWith("file_1")
	assert.Equal(t, []string{"group_1"}, groups)
}
