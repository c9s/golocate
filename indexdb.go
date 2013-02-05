package golocate

import (
  // "regexp"
  //"math"
  "strings"
  "bufio"
  "bytes"
  "log"
  "encoding/gob"
  "path/filepath"
  "os"
  "fmt"
)

const (
  Version = "1.0.0"
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
  FileItems []FileItem
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


func (p * IndexDb) SearchString(str string) {
  var dbPath string = p.GetDbPath()
  var err error
  var file *os.File
  var done chan bool = make(chan bool,5)
  var filestream chan string = make(chan string, 20)

  file, err = os.Open(dbPath)
  if err != nil {
    panic(err)
  }

  var reader = bufio.NewReader(file)
  var readLine = func() (string,error) {
    var line []byte
    var err error
    var isPrefix bool = true
    var ln []byte
    for isPrefix && err == nil {
      line, isPrefix, err = reader.ReadLine()
      ln = append(ln, line...)
    }
    return string(ln), err
  }

  // line matcher
  var matcher = func() {
    var line string
    line = <-filestream
    for line != "" {
      line = <-filestream
      if strings.Contains(line,str) {
        fmt.Printf("%s\n",line)
      }
    }
    done<-true
  }

  go matcher()
  go matcher()

  // line reader
  go func() {
    line, err := readLine()
    for err == nil {
      line , err = readLine()
      filestream<-line
    }
    close(filestream)
    done<-true
  }()

  <-done
  <-done
  <-done
  log.Println("Done")


  // split fileitems into chunks
//   var done = make(chan bool)
//   var size int = len(p.FileItems)
//   search := func(items []FileItem) {
//     for _, item := range(items) {
//       if strings.Contains(item.Path,str) {
//         fmt.Printf("%s\n",item.Path)
//       }
//     }
//     done <- true
//   }
//   go search(p.FileItems[ 0 : size / 2 ])
//   go search(p.FileItems[ size / 2 : size ])
//   <-done
}

/*
Make index from registered paths.
*/
func (p * IndexDb) MakeIndex() {
  var filepipe = make(chan *FileItem, FilePipeBufferLength )
  var done = make(chan bool, 5)
  var dbPath = p.GetDbPath()

  go func() {
    var fileitem *FileItem
    file, err := os.Create(dbPath)
    if err != nil {
      panic("Can not open db file to write index.")
    }

    for fileitem = <-filepipe ; fileitem != nil ; {
      p.FileItems = append(p.FileItems,*fileitem)
      if p.verbose {
        fmt.Printf("  Add\t%s\n", fileitem.Path)
      }
      file.Write( []byte(fileitem.Path + "\n") )
      fileitem = <-filepipe
    }
    file.Close()
    done <- true
  }()

  var path string
  for _ , path = range p.Config.SourcePaths {
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
  log.Println("Waiting for all GoRoutines to finish...")
  var waiting int = len(p.Config.SourcePaths)
  for ; waiting > 0 ; waiting-- {
    <-done
  }
  close(filepipe)
  <-done
  log.Println("All GoRoutines are finished.")

  p.Config.IndexedFiles = len(p.FileItems)
}

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

func (p * IndexDb) SaveConfig(path string, config *IndexDbConfig) error {
  var buf bytes.Buffer
  var enc *gob.Encoder = gob.NewEncoder(&buf)

  var encodeErr error = enc.Encode(config)
  if encodeErr != nil {
    log.Fatal("encode error:", encodeErr)
  }
  file, err := os.Create(path)
  file.Write( buf.Bytes() )
  file.Close()
  return err
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

func (p * IndexDb) GetConfigPath() string {
  return filepath.Join(p.GetDir(), "config")
}

func (p * IndexDb) GetDbPath() string {
  return filepath.Join(p.GetDir(), "db" )
}

func (p * IndexDb) Load() error {
  var err error
  var configPath string = p.GetConfigPath()
  // XXX: should be able to add more db paths for searching
  var dbPath string = p.GetDbPath()
  p.DbPaths = append(p.DbPaths, dbPath)
  p.Config, err = p.LoadIndexConfig(configPath)
  return err
}

func (p * IndexDb) Save() error {
  var err error
  p.SaveConfig(p.GetConfigPath(), p.Config)

//   file, err := os.Create(filepath)
//   file.Write( buf.Bytes() )
//   file.Close()
//   log.Printf("Done, %d files indexed.", len(p.FileItems) )
  return err
}





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

func (p * IndexDb) PrintInfo() {
  fmt.Printf("Indexed files: %d\n", p.Config.IndexedFiles)
  fmt.Printf("Indexed paths:\n")
  for _ , path := range p.Config.SourcePaths {
    fmt.Printf("  %s\n", path)
  }
}

