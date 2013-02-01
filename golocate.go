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
  root := flag.Arg(0)
  _ = root

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
    buf,err := db.MakeIndex(root)
    db.WriteIndexFile(buf)
    _ = err
  } else {
    // search from index
  }

  // err := filepath.Walk(root, visit)
  _ = db
  _ = flagIndex
  // fmt.Printf("filepath.Walk() returned %v\n", err)
}

