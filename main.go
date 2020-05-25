package main

import (
  "bufio"
  "fmt"
  "github.com/bluele/go-chord"
  "os"
  "strings"
  "time"
)

func main() {

  reader := bufio.NewReader(os.Stdin)
  fmt.Print("this address: ")
  text, _ := reader.ReadString('\n')
  // convert CRLF to LF
  text = strings.Replace(text, "\n", "", -1)

  userServer := text
  t, e := initTransport(userServer)
  if e != nil {
    return
  }

  chordConfig := chord.DefaultConfig(userServer)
  chordConfig.NumVnodes = 2
  chordConfig.NumSuccessors = 2

  var ring *chord.Ring = nil

  fmt.Println("Chord Shell")
  fmt.Println("---------------------")

  for {
    fmt.Print("-> ")
    text, _ = reader.ReadString('\n')
    // convert CRLF to LF
    text = strings.Replace(text, "\n", "", -1)

    tokens := strings.Split(text, " ")

    switch tokens[0] {
    case "c":
      ring, e = chord.Create(chordConfig, t)
      if e != nil {
        fmt.Println("Error creating ring")
        fmt.Println(e)
        continue
      }
      fmt.Println(userServer, "created chord ring")
    case "j":
      ring, e = chord.Join(chordConfig, t, tokens[1])
      if e != nil {
        fmt.Println("Error joining ring at ", tokens[1])
        fmt.Println(e)
        continue
      }
      fmt.Println(userServer, "joined", tokens[1])
    case "l":
      if ring == nil {
        fmt.Println("No ring")
        continue
      }
      vns, e := ring.Lookup(2, []byte(tokens[1]))
      if e != nil {
        fmt.Println("Error lookup")
        continue
      }
      for _, vn := range vns {
        fmt.Printf("%s lives on %s with id %s\n", tokens[1], vn.Host, string(vn.Id))
      }
    case "s":
      if ring == nil {
        fmt.Println("No ring")
        continue
      }
      s := ring.Len()
      fmt.Println("Len is", s)
    case "q":
      if ring == nil {
        fmt.Println("No ring")
        continue
      }
      e = ring.Leave()
      if e != nil {
        fmt.Println("Error leaving the ring")
        fmt.Println(e)
        continue
      }
      ring = nil
    case "exit":
      return
    default:
      fmt.Println("Doesn't make sense")
    }

  }
}

func initTransport(thisAddr string) (*chord.TCPTransport, error) {

  t, e := chord.InitTCPTransport(thisAddr, time.Second)
  if e != nil {
    fmt.Println(e)
    return nil, e
  }

  return t, e

}

//func InitTCPTransport(listen string, timeout time.Duration) (*TCPTransport, error)