package main

type Description struct {
	Para  []Para  `xml:"para"`
	Sect1 []Sect1 `xml:"sect1"`
	Sect2 []Sect2 `xml:"sect2"`
	Sect3 []Sect3 `xml:"sect3"`
}

func (d *Description) Dump(ctx DumpContext, w *Writer) error {
	reg := ctx.Reg
	for _, p := range d.Para {
		p.Dump(ctx, w)
		if reg.Option(ParaLine) {
			w.Println()
		}
	}

	for _, p := range d.Sect1 {
		p.Dump(ctx, w)
	}
	for _, p := range d.Sect2 {
		p.Dump(ctx, w)
	}
	for _, p := range d.Sect3 {
		p.Dump(ctx, w)
	}
	return nil
}

func (d *Description) Empty() bool {
	if len(d.Para) > 0 {
		return false
	}
	if len(d.Sect1) > 0 {
		return false
	}
	if len(d.Sect2) > 0 {
		return false
	}
	if len(d.Sect3) > 0 {
		return false
	}
	return true
}
