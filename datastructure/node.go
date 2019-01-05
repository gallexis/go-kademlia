package datastructure

import (
    "sync"
    "time"
)

type Node struct {
    ContactInfo         Contact
    LastMessageReceived time.Time
    LastFindNodeRequest time.Time
    mutex               sync.Mutex
}

func NewNode(contactInfo Contact) Node {
    return Node{
        ContactInfo: contactInfo,
    }
}

func (n *Node) RequestFindNode() bool {
    n.mutex.Lock()
    defer n.mutex.Unlock()

    now := time.Now()

    if n.LastFindNodeRequest.Add(time.Minute).After(now) {
        return false
    }

    n.LastFindNodeRequest = now
    return true
}

func (n *Node) IsGood() bool {
    return n.LastMessageReceived.Add(time.Minute * 15).After(time.Now())
}

func (n *Node) UpdateLastMessageReceived() {
    n.LastMessageReceived = time.Now()
}
