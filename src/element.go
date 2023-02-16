package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
)

type Attr struct {
	xml.Attr
	RefID string `xml:"refid,attr"`
	Kind  string `xml:"kind,attr"`
	URL   string `xml:"url,attr"`
}

type Element struct {
	Attr   Attr
	Values []Dumper
}

type Dumper interface {
	Dump(ctx DumpContext, w *Writer) error
}

func newElement(vals ...Dumper) Element {
	return Element{
		Values: vals,
	}
}

func (e *Element) Dump(ctx DumpContext, w *Writer) error {
	for _, v := range e.Values {
		if err := v.Dump(ctx, w); err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func (e *Element) DecodeElement(d *xml.Decoder, start xml.StartElement) (err error) {
	var (
		element Dumper
		name    = start.Name.Local
	)
	switch name {
	// links
	case "ref":
		element = new(Ref)
	case "ulink":
		element = new(Ulink)

	// style elements
	case "bold", "highlight", "computeroutput", "verbatim", "emphasis":
		element = newStyleElement(name)

	// lists
	case "parameterlist":
		element = newText("\n")
	case "orderedlist":
		element = newText("\n")
	case "itemizedlist":
		element = new(ItemizedList)

	case "parameteritem":
		element = new(Item)

	case "listitem":
		log.Println("listitem")
		element = new(ListItem)

	case "variablelist":
		element = newText("\n")
	case "variablelistentry":
		element = new(Item)

	case "term":
		element = newText("")

	case "xrefsect":
		element = newText("")

	case "parameternamelist":
		element = new(Item)
	case "parametername":
		element = new(Item)
	case "parameterdescription":
		element = new(Item)

	// table
	case "table":
		element = new(Table)
	case "row":
		element = new(Row)
	case "entry":
		element = new(Entry)

	// text sections
	case "programlisting":
		element = new(Listing)
	case "simplesect":
		element = new(SimpleSect)
	case "para":
		element = new(Para)

	// special charatecters
	case "lsquo":
		element = newText(" `")
	case "rsquo":
		element = newText(" ")
	case "linebreak":
		element = newText("\n")
	case "blockquote":
		element = newText("\n")
	case "ndash":
		element = newText(" -- ")
	case "sp":
		element = newText(" ")

	default:
		log.Printf("not implemented: %v", name)
		return fmt.Errorf("not implemented: %v", name)
	}

	e.add(element)

	err = d.DecodeElement(element, &start)
	return
}

func (e *Element) add(v Dumper) {
	e.Values = append(e.Values, v)
}

type tokenStream struct {
	tokens []xml.Token
}

func (ts *tokenStream) Token() (xml.Token, error) {
	if len(ts.tokens) == 0 {
		return nil, io.EOF
	}
	var t xml.Token
	t, ts.tokens = ts.tokens[0], ts.tokens[1:]
	return t, nil
}
func decodeTokenStream(v any, tokens ...xml.Token) error {
	ts := &tokenStream{
		tokens: tokens,
	}
	err := xml.NewTokenDecoder(ts).Decode(v)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (e *Element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if err := decodeTokenStream(&e.Attr, start, start.End()); err != nil {
		return err
	}

	for {
		t, err := d.Token()
		if t == nil || t == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("could not read token: %v", err)
		}
		switch t := t.(type) {
		case xml.CharData:
			e.add(newText(string(t)))
		case xml.StartElement:
			if err := e.DecodeElement(d, t); err != nil {
				return fmt.Errorf("could not decode element: %v", err)
			}
		case xml.EndElement:
			if t.Name.Local != start.Name.Local {
				return fmt.Errorf("unexpected end element: %v", t.Name.Local)
			}
			return nil
		}
	}
}
