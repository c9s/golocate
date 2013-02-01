package main

import (
  "os"
  "flag"
  "fmt"
)

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s %d\n", path, f.Size() )
  return nil
}


func main() {
  var flagVerbose *bool = flag.Bool("v",false,"Verbose output")
  var flagIndex *bool   = flag.Bool("i",false,"Create index file")

  flag.Usage = func() {
    fmt.Printf("Usage: --index", os.Args[0])
    os.Exit(1)
  }
  flag.Parse()

  /*
  if len(flag.Args()) == 0 {
    fmt.Printf("Usage")
  }
  */

  db := IndexDb{}

  if *flagVerbose {
    db.SetVerbose()
  }
  db.PrepareStructure()

  if *flagIndex {
    db.AddIgnoreString(".DS_Store")
    db.AddIgnoreString(".git")
    db.AddIgnoreString(".svn")
    db.AddIgnoreString(".hg")

    var dir string
    for _ , dir = range flag.Args() {
      db.MakeIndex(dir)
    }
    db.WriteIndexFile()
  } else {
    // search from index
  }

  _ = db
  _ = flagIndex
  // fmt.Printf("filepath.Walk() returned %v\n", err)
}

