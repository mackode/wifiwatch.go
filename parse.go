package main

import (
  "bufio"
  "fmt"
  "io"
  "os/exec"
  "regexp"
  "time"
)

func nmap(subnet string) (io.ReadCloser, error) {
  fmt.Printf("Running nmap in %s\n", subnet)
  cmd := exec.Command("./wifiscan")
  stdoutPipe, err := cmd.StdoutPipe()
  if err != nil {
    return nil, err
  }

  err = cmd.Start()
  if err != nil {
    stdoutPipe.Close()
    return nil, err
  }
  return stdoutPipe, nil
}

func parse(f io.ReadCloser, t time.Time, outCh chan<- Probe) error {
  probe := Probe{Device: Device{}}
  defer f.Close()
  scan := bufio.NewScanner(f)

  for scan.Scan() {
    line := scan.Text()
    ipRegex := regexp.MustCompile(`Nmap scan report for ([\d\.]+)`)
    macRegex := regexp.MustCompile(`MAC Address: ([\w:]+) \ ((.*?)\)`)
    if matches := ipRegex.FindStringSubmatch(line); matches != nil {
      if probe.IP != "" {
        probe = Probe{Device: Device{}}
      }
      probe.IP = matches[1]
    } else if matches := macRegex.FindStringSubmatch(line); matches != nil {
      probe.Device.MAC = matches[1]
      probe.Device.Product = matches[2]
      if !t.IsZero() {
        probe.Timestamp = t
      }
    }
    if probe.IP != "" && probe.Device.MAC != "" {
      outCh <- probe
    }
  }

  return scan.Err()
}
