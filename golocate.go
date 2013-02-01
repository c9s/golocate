package main

import (
  "bufio"
  "os"
  "flag"
  "fmt"
)

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s %d\n", path, f.Size() )
  return nil
}


var stdout *bufio.Writer
var stderr *bufio.Writer

func main() {
  var flagIndex *bool = flag.Bool("index",false,"Create index file")

  stderr = bufio.NewWriter(os.Stderr)
  stdout = bufio.NewWriter(os.Stdout)

  flag.Usage = func() {
    fmt.Fprintf(stderr, "Usage: --index", os.Args[0])
    os.Exit(1)
  }
  flag.Parse()

  /*
  if len(flag.Args()) == 0 {
    fmt.Printf("Usage")
  }
  */
  root := flag.Arg(0)
  _ = root

  db := IndexDb{}

  db.PrepareStructure()

  if *flagIndex {
    db.AddIgnoreString(".git")
    db.AddIgnoreString(".svn")
    db.AddIgnoreString(".hg")
    db.MakeIndex(root)
  } else {
    // search from index
  }

  // err := filepath.Walk(root, visit)
  _ = db
  _ = flagIndex
  // fmt.Printf("filepath.Walk() returned %v\n", err)
}

