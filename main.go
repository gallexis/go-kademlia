package main

import (
    "math/rand"
    "time"
)

func init(){
    rand.Seed(time.Now().UTC().UnixNano())
}


func main() {
    rt := InitRoutingTable(160, 20)

    for i:=0; i<100000; i++{
        go rt.Update(Contact{IP: "123.123.123.123", Port: 12345, Node: NewNode()})
    }

    rt.Display()
}