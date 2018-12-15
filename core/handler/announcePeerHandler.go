package handler

import (
    "fmt"
    "kademlia/network/krpc"
)

func OnAnnouncePeer(announcePeerChan chan krpc.AnnouncePeersResponse) {
    for {
        select {
        case announcePeer := <- announcePeerChan:
            fmt.Println(announcePeer)
        }
    }
}
