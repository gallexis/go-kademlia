package datastructure

import (
    "errors"
    "fmt"
    log "github.com/sirupsen/logrus"
    "math"
    "sync"
)

var (
    K     = 8
    Alpha = 3
)

type BucketPosition int

func (b BucketPosition) CloserThan(other BucketPosition) bool {
    return b > other
}

// Closest bucket = 159 (BitsInNodeID)
// Furthest bucket = 0
type RoutingTable struct {
    sync.Mutex
    ClosestBucketFilled BucketPosition

    KBuckets            [BitsInNodeID]KBucket
    selfNodeID          NodeId
    K                   int
    Alpha               int
}

func NewRoutingTable(nodeID NodeId, pingRequests chan Node) RoutingTable {
    return NewRoutingTableWithDetails(nodeID, K, Alpha, pingRequests)
}

func NewRoutingTableWithDetails(nodeID NodeId, k int, alpha int, pingRequests chan Node) RoutingTable {
    rt := RoutingTable{
        KBuckets:            [BitsInNodeID]KBucket{},
        selfNodeID:          nodeID,
        K:                   k,
        Alpha:               alpha,
        ClosestBucketFilled: 0,
    }

    for i := 0; i < BitsInNodeID; i++ {
        rt.KBuckets[i] = NewKBucket(k)
        rt.KBuckets[i].RefreshLoop(pingRequests)
    }

    return rt
}

func (rt *RoutingTable) SetClosestBucketFilled(position BucketPosition) {
    rt.Lock()
    rt.ClosestBucketFilled = position
    rt.Unlock()
}

func (rt *RoutingTable) GetClosestBucketFilled() BucketPosition {
    rt.Lock()
    defer rt.Unlock()
    return rt.ClosestBucketFilled
}

func (rt RoutingTable) String() (content string) {
    rt.Lock()
    defer rt.Unlock()

    content += "----Display RT--------------------------\n"
    for i, b := range rt.KBuckets {
        if b.Nodes.Len() > 0 {
            content += fmt.Sprintln(i, ": ", b.Nodes.Len())
        }
    }
    content += "----------------------------------------\n"

    return content
}

func (rt *RoutingTable) DisplayBucket(bucketNumber int) string {
    return fmt.Sprint(bucketNumber, ": ", rt.KBuckets[bucketNumber].Nodes.Len())
}

func (rt *RoutingTable) Remove(newNode Node) {
    xoredID := rt.selfNodeID.XOR(newNode.NodeID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].Lock()
    rt.KBuckets[bucketNumber].Nodes.Remove(newNode.NodeID)
    rt.KBuckets[bucketNumber].Unlock()
}

func (rt *RoutingTable) Insert(newNode Node, force bool) (bool, error) {
    if newNode.NodeID.Equals(rt.selfNodeID) {
        return false, errors.New("found myself")
    }

    xoredID := rt.selfNodeID.XOR(newNode.NodeID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].Lock()
    defer rt.KBuckets[bucketNumber].Unlock()

    ok, err := rt.KBuckets[bucketNumber].Insert(&newNode, force)
    if ok && bucketNumber.CloserThan(rt.GetClosestBucketFilled()) {
        rt.SetClosestBucketFilled(bucketNumber)
    }

    return ok, err
}

func (rt *RoutingTable) GetRandomNodes(bucketPosition int) []Node {
    rt.KBuckets[bucketPosition].Lock()
    defer rt.KBuckets[bucketPosition].Unlock()

    return rt.KBuckets[bucketPosition].GetRandomNodes(Alpha)
}

func (rt *RoutingTable) GetClosestNodes() (nodes []Node) {
    bucketNumber := rt.GetClosestBucketFilled()
    rt.KBuckets[bucketNumber].Lock()
    defer rt.KBuckets[bucketNumber].Unlock()

    rt.fillNodesByBucketNumber(bucketNumber, &nodes)

    if len(nodes) >= rt.K {
        return
    }

    alternatePositions := generateClosestNeighboursPositions(bucketNumber)
    rt.getClosestNeighbours(alternatePositions, &nodes)

    return
}

func (rt *RoutingTable) GetK(otherID NodeId) (nodes []Node) {
    xoredID := rt.selfNodeID.XOR(otherID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].Lock()
    defer rt.KBuckets[bucketNumber].Unlock()

    rt.fillNodesByBucketNumber(bucketNumber, &nodes)

    if len(nodes) >= rt.K {
        return
    }

    alternatePositions := generateClosestNeighboursPositions(bucketNumber)
    rt.getClosestNeighbours(alternatePositions, &nodes)

    return
}

func (rt *RoutingTable) PeekOne(otherID NodeId) (Node, bool) {
    xoredID := rt.selfNodeID.XOR(otherID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].Lock()
    defer rt.KBuckets[bucketNumber].Unlock()

    node, exists := rt.KBuckets[bucketNumber].Peek(otherID)
    return *node, exists
}

func (rt *RoutingTable) UpdateNodeStatus(nodeId NodeId) (exists bool) {
    xoredID := rt.selfNodeID.XOR(nodeId)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].Lock()
    defer rt.KBuckets[bucketNumber].Unlock()

    if node, exists := rt.KBuckets[bucketNumber].Get(nodeId); exists {
        node.UpdateLastMessageReceived()
    } else {
        log.Error("Node doesn't exist")
    }

    return exists
}

func (rt *RoutingTable) UpdateLastRequestFindNode(node Node) (exists bool) {
    xoredID := rt.selfNodeID.XOR(node.NodeID)
    bucketNumber := rt.selfNodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].Lock()
    defer rt.KBuckets[bucketNumber].Unlock()

    if node, exists := rt.KBuckets[bucketNumber].Peek(node.NodeID); exists {
        node.UpdateLastRequestFindNode()
    } else {
        log.Error("Node doesn't exist")
    }

    return exists
}

func (rt RoutingTable) fillNodesByBucketNumber(bucketNumber BucketPosition, nodes *[]Node) {
    nodesInterface := rt.KBuckets[bucketNumber].Nodes.Keys()

    for _, nodeID := range nodesInterface {
        if node, exists := rt.KBuckets[bucketNumber].Peek(nodeID.(NodeId)); exists {
            *nodes = append(*nodes, *node)
        } else {
            log.Error("Node doesn't exist")
        }
    }
}

func (rt RoutingTable) getClosestNeighbours(positions []BucketPosition, nodes *[]Node) {
    for _, bucketNumber := range positions {
        rt.fillNodesByBucketNumber(bucketNumber, nodes)

        if len(*nodes) >= rt.K {
            *nodes = (*nodes)[:rt.K]
            break
        }
    }
}

func generateClosestNeighboursPositions(origin BucketPosition) (positions []BucketPosition) {
    var after BucketPosition = 0
    var before BucketPosition = 0

    // TODO: write own Max function
    for i := 0; i < int(math.Max(float64(origin), float64(BitsInNodeID-origin))+1); i++ {
        after = (origin + BucketPosition(i)) % BitsInNodeID
        before = origin - BucketPosition(i)

        if after > origin {
            positions = append(positions, after)
        }

        if before >= 0 && before < origin {
            positions = append(positions, before)
        }
    }
    return
}
