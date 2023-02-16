package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

type Ref struct {
	XMLName xml.Name `xml:"ref"`
	Element
}

func newRef(value string, refid string) *Ref {
	return &Ref{
		Element: Element{
			Attr: Attr{
				RefID: refid,
			},
			Values: []Dumper{newText(value)},
		},
	}
}

func (r *Ref) Dump(fd io.Writer, reg *Registry) error {
	if !reg.Option(References) {
		r.Element.Dump(fd, reg)
	} else {
		fmt.Fprintf(fd, "[")
		r.Element.Dump(fd, reg)
		ref, err := reg.search(r.Attr.RefID)
		if err != nil {
			log.Fatal(err)
			return err
		}
		fmt.Fprintf(fd, "](%s)", ref)
	}
	return nil
}

type Ulink struct {
	XMLName xml.Name `xml:"ulink"`
	Element
}

func (e *Ulink) Dump(fd io.Writer, reg *Registry) error {
	if !reg.Option(References) {
		e.Element.Dump(fd, reg)
	} else {
		fmt.Fprintf(fd, "[")
		e.Element.Dump(fd, reg)
		fmt.Fprintf(fd, "](%s)", e.Attr.URL)
	}
	return nil
}
