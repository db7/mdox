package main

import (
	"encoding/xml"
	"log"
)

type SimpleSect struct {
	XMLName xml.Name `xml:"simplesect"`
	Element
}

func (e *SimpleSect) Dump(ctx DumpContext, w *Writer) (err error) {
	reg := ctx.Reg
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}

	switch e.Attr.Kind {
	case "return":
		w.Printf("\n**Returns:** ")
		err = e.Element.Dump(ctx, w)
		w.Println()
	case "note":
		w.Printf("\n**Note:** ")
		err = e.Element.Dump(ctx, w)
		w.Println()
	case "pre":
		w.Printf("\n**Precondition:** ")
		err = e.Element.Dump(ctx, w)
		w.Println()
	case "post":
		w.Printf("\n**Postcondition:** ")
		err = e.Element.Dump(ctx, w)
		w.Println()
	case "see":
		w.Printf(" (see ")
		err = e.Element.Dump(ctx, w)
		w.Println(")")
	default:
		err = e.Element.Dump(ctx, w)
		w.Println()
		log.Printf("not implemented: %v", e.Attr.Kind)
	}
	return
}
