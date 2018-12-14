package handler

import (
    "fmt"
    "kademlia/network/messages"
)

func OnFindNodes(findNodesChan chan messages.FindNodeResponse) {
    for {
        select {
        case findNodes := <-findNodesChan:
            fmt.Println(findNodes)
        }
    }
}
