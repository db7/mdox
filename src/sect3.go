package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

type Sect3 struct {
	XMLName xml.Name `xml:"sect3"`
	Id      string   `xml:"id,attr"`
	Para    []Para   `xml:"para"`
	Title   string   `xml:"title"`
}

func (se *Sect3) Dump(fd io.Writer, reg *Registry) error {
	fmt.Fprintf(fd, "### %s\n\n", se.Title)

	for _, p := range se.Para {
		if err := p.Dump(fd, reg); err != nil {
			return err
		}
	}
	fmt.Fprintln(fd)
	return nil
}
