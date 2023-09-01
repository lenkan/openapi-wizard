package main

import (
	"flag"
	"fmt"

	"github.com/lenkan/openapi-wizard/internal/openapi"
)

var print bool
var filename string

func main() {
	flag.BoolVar(&print, "print", false, "Just print the file")
	flag.StringVar(&filename, "filename", "", "")
	flag.Parse()

	if filename == "" {
		flag.Usage()
		return
	}

	spec := openapi.Load(filename)

	if print == true {
		fmt.Println(spec.Print())
	} else {
		result := openapi.FormatTypescriptClient(spec)
		fmt.Println(result)
	}
}
