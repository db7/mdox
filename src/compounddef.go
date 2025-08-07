package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"path/filepath"
	"sort"
)

type CompoundDef struct {
	XMLName      xml.Name     `xml:"compounddef"`
	SectionDef   []SectionDef `xml:"sectiondef"`
	Brief        Description  `xml:"briefdescription"`
	Detailed     Description  `xml:"detaileddescription"`
	Id           string       `xml:"id,attr"`
	Kind         string       `xml:"kind,attr"`
	Language     string       `xml:"language,attr"`
	InnerFile    []InnerFile  `xml:"innerfile"`
	InnerDir     []InnerDir   `xml:"innerdir"`
	InnerGroup   []InnerGroup `xml:"innergroup"`
	CompoundName string       `xml:"compoundname"`
	Location     Location     `xml:"location"`
	Title        string       `xml:"title"`
}

func (c *CompoundDef) Dump(ctx DumpContext, w *Writer) (err error) {
	addGroups := false
	switch c.Kind {
	case "page":
	case "file":
		addGroups = true
	case "group":
		addGroups = true
	case "dir":
		addGroups = true
	default:
	}

	c.dumpContent(ctx, w)
	c.dumpInnerFiles(ctx, w, addGroups)
	c.dumpInnerDirs(ctx, w)
	c.dumpSubgroups(ctx, w)
	return
}

func (c *CompoundDef) getPath(ctx DumpContext) string {
	if c.Kind == "group" {
		var commonDir string
		// find deepest common dir
		for _, f := range c.InnerFile {
			rel, _, has := ctx.Reg.getFilePath(f.RefID)
			if !has {
				continue
			}
			dir := filepath.Dir(rel)
			if commonDir == "" {
				commonDir = dir
			} else {
				commonDir = getCommonPrefix(commonDir, dir)
			}
		}
		if d, _ := filepath.Split(c.Location.File); d != "" {
			log.Fatal(c.Location.File, "has a directory")
		}

		return filepath.Join(commonDir, "GROUP_"+c.Location.File)
	}
	return c.Location.File
}

// dumpContent writes content of compound to w.
func (c *CompoundDef) dumpContent(ctx DumpContext, w *Writer) {
	// Prepare title with links to parent directories
	var ee Dumper
	if c.Title != "" {
		ee = newText(c.Title)
	} else {
		e := elem(c.Location.File)
		ee = &e
	}

	// Dump page title
	w.Print("#  ")
	ee.Dump(ctx, w)
	w.Println()

	// Dump descriptions
	ctx.Reg.Style = SEmphasis
	c.Brief.Dump(ctx, w)
	ctx.Reg.Style = Default

	if groups := ctx.Reg.getGroupsWith(c.Id); len(groups) > 0 {
		log.Println(groups)
		w.Print("**Groups:** ")
		sort.Strings(groups)
		for i, g := range groups {
			if i > 0 {
				w.Print(", ")
			}
			gg := ctx.Reg._groups[g]
			newRef(gg.Title, g).Dump(ctx, w)
		}
		w.Println()
		w.Println()
	}

	//w.Println("```c")
	//w.Printf("#include <%s>\n", c.Location.File)
	//w.Println("```")
	//w.Println("# Description ")
	c.Detailed.Dump(ctx, w)

	// Dump the sections (functions, macros, etc)
	for _, s := range c.SectionDef {
		s.Dump(ctx, w)
	}
}

// dumpInnerFiles writes a table of innerFiles.
func (c *CompoundDef) dumpInnerFiles(ctx DumpContext, w *Writer, addGroups bool) {
	reg := ctx.Reg
	var gorder []string
	groups := make(map[string]bool)
	fileInGroups := make(map[string]bool)
	if addGroups {
		for _, f := range c.InnerFile {
			for _, g := range ctx.Reg.getGroupsWith(f.RefID) {
				if _, has := groups[g]; !has {
					gorder = append(gorder, g)
					groups[g] = true
				}
				fileInGroups[fmt.Sprintf("%s--%s", f.RefID, g)] = true
			}
		}
	}
	tab := Table{
		Cols: 2 + len(groups),
		Row: []Row{{
			Entry: []Entry{
				newEntry(newText("File")),
				newEntry(newText("Description")),
			}},
		},
	}
	sort.Strings(gorder)
	for _, g := range gorder {
		e := newEntry(newText(reg.groupName(g)))
		tab.Row[0].Entry = append(tab.Row[0].Entry, e)
	}
	for _, f := range c.InnerFile {
		if shouldIgnoreFile(f.Name) {
			continue
		}
		ff := reg.file(f.RefID)
		ensureNotNil(ff)

		row := Row{
			Entry: []Entry{
				newEntry(newRef(ff.Location.File, f.RefID)),
				newEntry(&ff.Brief),
			},
		}
		for _, g := range gorder {
			log.Println(g)
			var e Entry
			if fileInGroups[fmt.Sprintf("%s--%s", f.RefID, g)] {
				// this renders on everything but on code-hub looks too small
				// similar to check mark
				e = newEntry(newText(" &#x2714; "))
			} else {
				// note :x: will not render on git code
				// so we use this instead
				// https://www.compart.com/de/unicode/U+274C
				e = newEntry(newText(" &#x274C; "))
			}
			row.Entry = append(row.Entry, e)
		}
		tab.Row = append(tab.Row, row)
	}
	if len(tab.Row) > 1 {
		w.Println("---")
		w.Println("## File Index")
		tab.Dump(ctx, w)
	}
}

// dumpInnerDir writes a table of innerDirs.
func (c *CompoundDef) dumpInnerDirs(ctx DumpContext, w *Writer) {
	tab := Table{
		Cols: 2,
		Row: []Row{{
			Entry: []Entry{
				newEntry(newText("Directory")),
				newEntry(newText("Description")),
			}},
		},
	}

	for _, f := range c.InnerDir {
		if shouldIgnore(f.Name) {
			continue
		}

		// Enforce references exist
		ff := ctx.Reg.dir(f.RefID)
		ensureNotNil(ff)

		row := Row{
			Entry: []Entry{
				newEntry(newRef(f.Name, f.RefID)),
				newEntry(&ff.Brief),
			},
		}
		tab.Row = append(tab.Row, row)
	}

	if len(tab.Row) > 1 {
		w.Println("---")
		w.Println("## Directory Index")
		tab.Dump(ctx, w)
	}
}

// dumpSubgroups writes a table with subgroups if available.
func (c *CompoundDef) dumpSubgroups(ctx DumpContext, w *Writer) {
	if len(c.InnerGroup) > 0 {
		w.Println("## Submodules")
		w.Println()
		for _, g := range c.InnerGroup {
			w.Print("- ")
			newRef(g.Name, g.RefID).Dump(ctx, w)
			w.Println()
		}
		w.Println()
	}
}

// dumpSubpages writes a table with subpages if available.
func (c *CompoundDef) dumpSubpages(ctx DumpContext, w *Writer) {
	log.Fatal("not implemented")
}
