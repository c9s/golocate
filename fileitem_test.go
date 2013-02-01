package main

import "testing"

func TestPrettySize (t * testing.T) {
  var bytes string
  bytes = PrettySize(900)
  if bytes != "900B" {
    t.Error("900 Should be 900B")
  }
  if PrettySize(1024) != "1KB" {
    t.Error("1024 Bytes Should 1KB: " + PrettySize(1024) )
  }
  if PrettySize(1024*1024) != "1MB" {
    t.Error("1024 KB Should 1MB: " + PrettySize(1024*1024) )
  }
}


