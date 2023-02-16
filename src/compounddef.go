package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"path/filepath"
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

func (c *CompoundDef) Dump(fd io.Writer, reg *Registry) (err error) {
	switch c.Kind {
	case "page":
		err = c.DumpPage(fd, reg)
	case "file":
		err = c.DumpFile(fd, reg)
	case "groupx":
		err = c.DumpGroup(fd, reg)
	case "dir":
		err = c.DumpDir(fd, reg)
		if err != nil {
			panic(err)
		}
	default:
	}
	return
}

func (c *CompoundDef) DumpHeader(fd io.Writer, reg *Registry) error {
	return nil
}

func (c *CompoundDef) DumpFile(fd io.Writer, reg *Registry) error {
	//fmt.Fprint(fd, `<a id="top"></a>`)
	e := elem(c.Location.File)
	ee := &e
	fmt.Fprintf(fd, "#  ")
	ee.Dump(fd, reg)
	fmt.Fprintf(fd, "\n\n")
	if err := c.Brief.Dump(fd, reg); err != nil {
		return err
	}
	if err := c.Detailed.Dump(fd, reg); err != nil {
		return err
	}
	for _, s := range c.SectionDef {
		if err := s.Dump(fd, reg); err != nil {
			return err
		}
	}
	return nil
}

func (c *CompoundDef) DumpDir(fd io.Writer, reg *Registry) error {
	//fmt.Fprint(fd, `<a id="top"></a>`)
	e := elem(c.Location.File)
	ee := &e
	fmt.Fprintf(fd, "#  ")
	ee.Dump(fd, reg)
	fmt.Fprintf(fd, "\n\n")

	if err := c.Brief.Dump(fd, reg); err != nil {
		return err
	}
	if err := c.Detailed.Dump(fd, reg); err != nil {
		return err
	}
	var gorder []string
	groups := make(map[string]bool)
	fileInGroups := make(map[string]bool)
	for _, f := range c.InnerFile {
		for _, g := range reg.getGroupsWith(f.RefID) {
			if _, has := groups[g]; !has {
				gorder = append(gorder, g)
				groups[g] = true
			}
			fileInGroups[fmt.Sprintf("%s--%s", f.RefID, g)] = true
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
	for _, g := range gorder {
		e := newEntry(newText(reg.groupName(g)))
		tab.Row[0].Entry = append(tab.Row[0].Entry, e)
	}
	for _, f := range c.InnerFile {
		ff := reg.file(f.RefID)
		target := filepath.Base(ff.Location.File) + ".md"
		row := Row{
			Entry: []Entry{
				newEntry(newRef(f.Name, target)),
				newEntry(&ff.Brief),
			},
		}
		for _, g := range gorder {
			var e Entry
			if fileInGroups[fmt.Sprintf("%s--%s", f.RefID, g)] {
				//e = newEntry(newText(" :heavy_check_mark: "))
				e = newEntry(newText(" + "))
			} else {
				e = newEntry(newText(" - "))
			}
			row.Entry = append(row.Entry, e)
		}
		tab.Row = append(tab.Row, row)
	}
	tab.Dump(fd, reg)

	tab = Table{
		Cols: 2,
		Row: []Row{{
			Entry: []Entry{
				newEntry(newText("Directory")),
				newEntry(newText("Description")),
			}},
		},
	}

	for _, f := range c.InnerDir {
		ff := reg.dir(f.RefID)
		if ff == nil {
			log.Fatal(f.RefID, "is nil")
		}
		row := Row{
			Entry: []Entry{
				newEntry(newRef(f.Name, f.RefID)),
				newEntry(&ff.Brief),
			},
		}
		tab.Row = append(tab.Row, row)
	}

	return tab.Dump(fd, reg)
}

func (c *CompoundDef) DumpGroup(fd io.Writer, reg *Registry) error {
	//fmt.Fprint(fd, `<a id="top"></a>`)
	fmt.Fprintf(fd, "#  GROUP %s\n\n", c.Title)
	if err := c.Brief.Dump(fd, reg); err != nil {
		return err
	}
	if err := c.Detailed.Dump(fd, reg); err != nil {
		return err
	}
	if len(c.InnerGroup) > 0 {
		fmt.Fprintln(fd, "## Submodules")
		fmt.Fprintln(fd)
		for _, g := range c.InnerGroup {
			fmt.Fprint(fd, "- ")
			err := newRef(g.Name, g.RefID).Dump(fd, reg)
			if err != nil {
				return err
			}
			fmt.Fprintln(fd)
		}
		fmt.Fprintln(fd)
	}
	if len(c.InnerFile) > 0 {
		fmt.Fprintln(fd, "## Files")
		fmt.Fprintln(fd)
		tab := Table{
			Cols: 2,
			Rows: len(c.InnerFile) + 1,
		}
		row := Row{
			Entry: []Entry{
				newEntry(newText("File")),
				newEntry(newText("Description")),
			},
		}
		tab.Row = append(tab.Row, row)
		for _, g := range c.InnerFile {
			f := reg.file(g.RefID)
			if f != nil {
				return fmt.Errorf("not found: %v", g.RefID)
			}
			row := Row{
				Entry: []Entry{
					newEntry(newRef(f.Location.File, g.RefID)),
					newEntry(&f.Brief),
				},
			}
			tab.Row = append(tab.Row, row)

			fmt.Fprintln(fd)
		}
		tab.Dump(fd, reg)
	}
	return nil
}

func (c *CompoundDef) DumpPage(fd io.Writer, reg *Registry) error {
	//fmt.Fprint(fd, `<a id="top"></a>`)
	fmt.Fprintf(fd, "#  PAGE %s\n\n", c.Title)
	if err := c.Brief.Dump(fd, reg); err != nil {
		return err
	}
	if err := c.Detailed.Dump(fd, reg); err != nil {
		return err
	}
	for _, s := range c.SectionDef {
		if err := s.Dump(fd, reg); err != nil {
			return err
		}
	}
	return nil
}
