package datastructure

import (
    "errors"
    "fmt"
    "github.com/hashicorp/golang-lru"
    log "github.com/sirupsen/logrus"
    "math/rand"
    "sync"
    "time"
)

type KBucket struct {
    sync.Mutex
    Nodes             *lru.Cache
    K                 int
    InsertDurationMax time.Duration
}

func NewKBucket(k int) KBucket {
    nodes, err := lru.New(k)
    if err != nil {
        log.Fatalln(err)
    }
    return KBucket{
        Nodes:             nodes,
        K:                 k,
        InsertDurationMax: 3 * time.Second,
    }
}

func (kb *KBucket) Insert(newNode *Node, forceInsert bool) (bool, error) {
    newNodeId := newNode.NodeID

    if !kb.isInBucket(newNodeId) && kb.freeSpaceLeft(){
        kb.Nodes.Add(newNodeId, newNode)

    } else if kb.isInBucket(newNodeId) {
        kb.Nodes.Get(newNodeId) // put at tail of LRU

    } else { // Insert when full KB
        keys := kb.Nodes.Keys()
        oldestNodeInterface, ok := kb.Nodes.Peek(keys[0])
        if !ok {
            return false, errors.New("peek not ok, node might have been removed (Mutex problem ?)")
        }

        oldestNode := oldestNodeInterface.(*Node)
        if oldestNode.IsGood() {
            return false, nil
        }
fmt.Println("FORCE INSEERT")
        kb.Nodes.Add(newNodeId, newNode)
    }

    return true, nil
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
        if i >= alpha {
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
