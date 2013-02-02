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
  SourcePaths   []string
  IgnoreFilenames []string
  IgnoreStrings []string
  IgnorePatterns []string
  verbose bool
}

func (p * IndexDb) AddSourcePath(path string) {
  p.SourcePaths = append(p.SourcePaths, path)
}

func (p * IndexDb) GetLocateDbDir() string {
  return os.Getenv("HOME") + "/" + LocateDbDirName
}

func (p * IndexDb) SetVerbose() {
  p.verbose = true
}

func (p * IndexDb) IgnorePattern(pattern string) {
  p.IgnorePatterns = append(p.IgnorePatterns,pattern)
}

func (p * IndexDb) IgnoreString(pattern string) {
  p.IgnoreStrings = append(p.IgnoreStrings, pattern)
}

func (p * IndexDb) IgnoreFilename(filename string) {
  p.IgnoreFilenames = append(p.IgnoreFilenames,filename)
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


/*
Make index from registered paths.
*/
func (p * IndexDb) MakeIndex() {
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
  for _ , path = range p.SourcePaths {
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
  var waiting int = len(p.SourcePaths)
  for ; waiting > 0 ; waiting-- {
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

func (p * IndexDb) LoadIndexFile(filepath string) (error){
  var buf bytes.Buffer
  var dec *gob.Decoder = gob.NewDecoder(&buf)

  // var content []byte, err error = ioutil.ReadFile(filepath)
  file, err := os.Open(filepath)
  if err != nil {
    log.Fatal(err)
  }

  _, err = buf.ReadFrom(file)
  if err != nil {
    log.Fatal(err)
  }

  // db := IndexDb{}
  var decodeErr error = dec.Decode(p)
  if decodeErr != nil {
    log.Fatal("decode error:", decodeErr)
  }
  return err
}

/*
Write indexdb object to file

filepath string
*/
func (p * IndexDb) WriteIndexFile(filepath string) error {
  if p.verbose {
    log.Println("Writing index file...")
  }

  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  var encodeErr error = enc.Encode(p)
  if encodeErr != nil {
    log.Fatal("encode error:", encodeErr)
  }

  file, err := os.Create(filepath)
  file.Write( buf.Bytes() )
  file.Close()

  log.Printf("Done, %d files indexed.", len(p.FileItems) )
  return err
}

