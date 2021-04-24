package main

import (
  "log"
  "sync"
  "fmt"
  "bufio"
  "io"
  "os"
  "errors"
)

type MultiResults map[Hash]string
func (fr MultiResults) Print() {
  for h, r := range fr {
    row := fmt.Sprintf("%s (%s) = %s", h.Value, h.Type, r)
    fmt.Println(row)
  }
}
type Cracker func(Hash) (string, error)

func GetCrackers(ht HashType) []Cracker {
  cs, ok := map[HashType][]Cracker{
    MD5: []Cracker{ nitrxgen, md5decrypt, hashtoolkit },
  }[ht]
  if !ok {
    return []Cracker(nil)
  }
  return cs
}

func Crack(h Hash, cs []Cracker) (string, bool) {
  res := make(chan string, 1)
  wg := new(sync.WaitGroup)
  wg.Add(len(cs))

  go func(r chan string, w *sync.WaitGroup) {
    wg.Wait()
    close(r)
  }(res, wg)

  for _, crackIt := range cs {
    crackIt := crackIt // make a copy
    go func(r chan string, c Cracker, w *sync.WaitGroup) {
      defer wg.Done()

      result, err := c(h)
      if err != nil {
        log.Println("[Crack]", err)
        return
      }

      // try to send result, but dont block
      select {
      case r <- result:
      default:
      }
    }(res, crackIt, wg)
  }

  result, ok := <- res
  return result, ok
}

func CrackString(hashString string) (string, error) {
  hash, err := ParseHash(hashString)
  if err != nil {
    return "", err
  }

  return CrackHash(hash)
}

func CrackHash(hash Hash) (string, error) {
  crackers := GetCrackers(hash.Type)
  if len(crackers) == 0 {
    return "", errors.New("[NoCrackers] No Crackers supported for the given type")
  }
  
  result, ok := Crack(hash, crackers)
  if !ok {
    return "", errors.New("[Crack] Couldn't crack the hash")
  }

  return result, nil
}

func CrackFromFile(fname string, workers int) (good MultiResults, bad MultiResults, err error) {
  hashPipe, err := MakePipe(fname, workers)
  if err != nil {
    return
  }
  good = make(MultiResults)
  bad  = make(MultiResults)
  wg := new(sync.WaitGroup)
  wg.Add(workers)
  for i := 0; i < workers; i++ {
    go CrackWorker(hashPipe, good, bad, wg)
  }

  wg.Wait()
  return
}

func MakePipe(fname string, buffer int) (hashPipe <-chan string, err error) {
  in, err := os.Open(fname)
  if err != nil {
    return
  }
  hp := make(chan string, buffer)
  hashPipe = hp
  
  go func(in io.ReadCloser, hp chan<- string) {
    defer in.Close()
    scanner := bufio.NewScanner(in)
    for scanner.Scan() {
      hp <- scanner.Text()
    }
    close(hp)
  }(in, hp)
  return
}

// Bug(keogami): there are concurrent writes to map
// should use channels instead, but eh
func CrackWorker(hp <-chan string, good MultiResults, bad MultiResults, wg *sync.WaitGroup) {
  defer wg.Done()

  for hstring := range hp {
    hash, err := ParseHash(hstring)
    if err != nil {
      bad[hash] = "Invalid hash"
      continue
    }
    result, err := CrackHash(hash)
    if err != nil {
      bad[hash] = "Couldn't crack"
      continue
    }
    good[hash] = result
  }
}