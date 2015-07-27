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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	head = `package main

func main() {
`
	tail = `}`
)

func main() {
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
		fmt.Fprint(file, head)
		for _, line := range os.Args[1:] {
			fmt.Fprintf(file, "%s\n", line)
		}
		fmt.Fprint(file, tail)
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
