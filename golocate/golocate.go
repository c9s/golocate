package main

import (
  "golocate"
  "os"
  "flag"
  "fmt"
  "log"
)

func main() {
  var flagUpdate  *bool = flag.Bool("update",false,"Update index")
  var flagInfo    *bool = flag.Bool("info", false,"Show indexdb info")
  var flagVerbose *bool = flag.Bool("v",false,"Verbose output")
  var flagIndex *bool   = flag.Bool("build",false,"Create index")

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

  var indexFilePath string = db.GetLocateDbDir() + "/db"

  // db.SetDbDir( db.GetLocateDbDir() )

  log.Println("Preparing golocate db structure...")
  db.PrepareStructure()

  if *flagIndex {
    log.Println("Building default ignore list...")
    db.IgnoreString(".DS_Store")
    db.IgnoreString(".o")
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
    // db.WriteIndexFile(indexFilePath)
  } else if (*flagUpdate) {
    log.Println( "Updating index..." )
    db.EmptyFileItems()
    db.MakeIndex()
    db.WriteIndexFile(indexFilePath)
  } else if (*flagInfo) {
    db.PrintInfo()
  } else {
    db.SearchString( flag.Arg(0) )
  }
}

