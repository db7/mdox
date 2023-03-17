package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Location struct {
	File   string `xml:"file,attr"`
	Line   string `xml:"line,attr"`
	Column string `xml:"column,attr"`
}

type MemberDef struct {
	XMLName    xml.Name    `xml:"memberdef"`
	Name       string      `xml:"name"`
	Brief      Description `xml:"briefdescription"`
	Detailed   Description `xml:"detaileddescription"`
	Location   Location    `xml:"location"`
	Param      []Param     `xml:"param"`
	Kind       string      `xml:"kind,attr"`
	Definition string      `xml:"definition"`
	Args       string      `xml:"argsstring"`
	Id         string      `xml:"id,attr"`
	Type       string      `xml:"type"`
}

type MemberWrapper struct {
	MemberDef
}

func (m *MemberWrapper) UnmarshalXMLx(d *xml.Decoder, start xml.StartElement) error {
	err := d.DecodeElement(&m.MemberDef, &start)
	switch m.Kind {
	case "function":
		fmt.Println("parsed function", m)
	}
	return err
}

type Param struct {
	XMLName xml.Name  `xml:"param"`
	DefName []DefName `xml:"defname"`
}

type DefName string

func (m *MemberDef) nameString() (s string) {
	switch m.Kind {
	case "define":
		s += fmt.Sprintf("Macro `%s`\n\n", m.Name)
		if len(m.Param) == 0 {
			return s
		}
		s += fmt.Sprintf("```c\n%s", m.Name)
		s += "("
		// param
		for i, p := range m.Param {
			for _, d := range p.DefName {
				s += string(d)
			}
			if i < len(m.Param)-1 {
				s += ","
			}

		}
		s += ")\n```\n\n"
	case "function":
		s += fmt.Sprintf("Function `%s`", m.Name)
		s += fmt.Sprintf("\n\n```c\n%s%s\n```", m.Definition, m.Args)
	case "typedef":
		s += fmt.Sprintf("Type `%s`\n", m.Name)
	default:
		fmt.Printf("What to do here? %#v\n", m)
	}

	return s
}

// Dump writes to w all the information about the memberdef.
func (m *MemberDef) Dump(ctx DumpContext, w *Writer) error {
	if strings.HasPrefix(m.Name, "_") {
		return nil
	}
	w.Print("### ")
	w.Printf(" %s ", m.nameString())
	w.Println()

	// Description
	s := ctx.Reg.Style
	ctx.Reg.Style = SEmphasis
	m.Brief.Dump(ctx, w)
	ctx.Reg.Style = s

	/* TODO: add link to original source file
	// Link to the source
	// We are writing to outputDir/ctx.Path
	// Find the relative path to the source:
	in := filepath.Join(".", ctx.Path)
	out := filepath.Join(ctx.OutputDir, ctx.Path)
	common := getCommonPrefix(in, out)
	log.Printf("in: %s", in)
	log.Printf("out: %s", out)
	log.Printf("commin: %s", common)
	in = strings.TrimPrefix(in, common+"/")
	out = strings.TrimPrefix(out, common+"/")

	var pth []string
	for i := 0; i < getDepth(out); i++ {
		pth = append(pth, "..")
	}
	pth = append(pth, in)
	w.Printf( "[See source](%v)\n", filepath.Join(pth...))
	log.Printf("source: %s", filepath.Join(pth...))
	*/
	// Detailed description
	w.Println()
	m.Detailed.Dump(ctx, w)

	// Final space
	w.Println()
	return nil
}

// DumpRow writes to w the member as a row of table.
func (m *MemberDef) DumpRow(ctx DumpContext, w *Writer) error {
	if strings.HasPrefix(m.Name, "_") {
		return nil
	}
	reg := ctx.Reg

	// Disable paragraph new line. It will be controlled by table row.
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}

	w.Printf("| ")
	newRef(m.Name, m.Id).Dump(ctx, w)
	w.Printf(" | ")
	m.Brief.Dump(ctx, w)
	w.Print(" |\n")

	return nil
}
