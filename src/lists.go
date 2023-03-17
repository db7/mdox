package main

import (
	"encoding/xml"
)

type Item struct {
	Element
}

func (i *Item) Dump(ctx DumpContext, w *Writer) error {
	w.Print("- ")
	i.Element.Dump(ctx, w)
	return nil
}

type ListItem struct {
	XMLName xml.Name `xml:"listitem"`
	Item
}

type ItemizedList struct {
	XMLName xml.Name   `xml:"itemizedlist"`
	Kind    string     `xml:"kind"`
	Item    []ListItem `xml:"listitem"`
}

func (il *ItemizedList) Dump(ctx DumpContext, w *Writer) error {
	if ctx.Reg.Disable(ParaLine) {
		defer ctx.Reg.Enable(ParaLine)
	}
	for _, i := range il.Item {
		i.Dump(ctx, w)
	}
	return nil
}

type ParameterList struct {
	XMLName xml.Name `xml:"parameterlist"`
	Element
}

func (s *ParameterList) Dump(ctx DumpContext, w *Writer) error {
	switch s.Attr.Kind {
	case "param":
		w.Printf("**Parameters:**\n\n")
	default:
		//		log.Printf("not implemented: %v", s.Attr.Kind)
	}
	//if ctx.Reg.Disable(ParaLine) {
	//	defer ctx.Reg.Enable(ParaLine)
	//}
	if ctx.Reg.Disable(PreserveNewLines) {
		defer ctx.Reg.Enable(PreserveNewLines)
	}
	s.Element.Dump(ctx, w)
	w.Println()

	return nil
}

type Table struct {
	XMLName xml.Name `xml:"table"`
	Rows    int      `xml:"rows,attr"`
	Cols    int      `xml:"cols,attr"`
	Row     []Row    `xml:"row"`
}

func (t *Table) header(w *Writer) {
	w.Print("|")
	for i := 0; i < t.Cols; i++ {
		w.Print(" --- |")
	}
	w.Println()
}

func (t *Table) Dump(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}

	header := false
	w.Println()
	w.Println()
	for _, r := range t.Row {
		if err := r.Dump(ctx, w); err != nil {
			return err
		}
		if !header {
			t.header(w)
			header = true
		}
	}
	w.Println()
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
func (e *Entry) Dump(ctx DumpContext, w *Writer) error {
	if e == nil {
		return nil
	}
	return e.Para.Dump(ctx, w)
}

func emptyEntries(n int) []Entry {
	return make([]Entry, n)
}

type Row struct {
	XMLName xml.Name `xml:"row"`
	Entry   []Entry  `xml:"entry"`
}

func (r *Row) Dump(ctx DumpContext, w *Writer) error {
	w.Print("| ")
	for _, e := range r.Entry {
		e.Dump(ctx, w)
		w.Print("|")
	}
	w.Println()
	return nil
}
