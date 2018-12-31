package datastructure

import "time"

type Node struct {
    ContactInfo         Contact
    LastMessageReceived time.Time
    LastFindNodeRequest time.Time
}

func NewNode(contactInfo Contact) Node {
    return Node{
        ContactInfo: contactInfo,
    }
}

func (n *Node) RequestFindNode() bool {
    now := time.Now()

    if n.LastFindNodeRequest.Add(time.Second * 10).After(now) {
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
