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
    LastSeenNode      time.Time
    Tick         <-chan time.Time
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
        LastSeenNode:      time.Time{},
        Tick:              time.Tick(time.Second * 80),
    }
}

func (kb *KBucket) RefreshLoop(pingRequests chan Node) {
    go func() {
        for {
            select {

            case <-kb.Tick:
                keys := kb.Nodes.Keys()

                if len(keys) <= 0{
                    continue
                }

                oldestNodeInterface, ok := kb.Nodes.Peek(keys[0]) // todo: UNSAFE
                if !ok {
                    log.Error("Error Peeking node")
                    continue
                }

                oldestNode := oldestNodeInterface.(*Node)
                if !oldestNode.IsGood() {
                    fmt.Print(time.Now().Clock())
                    fmt.Println("- Refreshing :", oldestNode.NodeID)
                    pingRequests <- *oldestNode
                }
            }
        }
    }()
}

func (kb *KBucket) Insert(newNode *Node, forceInsert bool) (inserted bool, err error) {
    newNodeId := newNode.NodeID
    inserted = true

    if !kb.isInBucket(newNodeId) && kb.freeSpaceLeft() {
        kb.Nodes.Add(newNodeId, newNode)

    } else if kb.isInBucket(newNodeId) {
        kb.Nodes.Get(newNodeId) // put at tail of LRU

    } else { // Insert when full KB
        keys := kb.Nodes.Keys()
        if len(keys) <= 0{
            log.Error("should not be here")
            return
        }
        oldestNodeInterface, ok := kb.Nodes.Peek(keys[0]) // todo: UNSAFE
        if !ok {
            inserted = false
            err = errors.New("peek not ok, node might have been removed (Mutex problem ?)")
            return
        }

        if forceInsert{
            fmt.Println("FORCE INSERT")
            kb.Nodes.Add(newNodeId, newNode)
            return
        }

        oldestNode := oldestNodeInterface.(*Node)
        if !oldestNode.IsGood() {
            return false, nil
        }
    }

    kb.LastSeenNode = time.Now()
    return
}

func (kb KBucket) Peek(nodeID NodeId) (*Node, bool) {
    if value, exists := kb.Nodes.Peek(nodeID); exists {
        node := value.(*Node)
        return node, true
    } else {
        return &Node{}, false
    }
}

func (kb KBucket) Get(nodeID NodeId) (*Node, bool) {
    if value, exists := kb.Nodes.Get(nodeID); exists {
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
