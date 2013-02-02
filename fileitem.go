package golocate

import (
  "strconv"
)

type FileItem struct {
    Name string
    Path string // full path
}

func (p * FileItem) String() string {
  return p.Path
}

func ConcatFileItems(old1, old2 []FileItem) []FileItem {
   newslice := make([]FileItem, len(old1) + len(old2))
   copy(newslice, old1)
   copy(newslice[len(old1):], old2)
   return newslice
}

func PrettySize(bytes int) string {
  if bytes < 1024 {
    return strconv.Itoa(bytes) + "B"
  }
  if bytes < 1024 * 1024 {
    return strconv.Itoa(bytes / 1024) + "KB"
  }
  if bytes < 1024 * 1024 * 1024 {
    return strconv.Itoa(bytes / 1024 / 1024) + "MB"
  }
  if bytes >= 1024 * 1024 * 1024 {
    // if bytes < 1024 * 1024 * 1024 * 1024 {
    return strconv.Itoa(bytes / 1024 / 1024 / 1024) + "GB"
  }
  return ""
}

