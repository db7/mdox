package main

import (
	"log"
	"path/filepath"
)

type Registry struct {
	Style      Style
	Options    map[RegistryOption]bool
	_files     map[string]*CompoundDef
	_groups    map[string]*CompoundDef
	_page      map[string]*CompoundDef
	_dir       map[string]*CompoundDef
	srcPath    string
	filePrefix string
	//translate map[string]string
	// 	funcs    map[string]*CompoundDef
	// 	macros   string
	// 	typedefs string
}

func newRegistry(srcPath string) *Registry {
	return &Registry{
		_files:  make(map[string]*CompoundDef),
		_page:   make(map[string]*CompoundDef),
		_groups: make(map[string]*CompoundDef),
		_dir:    make(map[string]*CompoundDef),
		Options: map[RegistryOption]bool{
			ParaLine:   true,
			References: true,
		},
		srcPath: srcPath,
	}
}

func getCommonPrefix(p1, p2 string) string {
	shortLen := len(p1)
	if l := len(p2); l < shortLen {
		shortLen = l
	}
	i := 0
	for ; i < shortLen && p1[i] == p2[i]; i++ {

	}
	return p1[:i]
}

func (r *Registry) Register(c CompoundDef) {
	log.Println("Registering", c.CompoundName, c.Id)
	switch c.Kind {
	case "file":
		dp := filepath.Dir(c.Location.File)
		if len(r._files) == 0 {
			r.filePrefix = dp
		}
		r.filePrefix = getCommonPrefix(r.filePrefix, dp)
		r._files[c.Id] = &c
	case "group":
		r._groups[c.Id] = &c
	case "page":
		r._page[c.Id] = &c
	case "dir":
		r._dir[c.Id] = &c
	}
}
func (r *Registry) getFilePath(id string) (rel, abs string, has bool) {
	f, has := r._files[id]
	if !has {
		return "", "", has
	}
	return f.Location.File, r.filePrefix, true

}
func (r *Registry) file(id string) *CompoundDef {
	return r._files[id]
}
func (r *Registry) dir(id string) *CompoundDef {
	return r._dir[id]
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
	return id, nil
	//return "", nil
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
