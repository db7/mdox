package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

type SectionDef struct {
	XMLName   xml.Name        `xml:"sectiondef"`
	Kind      string          `xml:"kind,attr"`
	MemberDef []MemberWrapper `xml:"memberdef"`
}

func (s *SectionDef) getMember(kind string) (members []MemberWrapper) {
	for _, m := range s.MemberDef {
		if m.Kind == kind {
			members = append(members, m)
		}
	}
	return
}

type memberSelection struct {
	kind   string
	title  string
	header string
}

func (s *SectionDef) Dump(fd io.Writer, reg *Registry) error {
	for _, p := range []memberSelection{
		{"define", "Macros", "| Macro | Description |\n|-|-|"},
		{"function", "Functions", "| Function | Description |\n|-|-|"},
		//{"typedef", "Type definitons", "| Type | Description |\n|-|-|"},
	} {
		members := s.getMember(p.kind)
		if len(members) == 0 {
			continue
		}
		fmt.Fprintf(fd, "## %s \n\n", p.title)

		fmt.Fprintf(fd, "%s\n", p.header)
		for _, m := range members {
			m.DumpRow(fd, reg)
		}
		fmt.Fprintln(fd)

		for _, m := range members {
			m.Dump(fd, reg)
		}
	}
	return nil
}
