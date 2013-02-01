package main

import (
  "path/filepath"
  "bufio"
  "os"
  "flag"
  "fmt"
  "bytes"
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
  ignores []string
}

func (p * IndexDb) AddIgnore(pattern string) {
  p.ignores = append(p.ignores,pattern)
}

func (p * IndexDb) AddFile(fileitem FileItem) {
  p.files = append(p.files, fileitem)
}

func (p * IndexDb) MakeIndex(root string) error {
  var buf bytes.Buffer        // Stand-in for a network connection
  enc := gob.NewEncoder(&buf) // Will write to network.
  var err error = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
	fileitem := FileItem{ Size: fi.Size(), Name: fi.Name(), Path: path }
	p.AddFile(fileitem)

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
    db.AddIgnore(".git")
    db.AddIgnore(".svn")
    db.AddIgnore(".hg")
    db.MakeIndex(root)
  } else {
    // search from index
  }

  // err := filepath.Walk(root, visit)
  _ = db
  _ = flagIndex
  // fmt.Printf("filepath.Walk() returned %v\n", err)
}

