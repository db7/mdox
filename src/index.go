package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Member struct {
	XMLName xml.Name `xml:"member"`
	RefID   string   `xml:"refid,attr"`
	Kind    string   `xml:"kind,attr"`
	Name    string   `xml:"name"`
}

type Compound struct {
	XMLName xml.Name `xml:"compound"`
	RefID   string   `xml:"refid,attr"`
	Kind    string   `xml:"kind,attr"`
	Name    string   `xml:"name"`
	Member  []Member `xml:"member"`
}

type InnerFile struct {
	XMLName xml.Name `xml:"innerfile"`
	RefID   string   `xml:"refid,attr"`
	Name    string   `xml:",innerxml"`
}

type InnerDir struct {
	XMLName xml.Name `xml:"innerdir"`
	RefID   string   `xml:"refid,attr"`
	Name    string   `xml:",innerxml"`
}

type InnerGroup struct {
	XMLName xml.Name `xml:"innergroup"`
	RefID   string   `xml:"refid,attr"`
	Name    string   `xml:",innerxml"`
}

type Index struct {
	XMLName  xml.Name   `xml:"doxygenindex"`
	Compound []Compound `xml:"compound"`
}

func LoadIndex(dir string) (*Index, error) {
	file, err := os.Open(dir + "/index.xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var index Index
	if err := xml.NewDecoder(file).Decode(&index); err != nil {
		return nil, err
	}
	return &index, nil
}

func (index *Index) Dump(fd io.Writer) error {
	var (
		s  string
		fn string
		k  string
	)

	for _, f := range index.Compound {
		if f.Kind != "file" {
			continue
		}
		if f.Name != fn {
			fn = f.Name
			k = ""

			s = fmt.Sprintf("# %s\n", fn)
			fmt.Fprint(fd, s)
		}
		for _, m := range f.Member {
			if k != m.Kind {
				fmt.Fprintf(fd, "## %s\n", m.Kind)
				k = m.Kind
			}

			s := fmt.Sprintf("- [%s](%s)\n", m.Name, m.RefID)
			fmt.Fprint(fd, s)
		}
	}
	return nil
}
