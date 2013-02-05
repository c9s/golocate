golocate
=========

golocate is a locate-like tool to build filelist index for searching, which is
written in Go with concurrency support.

The difference between `locate` and `golocate` is `golocate` saves your 
prefered paths, ignore list for indexing, so the indexdb is pretty small for searching,
and very easy for updating.

golocate uses separate goroutines to build/search index from custom paths, it's fast enough.

## Install

    go get github.com/c9s/golocate

## Build index

    goloc -build ~/Desktop /usr/local/include /etc

To build index with verbose messages:

    goloc -build -v ~/Desktop ~/Downloads

## Update index

    goloc -update

## Search from indexdb

    goloc [pattern]


