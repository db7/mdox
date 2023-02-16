package main

import (
	"encoding/xml"
	"os"
)

type File struct {
	XMLName     xml.Name      `xml:"doxygen"`
	CompoundDef []CompoundDef `xml:"compounddef"`
}

func loadFile(fn string) (*File, error) {
	fd, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var file File
	if err := xml.NewDecoder(fd).Decode(&file); err != nil {
		return nil, err
	}
	return &file, nil
}
