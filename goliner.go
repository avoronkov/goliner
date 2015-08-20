// Execute some Go one-liners.
//
// For example:
//
// $ goliner 'println("hello world!")'
//
// hello world!
//
// You can also use stdlib modules like:
//
// $ goliner 'fmt.Println(strings.Join([]string{"foo", "bar"}, ", "))'
//
// foo, bar
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	pkgHead  = "package main"
	mainHead = "func main() {"
	tail     = "}"

	imports Strings
)

func init() {
	flag.Var(&imports, "i", "specify import explicitly")

}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [ -i <import_path> ] <codeline> ...\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Run Golang one-liner.\n")
	fmt.Fprintf(os.Stderr, "For example: %s 'fmt.Println(\"Hello world!\")'\n\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	var (
		err               error
		file              *os.File
		tempPath, srcName string
	)
	defer func() {
		if err == nil && tempPath != "" {
			if e := os.Remove(tempPath); e != nil {
				log.Printf("Warning: %v", e)
			}
		}
		if err != nil {
			log.Printf("Error1: %v\n", err)
			os.Exit(1)
		}
	}()
	if len(os.Args) < 2 {
		err = fmt.Errorf("Not enough arguments")
	}
	if err == nil {
		file, err = ioutil.TempFile("", "goliner.src")
	}
	if err == nil {
		tempPath = file.Name()
		fmt.Fprintln(file, pkgHead)
		for _, imp := range imports {
			fmt.Fprintf(file, "import \"%s\"\n", imp)
		}
		fmt.Fprintln(file, mainHead)
		for _, line := range flag.Args() {
			fmt.Fprintf(file, "%s\n", line)
		}
		fmt.Fprintln(file, tail)
		err = file.Close()
	}
	if err == nil {
		srcName = file.Name() + ".go"
		err = os.Rename(file.Name(), srcName)
	}
	if err == nil {
		tempPath = srcName
		cmd := exec.Command("goimports", "-w", srcName)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		err = cmd.Run()
	}
	if err == nil {
		cmd := exec.Command("go", "run", srcName)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		err = cmd.Run()
	}
}

type Strings []string

func (i Strings) String() string {
	return strings.Join([]string(i), ", ")
}

func (i *Strings) Set(value string) error {
	*i = append(*i, value)
	return nil
}
