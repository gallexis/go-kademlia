package handler

import (
    "fmt"
    "kademlia/network/krpc"
)

func OnPing(pingChan chan krpc.PingResponse) {
    for {
        select {
        case ping := <-pingChan:
            fmt.Println(ping)
        }
    }
}
