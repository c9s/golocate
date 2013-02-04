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

  // db.SetDbDir( db.GetLocateDbDir() )

  log.Println("Preparing golocate db structure...")
  db.PrepareStructure()

  if *flagIndex {
    log.Println("Building default ignore list...")

    // Initialize indexdb
    db.Config = &golocate.IndexDbConfig{}
    db.Config.IgnoreString(".DS_Store")
    db.Config.IgnoreString(".o")
    db.Config.IgnoreString(".git")
    db.Config.IgnoreString(".svn")
    db.Config.IgnoreString(".hg")
    db.Config.IgnoreString(".sass-cache")

    log.Println("Building index")
    for _ , path := range flag.Args() {
      db.Config.AddSourcePath(path)
    }

    db.MakeIndex()

    // write index to file
    // db.WriteIndexFile(indexFilePath)
  } else if (*flagUpdate) {
    log.Println( "Updating index..." )
    db.Load()

    db.MakeIndex()
    // db.WriteIndexFile(indexFilePath)
  } else if (*flagInfo) {
    // db.PrintInfo()
  } else {
    // db.SearchString( flag.Arg(0) )
  }
}

