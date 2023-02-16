package main

import (
	"encoding/xml"
)

type Sect struct {
	Id    string `xml:"id,attr"`
	Title string `xml:"title"`
	Para  []Para `xml:"para"`
}

func (se *Sect) Dump(ctx DumpContext, w *Writer) {
	for _, p := range se.Para {
		p.Dump(ctx, w)
		if ctx.Reg.Option(ParaLine) {
			w.Println()
		}
	}
}

type Sect3 struct {
	XMLName xml.Name `xml:"sect3"`
	Sect
}

func (se *Sect3) Dump(ctx DumpContext, w *Writer) error {
	w.Printf("### %s\n\n", se.Title)
	se.Sect.Dump(ctx, w)
	return nil
}

type Sect2 struct {
	XMLName xml.Name `xml:"sect2"`
	Sect
	Sect3 []Sect3 `xml:"sect3"`
}

func (se *Sect2) Dump(ctx DumpContext, w *Writer) error {
	w.Printf("## %s\n\n", se.Title)
	se.Sect.Dump(ctx, w)
	for _, p := range se.Sect3 {
		p.Dump(ctx, w)
	}
	return nil
}

type Sect1 struct {
	XMLName xml.Name `xml:"sect1"`
	Sect
	Sect2 []Sect2 `xml:"sect2"`
	Sect3 []Sect3 `xml:"sect3"`
}

func (se *Sect1) Dump(ctx DumpContext, w *Writer) error {
	w.Printf("# %s\n\n", se.Title)
	se.Sect.Dump(ctx, w)
	for _, p := range se.Sect2 {
		p.Dump(ctx, w)
	}
	for _, p := range se.Sect3 {
		p.Dump(ctx, w)
	}

	return nil
}
