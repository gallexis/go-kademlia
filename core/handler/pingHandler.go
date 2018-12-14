package handler

import (
    "fmt"
    "kademlia/network/messages"
)

func OnPing(pingChan chan messages.PingResponse) {
    for {
        select {
        case ping := <-pingChan:
            fmt.Println(ping)
        }
    }
}
