package golocate

import (
  "regexp"
  "strings"
)

type IndexDbConfig struct {
  IgnoreFilenames []string
  IgnoreStrings []string
  IgnorePatterns []string
  SourcePaths   []string
}



func (p * IndexDbConfig) AddSourcePath(path string) {
  p.SourcePaths = append(p.SourcePaths, path)
}

func (p * IndexDbConfig) IgnorePattern(pattern string) {
  p.IgnorePatterns = append(p.IgnorePatterns,pattern)
}

func (p * IndexDbConfig) IgnoreString(pattern string) {
  p.IgnoreStrings = append(p.IgnoreStrings, pattern)
}

func (p * IndexDbConfig) IgnoreFilename(filename string) {
  p.IgnoreFilenames = append(p.IgnoreFilenames,filename)
}

func (p * IndexDbConfig) IsFileAcceptable(path string) bool {
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
