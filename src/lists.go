package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Item struct {
	Element
}

func (i *Item) Dump(fd io.Writer, reg *Registry) error {
	fmt.Fprint(fd, "- ")
	i.Element.Dump(fd, reg)
	return nil
}

type ParameterList struct {
	XMLName xml.Name `xml:"parameterlist"`
	Element
}

func (s *ParameterList) Dump(fd io.Writer, reg *Registry) error {
	fmt.Fprintln(fd)
	switch s.Attr.Kind {
	case "param":
		fmt.Fprintf(fd, "**Parameters:**\n ")
	default:
		//		log.Printf("not implemented: %v", s.Attr.Kind)
	}
	return nil
}

type ParameterItem struct {
	XMLName xml.Name `xml:"parameteritem"`
	Item
}

type Table struct {
	XMLName xml.Name `xml:"table"`
	Rows    int      `xml:"rows,attr"`
	Cols    int      `xml:"cols,attr"`
	Row     []Row    `xml:"row"`
}

func (t *Table) header(fd io.Writer) error {
	fmt.Fprint(fd, "|")
	for i := 0; i < t.Cols; i++ {
		fmt.Fprint(fd, " --- |")
	}
	fmt.Fprintln(fd)
	return nil
}

func (t *Table) Dump(fd io.Writer, reg *Registry) error {
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}

	header := false
	fmt.Fprintln(fd)
	fmt.Fprintln(fd)
	for _, r := range t.Row {
		if err := r.Dump(fd, reg); err != nil {
			return err
		}
		if !header {
			if err := t.header(fd); err != nil {
				return err
			}
			header = true
		}
	}
	fmt.Fprintln(fd)
	return nil
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Thead   string   `xml:"thead,attr"`
	Para    Para     `xml:"para"`
}

func newEntry(v ...Dumper) Entry {
	return Entry{
		Para: Para{
			Element: newElement(v...),
		},
	}
}
func (e *Entry) Dump(fd io.Writer, reg *Registry) error {
	if e == nil {
		return nil
	}
	return e.Para.Dump(fd, reg)
}

func emptyEntries(n int) []Entry {
	return make([]Entry, n)
}

type Row struct {
	XMLName xml.Name `xml:"row"`
	Entry   []Entry  `xml:"entry"`
}

func (r *Row) Dump(fd io.Writer, reg *Registry) error {
	fmt.Fprint(fd, "| ")
	for _, e := range r.Entry {
		e.Dump(fd, reg)
		fmt.Fprint(fd, "|")
	}
	fmt.Fprintln(fd)
	return nil
}
