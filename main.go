package main

import (
  "fmt"
  "log"
  "flag"
  "io/ioutil"
)

// Todo: switch to github.com/urfave/cli/v2 instead of flag
type Options struct {
  Hash    string
  File    string
  Workers int
  Verbose bool
}

func ParseOptions() Options {
  options := Options{}
  flag.StringVar(&options.Hash, "h", "", "the hash to crack")
  flag.StringVar(&options.File, "f", "", "the file from which the hashes should be loaded to crack")
  flag.BoolVar(&options.Verbose, "v", false, "whether the output should be verbose")
  flag.IntVar(&options.Workers, "w", 10, "how many workers to use while cracking")
  flag.Parse()
  return options
}

func main() {
  options := ParseOptions()
  if !options.Verbose {
    log.SetOutput(ioutil.Discard)
  }

  if options.Hash != "" {
    result, err := CrackString(options.Hash)
    if err != nil {
      fmt.Println(err)
      return
    }
    fmt.Println(result)
    return
  }

  if options.File != "" {
    good, bad, err := CrackFromFile(options.File, options.Workers)
    if err != nil {
      fmt.Println(err)
      return
    }
    fmt.Println("Cracked:")
    good.Print()
    fmt.Println("Uncracked:")
    bad.Print()
    return
  }
}