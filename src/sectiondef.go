package main

import (
	"encoding/xml"
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

// Dump all members of a section according to the section type.
func (s *SectionDef) Dump(ctx DumpContext, w *Writer) error {
	for _, p := range []memberSelection{
		{"define", "Macros", "| Macro | Description |\n|-|-|"},
		{"function", "Functions", "| Function | Description |\n|-|-|"},
		//{"typedef", "Type definitons", "| Type | Description |\n|-|-|"},
	} {
		var members []MemberWrapper
		// Pick only members with content
		for _, m := range s.getMember(p.kind) {
			if !m.Brief.Empty() {
				members = append(members, m)
			}
		}
		if len(members) == 0 {
			continue
		}

		w.Println("---")
		w.Printf("# %s \n\n", p.title)

		w.Printf("%s\n", p.header)
		for _, m := range members {
			m.DumpRow(ctx, w)
		}
		w.Println()

		for _, m := range members {
			m.Dump(ctx, w)
		}

	}
	return nil
}
