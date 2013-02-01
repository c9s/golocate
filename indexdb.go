package main

import (
  "regexp"
  "strings"
  "bytes"
  "log"
  "encoding/gob"
  "path/filepath"
  "os"
  _ "fmt"
)

const (
  LocateDbDir = ".golocate"
)

type IndexDb struct {
  createtime int32
  // paths []string
  Files []FileItem
  IgnoreStrings []string
  IgnorePatterns []string
  verbose bool
}

func (p * IndexDb) SetVerbose() {
  p.verbose = true
}

func (p * IndexDb) AddIgnorePattern(pattern string) {
  p.IgnorePatterns = append(p.IgnorePatterns,pattern)
}

func (p * IndexDb) AddIgnoreString(pattern string) {
  p.IgnoreStrings = append(p.IgnoreStrings, pattern)
}

func (p * IndexDb) AddFile(path string, fi os.FileInfo) error {
  for _, str := range p.IgnoreStrings {
    if strings.Contains(path,str) {
      return nil
    }
  }

  for _, pattern := range p.IgnorePatterns {
    if m, _ := regexp.MatchString(pattern, path); m {
      return nil
    }
  }

  fileitem := FileItem{ Size: fi.Size(), Name: fi.Name(), Path: path }
  // use regexp to compare the results
  p.Files = append(p.Files, fileitem)

  if p.verbose {
    log.Printf("Added: %s %d\n", path, fi.Size() )
  }
  return nil
}

func (p * IndexDb) PrepareStructure() error {
  return os.Mkdir(".golocate",0777)
}

func (p * IndexDb) SearchFile(pattern string) {

}

func (p * IndexDb) MakeIndex(root string) error {
  return filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
    p.AddFile(path,fi)
    return nil
  })
}

// write buffer to an index file
func (p * IndexDb) WriteIndexFile() error {
  if p.verbose {
    log.Println("Writing index file...")
  }

  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  var encodeErr error = enc.Encode(p)
  if encodeErr != nil {
    log.Fatal("encode error:", encodeErr)
  }

  var indexFileName string = LocateDbDir + "/db"
  file, err := os.Create(indexFileName)
  file.Write( buf.Bytes() )
  file.Close()

  if p.verbose {
    log.Println("Done")
  }
  return err
}

