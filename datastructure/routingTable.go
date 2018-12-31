package datastructure

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "math"
)

var (
    K     = 8
    Alpha = 3
)

type RoutingTable struct {
    KBuckets   [BitsInNodeID]KBucket
    selfNodeID NodeId
    K          int
    Alpha      int
}

func NewRoutingTable(nodeID NodeId) RoutingTable {
    return NewRoutingTableWithDetails(nodeID, K, Alpha)
}

func NewRoutingTableWithDetails(nodeID NodeId, k int, alpha int) RoutingTable {
    rt := RoutingTable{
        KBuckets:   [BitsInNodeID]KBucket{},
        selfNodeID: nodeID,
        K:          k,
        Alpha:      alpha,
    }

    for i := 0; i < BitsInNodeID; i++ {
        rt.KBuckets[i] = NewKBucket(k)
    }

    return rt
}

func (rt RoutingTable) String() string {
    var content string

    content += "----Display RT--------------------------\n"
    for i, b := range rt.KBuckets {
        if b.Nodes.Len() > 0 {
            content += fmt.Sprintln(159-i, ": ", b.Nodes.Len())
        }
    }
    content += "----------------------------------------\n"

    return content
}

func (rt *RoutingTable) DisplayBucket(bucketNumber int) string {
    return fmt.Sprint(159-bucketNumber, ": ", rt.KBuckets[bucketNumber].Nodes.Len())
}

func (rt *RoutingTable) Insert(newNode Node, pingNode func(chan bool)) {
    xoredID := rt.selfNodeID.XOR(newNode.ContactInfo.NodeID)
    position := rt.selfNodeID.GetBucketNumber(xoredID)

    if position < 0 {
        log.Error("Bucket position error: ", position, newNode, xoredID)
        return
    }

    rt.KBuckets[position].mutex.Lock()
    rt.KBuckets[position].Insert(newNode, pingNode)
    rt.KBuckets[position].mutex.Unlock()
}

func (rt *RoutingTable) GetRandomNodes(bucketPosition int) []Node {
    rt.KBuckets[bucketPosition].mutex.Lock()
    defer rt.KBuckets[bucketPosition].mutex.Unlock()

    return rt.KBuckets[bucketPosition].GetRandomNodes(Alpha)
}

func (rt *RoutingTable) GetClosestNodes() (nodes []Node) {
    bucketNumber := rt.GetLatestBucketFilled()

    rt.KBuckets[bucketNumber].mutex.Lock()
    defer rt.KBuckets[bucketNumber].mutex.Unlock()

    rt.fillNodesByBucketNumber(bucketNumber, &nodes)

    if len(nodes) >= rt.K {
        return
    }

    alternatePositions := generateClosestNeighboursPositions(bucketNumber)
    rt.getClosestNeighbours(alternatePositions, &nodes)

    return
}

func (rt *RoutingTable) Get(otherID NodeId) (nodes []Node) {
    xoredID := rt.selfNodeID.XOR(otherID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].mutex.Lock()
    defer rt.KBuckets[bucketNumber].mutex.Unlock()

    rt.fillNodesByBucketNumber(bucketNumber, &nodes)

    if len(nodes) >= rt.K {
        return
    }

    alternatePositions := generateClosestNeighboursPositions(bucketNumber)
    rt.getClosestNeighbours(alternatePositions, &nodes)

    return
}

func (rt *RoutingTable) GetLatestBucketFilled() int {
    last := 159

    for i := 0; i < K; i++ {
        if rt.KBuckets[i].Nodes.Len() > 0 {
            last = 159 - i
        }
    }
    return last
}

func (rt *RoutingTable) GetOne(otherID NodeId) (Node, bool) {
    xoredID := rt.selfNodeID.XOR(otherID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].mutex.Lock()
    defer rt.KBuckets[bucketNumber].mutex.Unlock()

    value, exists := rt.KBuckets[bucketNumber].Nodes.Peek(otherID)
    if exists {
        return value.(Node), exists
    } else {
        return Node{}, exists
    }
}

func (rt *RoutingTable) UpdateNodeStatus(otherID NodeId) bool{
    xoredID := rt.selfNodeID.XOR(otherID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].mutex.Lock()
    defer rt.KBuckets[bucketNumber].mutex.Unlock()

    value, exists := rt.KBuckets[bucketNumber].Nodes.Peek(otherID)
    if !exists {
        log.Error("Node doesn't exist")
        return false
    }

    node := value.(Node)
    node.UpdateLastMessageReceived()
    rt.KBuckets[bucketNumber].Nodes.Add(node.ContactInfo.NodeID, node)
    return true
}


func (rt RoutingTable) fillNodesByBucketNumber(bucketNumber int, nodes *[]Node) {
    nodesInterface := rt.KBuckets[bucketNumber].Nodes.Keys()

    for _, nodeID := range nodesInterface {
        node, _ := rt.KBuckets[bucketNumber].Nodes.Peek(nodeID.(NodeId))
        *nodes = append(*nodes, node.(Node))
    }
}

func (rt RoutingTable) getClosestNeighbours(positions []int, nodes *[]Node) {
    for _, bucketNumber := range positions {
        rt.fillNodesByBucketNumber(bucketNumber, nodes)

        if len(*nodes) >= rt.K {
            *nodes = (*nodes)[:rt.K]
            break
        }
    }
}

func generateClosestNeighboursPositions(origin int) []int {
    var positions []int
    after := 0
    before := 0

    // TODO: write own Max function
    for i := 0; i < int(math.Max(float64(origin), float64(BitsInNodeID-origin))+1); i++ {
        after = (origin + i) % BitsInNodeID
        before = origin - i

        if after > origin {
            positions = append(positions, after)
        }

        if before >= 0 && before < origin {
            positions = append(positions, before)
        }
    }
    return positions
}
