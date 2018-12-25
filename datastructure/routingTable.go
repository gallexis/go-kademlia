package datastructure

import (
    log "github.com/sirupsen/logrus"
    "math"
)

var (
    K     = 8
    Alpha = 3
)

type RoutingTable struct {
    KBuckets [BitsInNodeID]KBucket
    NodeID   NodeID
    K        int
    Alpha    int
}

func NewRoutingTable(nodeID NodeID) RoutingTable {
    return NewRoutingTableWithDetails(nodeID, K, Alpha)
}

func NewRoutingTableWithDetails(nodeID NodeID, k int, alpha int) RoutingTable {
    rt := RoutingTable{
        KBuckets: [BitsInNodeID]KBucket{},
        NodeID:   nodeID,
        K:        k,
        Alpha:    alpha,
    }

    for i := 0; i < BitsInNodeID; i++ {
        rt.KBuckets[i] = NewKBucket(k)
    }

    return rt
}

func (rt *RoutingTable) Display() {
    log.Debug("----Display RT--------------------------")
    for i, b := range rt.KBuckets {
        if b.Contacts.Len() > 0 {
            log.Debug(159-i, ": ", b.Contacts.Len())
        }
    }
    log.Debug("----------------------------------------")
}

func (rt *RoutingTable) DisplayBucket(bucketNumber int) {
    log.Debug(159-bucketNumber, ": ", rt.KBuckets[bucketNumber].Contacts.Len())
}

func (rt *RoutingTable) Insert(newContact Contact, pingNode func(chan bool)) {
    xoredID := rt.NodeID.XOR(newContact.NodeID)
    position := rt.NodeID.GetBucketNumber(xoredID)

    rt.KBuckets[position].mutex.Lock()
    rt.KBuckets[position].Insert(newContact, pingNode)
    rt.KBuckets[position].mutex.Unlock()
}

func (rt *RoutingTable) GetRandomContacts(bucketPosition int) []Contact {
    rt.KBuckets[bucketPosition].mutex.Lock()
    defer rt.KBuckets[bucketPosition].mutex.Unlock()

    return rt.KBuckets[bucketPosition].GetRandomContacts(Alpha)
}

func (rt *RoutingTable) GetClosest() (contacts []Contact) {
    bucketNumber := rt.GetLatestBucketFilled()

    rt.KBuckets[bucketNumber].mutex.Lock()
    defer rt.KBuckets[bucketNumber].mutex.Unlock()

    rt.fillContactsByBucketNumber(bucketNumber, &contacts)

    if len(contacts) >= rt.K {
        return
    }

    alternatePositions := generateClosestNeighboursPositions(bucketNumber)
    rt.getClosestNeighbours(alternatePositions, &contacts)

    return
}

func (rt *RoutingTable) Get(otherID NodeID) (contacts []Contact) {
    xoredID := rt.NodeID.XOR(otherID)
    bucketNumber := rt.NodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].mutex.Lock()
    defer rt.KBuckets[bucketNumber].mutex.Unlock()

    rt.fillContactsByBucketNumber(bucketNumber, &contacts)

    if len(contacts) >= rt.K {
        return
    }

    alternatePositions := generateClosestNeighboursPositions(bucketNumber)
    rt.getClosestNeighbours(alternatePositions, &contacts)

    return
}

func (rt *RoutingTable) GetLatestBucketFilled() int {
    last := 159

    for i := 0; i < K; i++ {
        if rt.KBuckets[i].Contacts.Len() > 0 {
            last = 159 - i
        }
    }
    return last
}

func (rt *RoutingTable) GetOne(otherID NodeID) (Contact, bool) {
    xoredID := rt.NodeID.XOR(otherID)
    bucketNumber := rt.NodeID.GetBucketNumber(xoredID)

    rt.KBuckets[bucketNumber].mutex.Lock()
    defer rt.KBuckets[bucketNumber].mutex.Unlock()

    value, exists := rt.KBuckets[bucketNumber].Contacts.Peek(otherID)
    if exists {
        return value.(Contact), exists
    } else {
        return Contact{}, exists
    }
}

func (rt RoutingTable) fillContactsByBucketNumber(bucketNumber int, contacts *[]Contact) {
    contactsInterface := rt.KBuckets[bucketNumber].Contacts.Keys()

    for _, nodeID := range contactsInterface {
        contact, _ := rt.KBuckets[bucketNumber].Contacts.Peek(nodeID.(NodeID))
        *contacts = append(*contacts, contact.(Contact))
    }
    return
}

func (rt RoutingTable) getClosestNeighbours(positions []int, contacts *[]Contact) {
    for _, bucketNumber := range positions {
        rt.fillContactsByBucketNumber(bucketNumber, contacts)

        if len(*contacts) >= rt.K {
            *contacts = (*contacts)[:rt.K]
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
