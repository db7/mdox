package main

import (
	"encoding/xml"

	_ "github.com/dennwc/go-doxy/xmlfile"
)

type Para struct {
	XMLName xml.Name `xml:"para"`
	Element
}

// Dump writes to fd a paragraph and adds a new line if ParaLine option is
// enabled.
func (e *Para) Dump(ctx DumpContext, w *Writer) error {
	e.Element.Dump(ctx, w)
	if ctx.Reg.Option(ParaLine) {
		w.Println()
	}
	return nil
}
