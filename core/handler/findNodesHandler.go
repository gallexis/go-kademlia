package handler

import (
    "fmt"
    "kademlia/network/krpc"
)

func OnFindNodes(findNodesChan chan krpc.FindNodeResponse) {
    for {
        select {
        case findNodes := <-findNodesChan:
            fmt.Println(findNodes)
        }
    }
}
