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
  LocateDbDirName = ".golocate"
)

type IndexDb struct {
  createtime int32
  // paths []string
  FileItems []FileItem
  IgnoreStrings []string
  IgnorePatterns []string
  verbose bool
}

func (p * IndexDb) GetLocateDbDir() string {
  return os.Getenv("HOME") + "/" + LocateDbDirName
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

func (p * IndexDb) fileAcceptable(path string) bool {
  for _, str := range p.IgnoreStrings {
    if strings.Contains(path,str) {
      return false
    }
  }

  for _, pattern := range p.IgnorePatterns {
    if m, _ := regexp.MatchString(pattern, path); m {
      return false
    }
  }
  return true
}

func (p * IndexDb) PrepareStructure() error {
  return os.Mkdir( p.GetLocateDbDir() ,0777)
}

func (p * IndexDb) SearchFile(pattern string) {

}

func (p * IndexDb) MakeIndex(paths []string) {
  var ch = make(chan []FileItem, 10)
  var path string
  for _ , path = range paths {
    log.Println("Building index from " + path)
    fileitems, err := p.TraverseDirectory(path,ch)
    if err != nil {
      log.Fatal(err)
      continue
    }
    p.FileItems = ConcatFileItems(p.FileItems, fileitems)
  }
}

func (p * IndexDb) TraverseDirectory(root string, ch chan []FileItem) ([]FileItem, error) {
  var fileitems []FileItem
  var err error = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
    if ! p.fileAcceptable(path) {
      fmt.Println("  Skip\t" + path)
      if fi.IsDir() {
        return filepath.SkipDir
      }
      return nil
    }

    fileitem := FileItem{ Size: fi.Size(), Name: fi.Name(), Path: path }
    fileitems = append(fileitems,fileitem)
    if p.verbose {
      fmt.Printf("  Add\t%s %s\n", path, PrettySize( int(fi.Size()) ) )
    }
    return nil
  })
  return fileitems, err
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

  // write index to file
  var indexFileName string = p.GetLocateDbDir() + "/db"
  file, err := os.Create(indexFileName)
  file.Write( buf.Bytes() )
  file.Close()

  if p.verbose {
    log.Println("Done")
  }
  return err
}

