package main

import (
  "fmt"
  "io/ioutil"
  "errors"
  "net/http"
)

// Todo: (this shares all its code with md5decrypt)
// so move generalize this logic
func nitrxgen(h Hash) (string, error) {
  const apiEndPoint = "https://www.nitrxgen.net/md5db/%s"
  requestUrl := fmt.Sprintf(apiEndPoint, h.Value)

  resp, err := GetClient().Get(requestUrl)
  if err != nil {
    return "", fmt.Errorf("nitrxgen: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return "", errors.New("nitrxgen: Invalid Status Code")
  }

  rbody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", fmt.Errorf("nitrxgen: %w", err)
  }

  if len(rbody) == 0 {
    return "", errors.New("nitrxgen: Empty Response")
  }

  result := string(rbody)

  return result, nil
}