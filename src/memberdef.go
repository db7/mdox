package main

import (
	"encoding/xml"
	"fmt"
	"io"
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

func (m *MemberDef) NameString() (s string) {
	switch m.Kind {
	case "define":
		s += "Macro `" + m.Name
		defer func() { s += "`" }()
		if len(m.Param) == 0 {
			return s
		}
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
		s += ")"
	case "function":
		s += fmt.Sprintf("`Function %s`", m.Name)
		s += fmt.Sprintf("\n\n```c\n%s%s\n```", m.Definition, m.Args)
	case "typedef":
		s += fmt.Sprintf("Type `%s`\n", m.Name)
	default:
		//		fmt.Printf("What to do here: %#v\n", m)
	}

	return s
}

func (m *MemberDef) Dump(fd io.Writer, reg *Registry) error {
	fmt.Fprint(fd, "### ")
	fmt.Fprintf(fd, defAnchor(m.Id))
	fmt.Fprintf(fd, " %s ", m.NameString())
	//fmt.Fprintf(fd, "[&uparrow](#%s)\n\n", m.Kind)
	fmt.Fprintln(fd)
	if err := m.Brief.Dump(fd, reg); err != nil {
		return err
	}
	if err := m.Detailed.Dump(fd, reg); err != nil {
		return err
	}
	fmt.Fprintln(fd)
	//fmt.Fprintf(fd, "%v\n", m.Location)
	return nil
}

func (m *MemberDef) DumpRow(fd io.Writer, reg *Registry) error {
	if reg.Disable(ParaLine) {
		defer reg.Enable(ParaLine)
	}

	nameLink := fmt.Sprintf("[`%s`](#%s)", m.Name, m.Id)
	fmt.Fprintf(fd, "| %s | ", nameLink)
	m.Brief.Dump(fd, reg)
	fmt.Fprint(fd, " |\n")
	return nil
}
