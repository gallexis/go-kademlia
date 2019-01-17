package datastructure

import (
    "github.com/hashicorp/golang-lru"
    log "github.com/sirupsen/logrus"
    "math/rand"
    "sync"
    "time"
)

type KBucket struct {
    sync.Mutex
    Nodes *lru.Cache
    K     int
    InsertDurationMax time.Duration
}

func NewKBucket(k int) KBucket {
    nodes, err := lru.New(k)
    if err != nil {
        log.Fatalln(err)
    }
    return KBucket{
        Nodes: nodes,
        K:        k,
        InsertDurationMax: 3 * time.Second,
    }
}

func (kb *KBucket) Insert(newNode *Node, pingNode func(chan bool)){
    newNodeId := newNode.NodeID

    if !kb.isInBucket(newNodeId) && kb.freeSpaceLeft() {
        kb.Nodes.Add(newNodeId, newNode)

    } else if kb.isInBucket(newNodeId) {
        kb.Nodes.Get(newNodeId) // put in top of LRU

    } else { // Insert when full KB
        keys := kb.Nodes.Keys()
        oldestNodeInterface, ok := kb.Nodes.Peek(keys[0])
        if !ok{
            log.Error("Peek not ok, node might have been removed (Mutex problem ?)")
            return
        }
        oldestNode := oldestNodeInterface.(*Node)
        if oldestNode.IsGood(){ // don't remove the node if we know it is good
            return
        }

        oldestNodeID := oldestNode.NodeID
        pingChan := make(chan bool)
        tick := time.Tick(kb.InsertDurationMax)

        pingNode(pingChan)

        select {
        case <-tick:
            kb.Nodes.Add(newNodeId, newNode) //Add will remove oldestContact then add newNode
            //log.Println("add new")
        case <-pingChan:
            oldestNode.UpdateLastMessageReceived()
            kb.Nodes.Add(oldestNodeID, oldestNode) // if the oldest answers, put it back to the tail
            log.Info("add old node")
        }
    }
}

func (kb KBucket) Get(nodeID NodeId) (*Node, bool) {
    if value, exists := kb.Nodes.Peek(nodeID); exists {
        node := value.(*Node)
        return node, true
    } else {
        return &Node{}, false
    }
}

func (kb KBucket) GetRandomNodes(alpha int) []Node {
    keys := kb.Nodes.Keys()
    contactsLength := len(keys)
    var nodes []Node

    for i, key := range rand.Perm(contactsLength) {
        if i >= alpha{
            break
        }
        node, _ := kb.Nodes.Peek(keys[key].(NodeId))
        nodes = append(nodes, node.(Node))
    }
    return nodes
}

func (kb KBucket) freeSpaceLeft() bool {
    return kb.Nodes.Len() < kb.K
}

func (kb KBucket) isInBucket(nid NodeId) bool {
    return kb.Nodes.Contains(nid)
}
