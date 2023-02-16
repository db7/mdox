package main

import (
	"encoding/xml"
	"fmt"
	"io"

	_ "github.com/dennwc/go-doxy/xmlfile"
)

type Para struct {
	XMLName xml.Name `xml:"para"`
	Element
}

func (e *Para) Dump(fd io.Writer, reg *Registry) error {
	e.Element.Dump(fd, reg)
	if reg.Option(ParaLine) {
		fmt.Fprintln(fd)
	}
	return nil
}
