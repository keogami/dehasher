package main

import (
  "net/http"
)

// Todo: Configure client to have better timeouts
var httpClient = &http.Client{}

func GetClient() *http.Client {
  return httpClient
}