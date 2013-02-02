package main

import (
  "golocate"
  "os"
  "flag"
  "fmt"
  "log"
  // "bufio"
)

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s %d\n", path, f.Size() )
  return nil
}

func main() {
  var flagUpdate  *bool = flag.Bool("u",false,"Update index")
  var flagVerbose *bool = flag.Bool("v",false,"Verbose output")
  var flagIndex *bool   = flag.Bool("i",false,"Create index")

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

  var indexFilepath string = db.GetLocateDbDir() + "/db"

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

    // write index to file
    db.WriteIndexFile(indexFilepath)
  } else if (*flagUpdate) {
    log.Println( "Updating index..." )
    db.LoadIndexFile(indexFilepath)
    db.EmptyFileItems()
    db.MakeIndex()
    db.WriteIndexFile(indexFilepath)
  } else {
    log.Println( "Loading index..." )
    db.LoadIndexFile(indexFilepath)
    // var out *bufio.Writer
    // out = bufio.NewWriter(os.Stdout)
    // search from index
    // fmt.Printf( "%s\n", flag.Arg(0) )
    db.SearchString( flag.Arg(0) )
  }

  _ = db
  _ = flagIndex
  // fmt.Printf("filepath.Walk() returned %v\n", err)
}

