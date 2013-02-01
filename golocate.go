package main

import (
  "path/filepath"
  "bufio"
  "os"
  "flag"
  "fmt"
  "bytes"
  "regexp"
  "strings"
  "encoding/gob"
  "log"
)

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s %d\n", path, f.Size() )
  return nil
}

type IndexDb struct {
  createtime int32
  // paths []string
  files []FileItem
  ignoreStrings []string
  ignorePatterns []string
}

func (p * IndexDb) AddIgnorePattern(pattern string) {

}

func (p * IndexDb) AddIgnoreString(pattern string) {
  p.ignoreStrings = append(p.ignoreStrings, pattern)
}

func (p * IndexDb) AddFile(path string, fi os.FileInfo) error {
  for _, str := range p.ignoreStrings {
    if strings.Contains(path,str) {
      return nil
    }
  }

  for _, pattern := range p.ignorePatterns {
    if m, _ := regexp.MatchString(pattern, path); m {
      return nil
    }
  }

  fileitem := FileItem{ Size: fi.Size(), Name: fi.Name(), Path: path }
  // use regexp to compare the results
  p.files = append(p.files, fileitem)
  return nil
}

func (p * IndexDb) SearchFile(pattern string) {
}

func (p * IndexDb) MakeIndex(root string) error {
  var buf bytes.Buffer        // Stand-in for a network connection
  enc := gob.NewEncoder(&buf) // Will write to network.
  var err error = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
    p.AddFile(path,fi)

    var encodeErr error = enc.Encode(path)
    if encodeErr != nil {
      log.Fatal("encode error:", encodeErr)
    }
    fmt.Printf("Visited: %s %d\n", path, fi.Size() )
    return nil
  })

  // write buffer to an index file
  var indexFileName string = ".golocate.db"
  file, err := os.Create(indexFileName)
  file.Write( buf.Bytes() )
  file.Close()
  _ = enc
  return err
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

