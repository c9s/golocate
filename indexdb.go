package golocate

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
  FilePipeBufferLength = 10
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
  var filepipe = make(chan *FileItem, FilePipeBufferLength )
  var done = make(chan bool, 5)

  go func() {
    var fileitem *FileItem
    for fileitem = <-filepipe ; fileitem != nil ; {
      p.FileItems = append(p.FileItems,*fileitem)
      if p.verbose {
        fmt.Printf("  Add\t%s %s\n", fileitem.Path, PrettySize( int(fileitem.Size) ) )
      }
      fileitem = <-filepipe
    }
    done <- true
  }()

  var path string
  for _ , path = range paths {
    log.Println("Building index from " + path)
    // Launch Goroutine
    go func(path string) {
      err := p.TraverseDirectory(path,filepipe)
      if err != nil {
        log.Fatal(err)
      }
      done <- true
    }(path)
  }

  // waiting for all goroutines finish
  for i := 0 ; i < len(paths); i++ {
    <-done
  }
  close(filepipe)
  <-done
}

func (p * IndexDb) TraverseDirectory(root string, ch chan<- *FileItem) (error) {
  var err error = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
    if ! p.fileAcceptable(path) {

      if p.verbose {
        fmt.Println("  Skip\t" + path)
      }
      if fi.IsDir() {
        return filepath.SkipDir
      }
      return nil
    }

    var fileitem FileItem = FileItem{ Size: fi.Size(), Name: fi.Name(), Path: path }
    ch <- &fileitem

    return nil
  })
  return  err
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

  log.Println("Done")
  return err
}

