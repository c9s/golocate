package golocate

import (
  // "regexp"
  //"math"
  // "strings"
  "bytes"
  // "log"
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

  // The directory that contains index db and config
  Dir string

  // database path
  DbPaths []string

  // config object (which is encoded/decoded with Gob)
  Config *IndexDbConfig

  // paths []string
  // FileItems []FileItem
  verbose bool
}

func (p * IndexDb) GetDir() string {
  if p.Dir != "" {
    return p.Dir
  }
  return os.Getenv("HOME") + "/" + LocateDbDirName
}

func (p * IndexDb) SetDbDir(path string) {
  p.Dir = path
}

func (p * IndexDb) SetVerbose() {
  p.verbose = true
}

func (p * IndexDb) PrepareStructure() error {
  return os.Mkdir( p.GetDir() ,0777)
}




/*
func (p * IndexDb) EmptyFileItems() {
  p.FileItems = []FileItem{}
}
*/

/*
func (p * IndexDb) AppendFileItems(old2 []FileItem) []FileItem {
  old1 := p.FileItems
  newslice := make([]FileItem, len(old1) + len(old2))
  copy(newslice, old1)
  copy(newslice[len(old1):], old2)
  return newslice
}
*/


/*
func (p * IndexDb) SearchString(str string) {
  // split fileitems into chunks

  var done = make(chan bool)
  var size int = len(p.FileItems)
  search := func(items []FileItem) {
    for _, item := range(items) {
      if strings.Contains(item.Path,str) {
        fmt.Printf("%s\n",item.Path)
      }
    }
    done <- true
  }
  go search(p.FileItems[ 0 : size / 2 ])
  go search(p.FileItems[ size / 2 : size ])
  <-done
}
*/

/*
Make index from registered paths.
*/
/*
func (p * IndexDb) MakeIndex() {
  var filepipe = make(chan *FileItem, FilePipeBufferLength )
  var done = make(chan bool, 5)

  go func() {
    var fileitem *FileItem
    for fileitem = <-filepipe ; fileitem != nil ; {
      p.FileItems = append(p.FileItems,*fileitem)
      if p.verbose {
        fmt.Printf("  Add\t%s\n", fileitem.Path)
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
*/

func (p * IndexDb) TraverseDirectory(root string, ch chan<- *FileItem) (error) {
  var err error = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
    if ! p.Config.IsFileAcceptable(path) {

      if p.verbose {
        fmt.Println("  Skip\t" + path)
      }
      if fi.IsDir() {
        return filepath.SkipDir
      }
      return nil
    }

    var fileitem FileItem = FileItem{ Name: fi.Name(), Path: path }
    ch <- &fileitem
    return nil
  })
  return  err
}

func (p * IndexDb) LoadIndexConfig(path string)  (*IndexDbConfig,error) {
  // initialize a buffer object
  var buf bytes.Buffer

  // initialize gob decoder
  var dec *gob.Decoder = gob.NewDecoder(&buf)

  // open the config file
  file, err := os.Open(path)
  if err != nil {
    return nil, err
    // log.Fatal(err)
  }

  // decode from the content of the file.
  _, err = buf.ReadFrom(file)
  if err != nil {
    return nil, err
    // log.Fatal(err)
  }

  var config IndexDbConfig
  err = dec.Decode(&config)
  if err != nil {
    return nil, err
    // log.Fatal("decode error:", decodeErr)
  }
  return &config, err
}

func (p * IndexDb) LoadIndexDb() error {
  var err error
  var configPath string = filepath.Join(p.GetDir(), "config" )
  // XXX: should be able to add more db paths for searching
  var dbPath string = filepath.Join(p.GetDir(), "db" )
  p.DbPaths = append(p.DbPaths, dbPath)
  p.Config, err = p.LoadIndexConfig(configPath)
  return err
}



/*
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
*/

/*
Write indexdb object to file

filepath string
*/
/*
func (p * IndexDb) WriteIndexFile(filepath string) error {
  if p.verbose {
    log.Println("Writing index file...")
  }

  var buf bytes.Buffer
  enc := gob.NewEncoder(&buf)
  var encodeErr error = enc.Encode(*p)
  if encodeErr != nil {
    log.Fatal("encode error:", encodeErr)
  }

  file, err := os.Create(filepath)
  file.Write( buf.Bytes() )
  file.Close()

  log.Printf("Done, %d files indexed.", len(p.FileItems) )
  return err
}
*/

/*
func (p * IndexDb) PrintInfo() {
  fmt.Printf("Indexed files: %d\n", len(p.FileItems) )
  fmt.Printf("Indexed paths:\n")
  for _ , path := range p.SourcePaths {
    fmt.Printf("  %s\n", path)
  }
}
*/

