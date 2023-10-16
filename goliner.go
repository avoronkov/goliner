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
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	goimports "golang.org/x/tools/imports"
)

var (
	imports Strings
)

func init() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Var(&imports, "i", "specify import explicitly")

}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [ -i <import_path> ] <codeline> ...\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "Run Golang one-liner.\n")
	fmt.Fprintf(flag.CommandLine.Output(), "For example: %s 'fmt.Println(\"Hello world!\")'\n\n", os.Args[0])
	flag.PrintDefaults()
}

func prepareSourceFile(imports, args []string) (path string, err error) {
	const (
		pkgHead  = "package main"
		mainHead = "func main() {"
		tail     = "}"
	)

	var buffer bytes.Buffer
	fmt.Fprintln(&buffer, pkgHead)
	for _, imp := range imports {
		fmt.Fprintf(&buffer, "import \"%s\"\n", imp)
	}
	fmt.Fprintln(&buffer, mainHead)
	for _, line := range flag.Args() {
		fmt.Fprintf(&buffer, "%s\n", line)
	}
	fmt.Fprintln(&buffer, tail)

	// run imports
	code, err := goimports.Process("/dev/stdin", buffer.Bytes(), nil)
	if err != nil {
		return "", err
	}

	// save to file
	file, err := os.CreateTemp("", "goliner.src")
	if err != nil {
		return "", err
	}
	fmt.Fprintf(file, "%s", code)
	if err = file.Close(); err != nil {
		os.Remove(file.Name())
	}
	path = file.Name() + ".go"
	if err = os.Rename(file.Name(), path); err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Not enough arguments")
	}
	path, err := prepareSourceFile(imports, flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(path)

	cmd := exec.Command("go", "run", path)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed: %v", err)
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
