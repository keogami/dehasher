package main

import (
  "net/http"
  "fmt"
  "regexp"
  "errors"
  "io/ioutil"
)

// i dunno if this works, just ported the python function tbh
func hashtoolkit(h Hash) (string, error) {
  re := regexp.MustCompile(`(?m)/generate-hash/\?text=(.+)\"`)

  const apiEndPoint = "https://hashtoolkit.com/reverse-hash/?hash=%s"
  requestUrl := fmt.Sprintf(apiEndPoint, h.Value)

  resp, err := GetClient().Get(requestUrl)
  if err != nil {
    return "", fmt.Errorf("hashtoolkit: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return "", errors.New("hashtoolkit: Invalid Status Code")
  }

  rbody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", fmt.Errorf("hastoolkit: %w", err)
  }

  matches := re.FindStringSubmatch(string(rbody))
  if len(matches) < 2 {
    return "", errors.New("hashtoolkit: Couldn't find the result in response")
  }

  result := matches[1]

  if result == "" {
    return "", errors.New("hashtoolkit: The result is empty")
  }

  return result, nil
}