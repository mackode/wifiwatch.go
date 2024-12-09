package main

import (
  "flag"
  "time"
)

func main() {
  update := flag.Bool("update", false, "update db")
  flag.Parse()

  if *update {
    updater()
    return
  }

  tree()
}

func updater() {
  db, err := NewDB()
  if err != nil {
    panic(err)
  }

  f, err := nmap("192.168.0.0/24")
  if err != nil {
    panic(err)
  }

  ch := make(chan Probe)
  go func() {
    err = parse(f, time.Now(), ch)
    if err != nil {
      panic(err)
    }
    close(ch)
  }()

  for probe := range ch {
    db.Add(probe.IP, probe.Device.MAC, probe.Device.Product, probe.Timestamp)
  }
}
