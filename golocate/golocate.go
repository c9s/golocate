package main

import (
  "golocate"
  "os"
  "flag"
  "fmt"
  "log"
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

  db := golocate.IndexDb{}

  if *flagVerbose {
    db.SetVerbose()
  }

  log.Println("Preparing golocate db structure...")
  db.PrepareStructure()

  if *flagIndex {
    db.IgnoreString(".DS_Store")
    db.IgnoreString(".git")
    db.IgnoreString(".svn")
    db.IgnoreString(".hg")
    db.IgnoreString(".sass-cache")

    log.Println("Building index")

    for _ , path := range flag.Args() {
      db.AddSourcePath(path)
    }
    db.MakeIndex()
    db.WriteIndexFile()
  } else {
    // search from index
  }

  _ = db
  _ = flagIndex
  // fmt.Printf("filepath.Walk() returned %v\n", err)
}

