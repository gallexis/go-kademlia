package datastructure

import (
    "net"
)

type Node struct {
    Peer
    NodeID              NodeId
}

func NewNode(ip net.IP, port uint16, nodeId NodeId) Node {
    return Node{
        Peer:   NewPeer(ip, port),
        NodeID: nodeId,
    }
}

func (n *Node) Decode(data []byte){
    nodeID := NodeId{}
    nodeID.Decode(data[:20])

    n.NodeID = nodeID
    n.Peer.Decode(data[20:])
}

func (n *Node) Encode() []byte{
    return append(n.NodeID.Encode(), n.Peer.Encode()...)
}
