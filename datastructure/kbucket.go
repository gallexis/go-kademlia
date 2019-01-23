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
    K             int
    RefreshTicker <-chan time.Time

    sync.RWMutex
    Nodes         *lru.Cache
}

func NewKBucket(k int) KBucket {
    nodes, err := lru.New(k)
    if err != nil {
        log.Fatalln(err)
    }
    return KBucket{
        Nodes:         nodes,
        K:             k,
        RefreshTicker: time.Tick(time.Minute * 5),
    }
}

func (kb *KBucket) RefreshLoop(pingRequests chan Node) {
    go func() {
        for {
            select {

            case <-kb.RefreshTicker:

                kb.RLock()
                keys := kb.Nodes.Keys()

                if len(keys) <= 0 {
                    kb.RUnlock()
                    continue
                }

                oldestNodeInterface, ok := kb.Nodes.Peek(keys[0])
                kb.RUnlock()
                if !ok {
                    log.Error("Error Peeking node")
                    continue
                }

                oldestNode := oldestNodeInterface.(*Node)
                if !oldestNode.IsGood() {
                    pingRequests <- *oldestNode
                }
            }
        }
    }()
}

func (kb *KBucket) Remove(id NodeId){
    kb.Lock()
    defer kb.Unlock()

    kb.Nodes.Remove(id)
}

func (kb *KBucket) Insert(newNode *Node, forceInsert bool)  (bool, error) {
    newNodeId := newNode.NodeID
    var err error = nil

    kb.Lock()

    if !kb.isInBucket(newNodeId) && kb.freeSpaceLeft() {
        kb.Nodes.Add(newNodeId, newNode)

    } else if kb.isInBucket(newNodeId) {
        kb.Nodes.Get(newNodeId) // put at tail of LRU

    } else { // insert when full KB
        keys := kb.Nodes.Keys()
        if len(keys) <= 0 {
            log.Error("should not be here")

            kb.Unlock()
            return true, err
        }
        oldestNodeInterface, ok := kb.Nodes.Peek(keys[0]) // todo: UNSAFE
        if !ok {
            err = errors.New("peek not ok, node might have been removed (Mutex problem ?)")

            kb.Unlock()
            return false, err
        }

        if forceInsert {
            fmt.Println("FORCE INSERT")
            kb.Nodes.Add(newNodeId, newNode)

            kb.Unlock()
            return true, nil
        }

        oldestNode := oldestNodeInterface.(*Node)
        if !oldestNode.IsGood() {
            kb.Unlock()
            return false, nil
        }
    }

    // Last seen node?

    kb.Unlock()
    return true, err
}

func (kb KBucket) Keys() (keys []NodeId) {
    kb.RLock()
    defer kb.RUnlock()

    for _, key := range kb.Nodes.Keys(){
        keys = append(keys, key.(NodeId))
    }

    return
}

func (kb KBucket) Len() int {
    kb.RLock()
    defer kb.RUnlock()

    return kb.Nodes.Len()
}

func (kb KBucket) Peek(nodeID NodeId) (*Node, bool) {
    kb.RLock()
    defer kb.RUnlock()

    if value, exists := kb.Nodes.Peek(nodeID); exists {
        node := value.(*Node)
        return node, true
    } else {
        return &Node{}, false
    }
}

func (kb KBucket) Get(nodeID NodeId) (*Node, bool) {
    kb.RLock()
    defer kb.RUnlock()

    if value, exists := kb.Nodes.Get(nodeID); exists {
        node := value.(*Node)
        return node, true
    } else {
        return &Node{}, false
    }
}

func (kb KBucket) GetRandomNodes(alpha int) []Node {
    kb.Lock()
    defer kb.Unlock()

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
