package main

import (
	"encoding/xml"
)

type Listing struct {
	XMLName  xml.Name   `xml:"programlisting"`
	Filename string     `xml:"filename,attr"`
	CodeLine []CodeLine `xml:"codeline"`
}

func (l *Listing) Dump(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	w.Println()
	w.Print("```")
	switch l.Filename {
	case ".c":
		w.Print("c")
	default:
	}
	w.Println()
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}
	if reg.Disable(References) {
		defer reg.Enable(References)
	}

	for _, c := range l.CodeLine {
		if err := c.Dump(ctx, w); err != nil {
			return err
		}
		w.Println()
	}
	w.Print("```")
	w.Println()
	w.Println()
	return nil
}

type CodeLine struct {
	XMLName xml.Name `xml:"codeline"`
	Element
}

func (c *CodeLine) Dump(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	var style Style
	style, reg.Style = reg.Style, SListing
	c.Element.Dump(ctx, w)
	reg.Style = style
	return nil
}
