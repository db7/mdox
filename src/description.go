package main

import (
	"fmt"
	"io"
)

type Description struct {
	Para  []Para  `xml:"para"`
	Sect3 []Sect3 `xml:"sect3"`
}

func (d *Description) Dump(fd io.Writer, reg *Registry) error {
	for _, p := range d.Para {
		if err := p.Dump(fd, reg); err != nil {
			return err
		}
		if reg.Option(ParaLine) {
			fmt.Fprintln(fd)
		}
	}
	for _, p := range d.Sect3 {
		if err := p.Dump(fd, reg); err != nil {
			return err
		}
	}
	return nil
}
