package main

import "strings"

type Style int

const (
	Default Style = iota
	SHighlight
	SComputerOutput
	SVerbatim
	SBold
	SEmphasis
	SListing
)

type StyleElement struct {
	Style Style
	Element
}

func newStyleElement(name string) *StyleElement {
	var style Style
	switch name {
	case "bold":
		style = SBold
	case "verbatim":
		style = SVerbatim
	case "computeroutput":
		style = SComputerOutput
	case "emphasis":
		style = SEmphasis
	default:
		style = Default
	}
	return newStyleElementI(style)
}
func newStyleElementI(style Style) *StyleElement {
	return &StyleElement{
		Style: style,
	}
}

func (e *StyleElement) Dump(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	style := reg.Style

	// TODO(diogo): can we really disable this if?
	//if reg.Style == Default {
	reg.Style = e.Style
	//}
	err := e.Element.Dump(ctx, w)

	reg.Style = style
	return err
}

type TextElement struct {
	Text string
}

func (e *TextElement) start(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	var s string
	switch reg.Style {
	case SBold:
		s = "**"
	case SEmphasis:
		s = "_"
	case SComputerOutput, SVerbatim:
		s = "`"
	case SListing:
		return nil
	}
	w.Print(s)
	return nil
}

func (e *TextElement) end(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	var s string
	switch reg.Style {
	case SBold:
		s = "**"
	case SEmphasis:
		s = "_"
	case SComputerOutput, SVerbatim:
		s = "`"
	case SListing:
		return nil
	}
	w.Print(s)
	return nil
}

func (e *TextElement) Dump(ctx DumpContext, w *Writer) error {
	if e == nil {
		return nil
	}
	var (
		text      = e.Text
		postBlank bool
	)

	// In case a string start or ends with " ", we have to swap the blank with
	// the emphasis or bold markers to have correctly rendered markdown.
	if strings.HasPrefix(text, " ") {
		text = text[1:]
		w.Print(" ")
	}

	if strings.HasSuffix(text, " ") {
		text = text[:len(text)-1]
		postBlank = true
	}

	e.start(ctx, w)
	w.Printf("%s", text)
	e.end(ctx, w)
	if postBlank {
		w.Print(" ")
	}

	return nil
}

func newText(text string) *TextElement {
	return &TextElement{
		Text: text,
	}
}
