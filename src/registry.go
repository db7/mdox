package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type Registry struct {
	Style    Style
	Options  map[RegistryOption]bool
	_files   map[string]*CompoundDef
	_groups  map[string]*CompoundDef
	_page    map[string]*CompoundDef
	_dir     map[string]*CompoundDef
	_members map[string]*MemberWrapper
	srcPath  string
}

func newRegistry(srcPath string) *Registry {
	return &Registry{
		_files:   make(map[string]*CompoundDef),
		_page:    make(map[string]*CompoundDef),
		_groups:  make(map[string]*CompoundDef),
		_dir:     make(map[string]*CompoundDef),
		_members: make(map[string]*MemberWrapper),
		Options: map[RegistryOption]bool{
			ParaLine:   true,
			References: true,
		},
		srcPath: srcPath,
	}
}

func getCommonPrefix(p1, p2 string) string {
	if p1 == "" || p2 == "" {
		return ""
	}
	clean := func(p string) []string {
		p = filepath.ToSlash(filepath.Clean(p))
		if p == "." {
			return nil
		}
		p = strings.TrimPrefix(p, "/")
		if p == "" {
			return nil
		}
		return strings.Split(p, "/")
	}
	s1 := clean(p1)
	s2 := clean(p2)
	if len(s1) == 0 || len(s2) == 0 {
		return ""
	}
	limit := len(s1)
	if len(s2) < limit {
		limit = len(s2)
	}
	var common []string
	for i := 0; i < limit; i++ {
		if s1[i] != s2[i] {
			break
		}
		common = append(common, s1[i])
	}
	if len(common) == 0 {
		return ""
	}
	prefix := strings.Join(common, "/")
	if len(common) < len(s1) && len(common) < len(s2) {
		prefix += "/"
	}
	return prefix
}

func (r *Registry) Register(c CompoundDef) {
	log.Println("Registering", c.CompoundName, c.Id)
	switch c.Kind {
	case "file":
		r._files[c.Id] = &c
	case "group":
		r._groups[c.Id] = &c
		r._groups[c.Id].Location.File = fmt.Sprintf("%s.md", c.CompoundName)
	case "page":
		r._page[c.Id] = &c
	case "dir":
		r._dir[c.Id] = &c
	}

	for _, s := range c.SectionDef {
		for _, m := range s.MemberDef {
			mm := m
			log.Println("register", m.Name, "with id", m.Id)
			r._members[m.Id] = &mm
		}
	}
}
func (r *Registry) getFilePath(id string) (rel, abs string, has bool) {
	f, has := r._files[id]
	if !has {
		return "", "", has
	}
	return f.Location.File, "", true

}
func (r *Registry) file(id string) *CompoundDef {
	return r._files[id]
}
func (r *Registry) dir(id string) *CompoundDef {
	return r._dir[id]
}

type RegEntry struct {
	Location string
	Kind     Kind
	Name     string
}

func (r *Registry) get(ctx DumpContext, id string) *RegEntry {
	if e, has := r._files[id]; has {
		return &RegEntry{
			Location: e.getPath(ctx),
			Kind:     KindFile,
			Name:     id,
		}
	}
	if e, has := r._groups[id]; has {
		return &RegEntry{
			Location: e.getPath(ctx),
			Kind:     KindGroup,
			Name:     id,
		}
	}
	if e, has := r._page[id]; has {
		return &RegEntry{
			Location: e.getPath(ctx),
			Kind:     KindPage,
			Name:     id,
		}
	}
	if e, has := r._dir[id]; has {
		return &RegEntry{
			Location: e.getPath(ctx),
			Kind:     KindDir,
			Name:     id,
		}
	}
	if e, has := r._members[id]; has {
		var k Kind
		switch e.Kind {
		case "function":
			k = KindFunc
		case "define":
			k = KindMacro
		default:
			panic("oops")
		}
		return &RegEntry{
			Location: e.Location.File,
			Kind:     k,
			Name:     e.Name,
		}
	}
	return nil
}

func (r *Registry) search(id string) (string, error) {
	if e, has := r._files[id]; has {
		return e.Location.File, nil
	}
	if e, has := r._groups[id]; has {
		return e.Location.File, nil
	}
	if e, has := r._page[id]; has {
		return e.Location.File, nil
	}
	if e, has := r._dir[id]; has {
		path := filepath.Base(e.Location.File)
		return path + "/README.md", nil
	}

	if _, has := r._members[id]; has {
		panic("yes")
	}

	return id, nil
}
func (c *CompoundDef) hasInnerFile(id string) bool {
	for _, f := range c.InnerFile {
		if f.RefID == id {
			return true
		}
	}
	return false
}
func (r *Registry) getGroupsWith(id string) []string {
	var g []string
	for _, gg := range r._groups {
		if gg.hasInnerFile(id) {
			g = append(g, gg.Id)
		}
	}
	return g
}
func (r *Registry) groupName(id string) string {
	g, has := r._groups[id]
	if !has {
		panic("not has")
	}
	return g.Title
}

type RegistryOption int

const (
	ParaLine RegistryOption = iota
	PreserveNewLines
	References
)

func (r *Registry) SetOption(o RegistryOption, v bool) bool {
	cur := r.Options[o]
	r.Options[o] = v
	return cur
}

func (r *Registry) Disable(o RegistryOption) bool {
	return r.SetOption(o, false)
}
func (r *Registry) Enable(o RegistryOption) bool {
	return r.SetOption(o, true)
}
func (r *Registry) Option(o RegistryOption) bool {
	return r.Options[o]
}
