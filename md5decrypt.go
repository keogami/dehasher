package main

import (
  "fmt"
  "io/ioutil"
  "errors"
  "net/http"
)

func md5decrypt(h Hash) (string, error) {
  const apiEndPoint = "https://md5decrypt.net/Api/api.php?hash=%s&hash_type=%s&email=deanna_abshire@proxymail.eu&code=1152464b80a61728"
  requestUrl := fmt.Sprintf(apiEndPoint, h.Value, h.Type)

  resp, err := GetClient().Get(requestUrl)
  if err != nil {
    return "", fmt.Errorf("md5decrypt: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return "", errors.New("md5decrypt: Invalid Status Code")
  }

  rbody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", fmt.Errorf("md5decrypt: %w", err)
  }

  if len(rbody) == 0 {
    return "", errors.New("md5decrypt: Empty Response")
  }

  result := string(rbody)

  // Todo: check for server sent errors
  return result, nil
}