package main

import (
  "path/filepath"
  "os"
  "flag"
  "fmt"
)

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s %d\n", path, f.Size() )
  return nil
}

func makeIndex(root string) error {
  var err error = filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s %d\n", path, f.Size() )
	return nil
  })
  return err
}

func main() {
  var flagIndex *bool = flag.Bool("index",false,"Create index file")
  flag.Parse()
  root := flag.Arg(0)

  err := makeIndex(root)

  // err := filepath.Walk(root, visit)
  _ = flagIndex
  fmt.Printf("filepath.Walk() returned %v\n", err)
}
