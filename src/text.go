package main

import (
	"fmt"
	"io"
)

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

func (e *StyleElement) Dump(fd io.Writer, reg *Registry) error {
	style := reg.Style
	if reg.Style == Default {
		reg.Style = e.Style
	}
	err := e.Element.Dump(fd, reg)
	reg.Style = style
	return err
}

type TextElement struct {
	Text string
}

func (e *TextElement) start(fd io.Writer, reg *Registry) error {
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
	fmt.Fprint(fd, s)
	return nil
}

func (e *TextElement) end(fd io.Writer, reg *Registry) error {
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
	fmt.Fprint(fd, s)
	return nil
}

func (e *TextElement) Dump(fd io.Writer, reg *Registry) error {
	if e == nil {
		return nil
	}
	e.start(fd, reg)
	_, err := fmt.Fprintf(fd, "%s", e.Text)
	e.end(fd, reg)
	return err
}

func newText(text string) *TextElement {
	return &TextElement{
		Text: text,
	}
}
