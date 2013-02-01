package main

import (
  "regexp"
  "strings"
  "bytes"
  "log"
  "encoding/gob"
  "path/filepath"
  "os"
  "fmt"
)

const (
  LocateDbDir = ".golocate"
)

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

func (p * IndexDb) PrepareStructure() {
  os.Mkdir(".golocate",0777)
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

  log.Println("Writing index file...")

  var indexFileName string = LocateDbDir + "/db"
  file, err := os.Create(indexFileName)
  file.Write( buf.Bytes() )
  file.Close()
  _ = enc

  log.Println("Done")
  return err
}
