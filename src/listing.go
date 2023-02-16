package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Listing struct {
	XMLName  xml.Name   `xml:"programlisting"`
	Filename string     `xml:"filename,attr"`
	CodeLine []CodeLine `xml:"codeline"`
}

func (l *Listing) Dump(fd io.Writer, reg *Registry) error {
	fmt.Fprintln(fd)
	fmt.Fprint(fd, "```")
	switch l.Filename {
	case ".c":
		fmt.Fprint(fd, "c")
	default:
	}
	fmt.Fprintln(fd)
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}
	if reg.Disable(References) {
		defer reg.Enable(References)
	}

	for _, c := range l.CodeLine {
		if err := c.Dump(fd, reg); err != nil {
			return err
		}
		fmt.Fprintln(fd)
	}
	fmt.Fprint(fd, "```")
	fmt.Fprintln(fd)
	fmt.Fprintln(fd)
	return nil
}

type CodeLine struct {
	XMLName xml.Name `xml:"codeline"`
	Element
}

func (c *CodeLine) Dump(fd io.Writer, reg *Registry) error {
	var style Style
	style, reg.Style = reg.Style, SListing
	c.Element.Dump(fd, reg)
	reg.Style = style
	return nil
}
