package main

import (
	"encoding/xml"
	"os"
)

type Root struct {
	XMLName     xml.Name      `xml:"doxygen"`
	CompoundDef []CompoundDef `xml:"compounddef"`
}

func LoadCompounds(fn string) ([]CompoundDef, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root Root
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}
	return root.CompoundDef, nil
}

func LoadFile(fn string) (*Root, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root Root
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}
	return &root, nil
}
