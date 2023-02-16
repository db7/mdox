package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

type SimpleSect struct {
	XMLName xml.Name `xml:"simplesect"`
	Element
}

func (e *SimpleSect) Dump(fd io.Writer, reg *Registry) error {
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}

	switch e.Attr.Kind {
	case "return":
		fmt.Fprintf(fd, "\n**Returns:** ")
		e.Element.Dump(fd, reg)
		fmt.Fprintln(fd)
	case "note":
		fmt.Fprintf(fd, "\n**Note:** ")
		e.Element.Dump(fd, reg)
		fmt.Fprintln(fd)
	case "pre":
		fmt.Fprintf(fd, "\n**Precondition:** ")
		e.Element.Dump(fd, reg)
		fmt.Fprintln(fd)
	case "post":
		fmt.Fprintf(fd, "\n**Postcondition:** ")
		e.Element.Dump(fd, reg)
		fmt.Fprintln(fd)
	case "see":
		fmt.Fprintf(fd, " (see ")
		e.Element.Dump(fd, reg)
		fmt.Fprintln(fd, ")")
	default:
		e.Element.Dump(fd, reg)
		fmt.Fprintln(fd)
		log.Printf("not implemented: %v", e.Attr.Kind)
	}
	return nil
}
