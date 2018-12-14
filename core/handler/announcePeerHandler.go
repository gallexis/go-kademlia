package handler

import (
    "fmt"
    "kademlia/network/messages"
)

func OnAnnouncePeer(announcePeerChan chan messages.AnnouncePeersResponse) {
    for {
        select {
        case announcePeer := <- announcePeerChan:
            fmt.Println(announcePeer)
        }
    }
}
