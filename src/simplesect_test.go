package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func decode(s string, v any) error {
	return xml.NewDecoder(strings.NewReader(s)).Decode(v)
}

func Test_SimpleSect(t *testing.T) {
	s := `<simplesect kind="return">
	 <para>TEXT1</para>
	 </simplesect>`
	var e SimpleSect
	err := decode(s, &e)
	if err != nil {
		fmt.Println(err)
		t.Errorf("could not decode xml: %v", err)
		t.FailNow()
	}
	fmt.Println()
	e.Dump(os.Stdout, nil)
	//spew.Dump("Parsing:", para)
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Kind    string   `xml:"kind,attr"`
}

func Test_Attr(t *testing.T) {
	s := `<tag kind="something">value</tag>`
	var a Attr
	err := decode(s, &a)
	assert.Nil(t, err)
	assert.Equal(t, "something", a.Kind)
	fmt.Printf("------%#v\n", a)
}

type TestTag struct {
	XMLName xml.Name `xml:"tag"`
	Element
}

func Test_Element(t *testing.T) {
	s := `<tag kind="something">value</tag>`
	var e TestTag
	err := decode(s, &e)
	assert.Nil(t, err)

	assert.Equal(t, "something", e.Attr.Kind)
	assert.Equal(t, 1, len(e.Values))
	if len(e.Values) > 0 {
		v, is := e.Values[0].(*TextElement)
		assert.True(t, is)
		assert.Equal(t, "value", v.Text)
	}

}
func Test_TokenStream(t *testing.T) {
	start := xml.StartElement{Name: xml.Name{Local: "test"}, Attr: []xml.Attr{
		{
			Name:  xml.Name{Local: "kind"},
			Value: "test_attr",
		},
	}}
	tests := []struct {
		name   string
		tokens []xml.Token
		ok     bool
	}{
		{
			name: "OK",
			tokens: []xml.Token{
				start,
				start.End(),
			},
			ok: true,
		},
		{
			name: "Malformed",
			tokens: []xml.Token{
				start,
				xml.StartElement{Name: xml.Name{Local: "bad"}},
				start.End(),
			},
			ok: false,
		},
	}
	for _, tc := range tests {
		for _, eof := range []bool{true, false} {
			name := fmt.Sprintf("%s/earlyEOF=%v", tc.name, eof)
			t.Run(name, func(t *testing.T) {
				d := xml.NewTokenDecoder(&tokenStream{
					tokens: tc.tokens,
				})
				var v struct {
					XMLName xml.Name `xml:"test"`
					Kind    string   `xml:"kind,attr"`
				}
				err := d.Decode(&v)
				if tc.ok && err != nil {
					t.Fatalf("d.Decode: expected nil error, got %v", err)
				}
				if _, ok := err.(*xml.SyntaxError); !tc.ok && !ok {
					t.Errorf("d.Decode: expected syntax error, got %v", err)
				}
				fmt.Println(v)
			})
		}
	}
	var v struct {
		XMLName xml.Name `xml:"test"`
		Kind    string   `xml:"kind,attr"`
	}
	err := decodeTokenStream(&v, start)
	//assert.Equal(t, err)
	_ = err
	fmt.Printf("---> %#v\n", v)
}
