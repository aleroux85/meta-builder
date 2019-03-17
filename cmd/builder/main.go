package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	builder "github.com/aleroux85/meta-builder"
)

func main() {
	flag.Usage = func() {
		fmt.Println(`
	The Meta Build tool (meta-builder) copy and parses
	files into the project from a configuration file
	(by default from the "meta.json" file), files and
	templates (by default from the "meta" folder).`)

		fmt.Print("\nUsage:\n\n")
		fmt.Printf("\t%s [flags]\n\n", os.Args[0])
		fmt.Print("The flags are:\n\n")
		flag.PrintDefaults()
	}

	metaFilePath := flag.String("c", "meta.json", "configuration file")
	srcPath := flag.String("s", "meta", "source folder files")
	dstPath := flag.String("d", "", "destination folder files")
	force := flag.Bool("f", false, "force re/writing of all files")
	// watch := flag.Bool("w", false, "watch file changes")

	help := flag.Bool("h", false, "help message")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	c := builder.NewConfig(*srcPath, *dstPath)
	c.Load(*metaFilePath)
	c.BuildAll(*force)
	if c.Error() != nil {
		log.Fatalf("%+v\n", c.Error())
	}
}
