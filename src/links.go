package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"path/filepath"
	"strings"
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

// Dump a reference as a markdown link, ie, `[text](url)`.
func (r *Ref) Dump(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	if !reg.Option(References) {
		r.Element.Dump(ctx, w)
	} else {
		url := r.Attr.RefID
		if ref := reg.get(ctx, r.Attr.RefID); ref != nil {
			switch ref.Kind {
			case KindFile:
				url = getRelativePath(ctx.Path, ref.Location)
				url += ".md"
			case KindGroup:
				url = getRelativePath(ctx.Path, ref.Location)
				log.Println("HHHH", url)
			case KindPage:
				url = getRelativePath(ctx.Path, ref.Location)
			case KindDir:
				path := filepath.Base(ref.Location)
				url = path + "/README.md"
			case KindMacro:
				url = getRelativePath(ctx.Path, ref.Location)
				url += ".md"
				url = fmt.Sprintf("%s#macro-%s", url, strings.ToLower(ref.Name))
			case KindFunc:
				url = getRelativePath(ctx.Path, ref.Location)
				url += ".md"
				url = fmt.Sprintf("%s#function-%s", url, strings.ToLower(ref.Name))
			default:
				panic("what to do?")
			}
		}

		w.Printf("[")
		r.Element.Dump(ctx, w)
		w.Printf("](%s)", url)
	}
	return nil
}

type Ulink struct {
	XMLName xml.Name `xml:"ulink"`
	Element
}

func (e *Ulink) Dump(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	if !reg.Option(References) {
		e.Element.Dump(ctx, w)
	} else {
		url := e.Attr.URL
		w.Printf("[")
		e.Element.Dump(ctx, w)
		w.Printf("](%s)", url)
	}
	return nil
}
