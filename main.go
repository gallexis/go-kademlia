package main

import (
    "os"
    "fmt"
    "crypto/sha1"
    "io"
    "math/rand"
    "time"
)

func init(){
    rand.Seed(time.Now().UTC().UnixNano())
}

func f1() {
    f, err := os.Open("file.txt")
    if err != nil {
        fmt.Print(err)
    }
    defer f.Close()

    h2 := sha1.New()

    if _, err := io.Copy(h2, f); err != nil {
        fmt.Print(err)
    }

    //h2.Write([]byte("192.168.1.4"))
    fmt.Printf("  %x", h2.Sum(nil) )
}

func main() {
    rt := InitRoutingTable(160, 20)

    for i:=0; i<100000; i++{
        go rt.Update(Contact{IP: "123.123.123.123", Port: 12345, Node: NewNode()})
    }

    rt.Display()
}