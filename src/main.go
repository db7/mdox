package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	input  = flag.String("i", "", "input directory with Doxygen XML")
	output = flag.String("o", "output", "output directory")
)

func defAnchor(id string) string {
	return fmt.Sprintf(`<a id="%s"></a>`, id)
}

/*
	func anchor(text, id string) string {
		format := `<p style='text-align: right;'>[%s](%s)</p>`
		return fmt.Sprintf(format, text, id)
	}
*/

type fileInfo struct {
	Name  string
	Dir   string
	Level int
}

func main() {
	var (
		err error
		reg *Registry
	)
	flag.Parse()
	if reg, err = loadRegistry(); err != nil {
		log.Fatal(err)
	}
	if err = dumpFiles(reg); err != nil {
		log.Fatal(err)
	}
	/*	if err = dumpGroups(reg); err != nil {
			log.Fatal(err)
		}
		if err = dumpPages(reg); err != nil {
			log.Fatal(err)
		}*/
	if err = dumpDirs(reg); err != nil {
		log.Fatal(err)
	}
	/*
		if index, err := LoadIndex(*input); err != nil {
			log.Fatal(err)
		} else {
			fd, err := os.Create(filepath.Join(*output, "INDEX.md"))
			if err != nil {
				log.Fatal(err)
			}
			defer fd.Close()
			index.Dump(fd)
		}*/
}

func ensure(err error) {
	if err != nil {
		log.Panic(err)
	}
}
func reverse(s []Dumper) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
func elem(path string) Element {
	var (
		isdir   bool
		refs    []Dumper
		d, p, s string
	)

	if strings.HasSuffix(path, "/") {
		isdir = true
	}
	fmt.Println(path, isdir)
	path = filepath.Clean(path)
	d, _ = filepath.Split(path)
	for d != "" {
		d = filepath.Clean(d)
		d, p = filepath.Split(d)
		if isdir {
			s += "../"
		}
		isdir = true
		p = filepath.Clean(p)
		refs = append(refs, newText(" / "), newRef(p, s+"README.md"))
	}
	reverse(refs)
	refs = append(refs, newText(filepath.Base(path)))
	return newElement(refs...)
}

func loadIndex() (*Index, error) {
	index, err := LoadIndex(*input)
	if err != nil {
		return nil, err
	}
	fd, err := os.Create(filepath.Join(*output, "index.md"))
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return index, index.Dump(fd)
}

func loadRegistry() (*Registry, error) {
	reg := newRegistry(*input)
	files, err := ioutil.ReadDir(*input)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.Contains(file.Name(), ".xml") {
			continue
		}
		if file.Name() == "index.xml" {
			continue
		}
		log.Println("Parsing", file.Name())
		fn := filepath.Join(*input, file.Name())
		root, err := LoadFile(fn)
		if err != nil {
			return nil, fmt.Errorf("could not read %s: %v", file.Name(), err)
		}
		for _, c := range root.CompoundDef {
			reg.Register(c)
		}
	}
	return reg, nil
}

func dumpDirs(reg *Registry) error {
	for _, c := range reg._dir {
		dir := c.Location.File
		odir := filepath.Join(*output, dir)

		if err := os.MkdirAll(odir, 0775); err != nil {
			return err
		}
		fd, err := os.Create(filepath.Join(odir, "README.md"))
		if err != nil {
			return err
		}
		defer fd.Close()
		if err := c.Dump(fd, reg); err != nil {
			return err
		}
		if err := DumpFooter(fd, reg, strings.Count(dir, "/")); err != nil {
			return err
		}
	}
	return nil
}

func dumpFiles(reg *Registry) error {
	for _, c := range reg._files {
		log.Println("Processing", c.Location.File)
		fnn := fmt.Sprintf("%s.md", c.Location.File)
		fn := filepath.Join(*output, fnn)
		dir := filepath.Dir(fn)
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
		fd, err := os.Create(fn)
		if err != nil {
			return err
		}
		defer fd.Close()
		if err := c.Dump(fd, reg); err != nil {
			return err
		}
		if err := DumpFooter(fd, reg, strings.Count(dir, "/")); err != nil {
			return err
		}

	}

	return nil
}

func dumpGroups(reg *Registry) error {
	ifd, err := os.Create(filepath.Join(*output, "INDEX-modules.md"))
	if err != nil {
		return err
	}
	defer ifd.Close()
	fmt.Fprintln(ifd, "# Modules index")
	fmt.Fprintln(ifd)
	fmt.Fprintln(ifd, "|  |  |")
	fmt.Fprintln(ifd, "|--|--|")

	for _, c := range reg._groups {
		fn := filepath.Join(*output, fmt.Sprintf("%s.md", c.CompoundName))
		log.Println("Processing group", c.CompoundName)
		dir := filepath.Dir(fn)
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
		fd, err := os.Create(fn)
		if err != nil {
			return err
		}
		defer fd.Close()
		if err := c.Dump(fd, reg); err != nil {
			return err
		}

		pl := reg.Disable(ParaLine)
		fmt.Fprint(ifd, "| ")
		newRef(c.CompoundName, fn[:len(fn)-3]).Dump(ifd, reg)
		fmt.Fprint(ifd, "|")
		c.Brief.Dump(ifd, reg)
		fmt.Fprintln(ifd, "|")
		if pl {
			reg.Enable(ParaLine)
		}
	}

	return DumpFooter(ifd, reg, 1)
}

func dirUp(depth int) (path string) {
	for i := 0; i < depth; i++ {
		path += "../"
	}
	return
}

func DumpFooter(fd io.Writer, reg *Registry, depth int) error {
	fmt.Fprintln(fd)
	fmt.Fprintln(fd, "---")
	path := dirUp(depth)
	entries := []Entry{
		newEntry(newRef("Home", path+"README")),
		//newEntry(newRef("Modules", path+"INDEX-modules")),
		//newEntry(newRef("Files", path+"INDEX-files")),
		//newEntry(newText("Description")),
	}
	cols := len(entries)
	tab := Table{
		Cols: cols,
		Row: []Row{
			{Entry: emptyEntries(cols)},
			{Entry: entries}},
	}

	fmt.Fprintf(fd, "_Last updated on %s._\n", time.Now().Format("2006.01.02 15:04:05"))
	//tab.Dump(fd, reg)
	_ = tab

	return nil
}

func dumpPages(reg *Registry) error {
	for _, c := range reg._page {
		log.Println("Processing page", c.CompoundName)
		var name string
		if c.CompoundName == "index" {
			name = "README.md"
		} else {
			name = fmt.Sprintf("page-%s.md", c.CompoundName)
		}
		fn := filepath.Join(*output, name)
		dir := filepath.Dir(fn)
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
		fd, err := os.Create(fn)
		if err != nil {
			return err
		}
		defer fd.Close()
		if err := c.Dump(fd, reg); err != nil {
			return err
		}
	}
	return nil
}
