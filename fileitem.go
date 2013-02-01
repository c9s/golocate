package main

import (
  "strconv"
)

type FileItem struct {
    Name string
    Path string
    Size int64
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

