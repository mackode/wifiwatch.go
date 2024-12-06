package main

import (
  "fmt"
  "os/exec"
)

const subnet = "192.168.0.0/24"

func main() {
  cmd := exec.Command("/usr/local/bin/nmap", "-sn", subnet)
  output, err := cmd.Output()
  if err != nil {
    panic(err)
  }
  fmt.Println(string(output))
}
