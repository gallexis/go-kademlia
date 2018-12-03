package datastructure

import (
	"math"
)

var (
	K     = 20
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
		KBuckets: [160]KBucket{},
		NodeID:   nodeID,
		K:        k,
		Alpha:    alpha,
	}

	for i := 0; i < BitsInNodeID; i++ {
		rt.KBuckets[i] = NewKBucket(k)
	}

	return rt
}

func (rt *RoutingTable) Insert(newContact Contact) {
	xoredID := rt.NodeID.XOR(newContact.NodeID)
	position := rt.NodeID.getBucketNumber(xoredID)

	rt.KBuckets[position].mutex.Lock()
	rt.KBuckets[position].Insert(newContact)
	rt.KBuckets[position].mutex.Unlock()
}

func (rt *RoutingTable) Get(otherID NodeID) (contacts []Contact) {
	xoredID := rt.NodeID.XOR(otherID)
	bucketNumber := rt.NodeID.getBucketNumber(xoredID)

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

func (rt *RoutingTable) GetOne(otherID NodeID) (Contact, bool) {
	xoredID := rt.NodeID.XOR(otherID)
	bucketNumber := rt.NodeID.getBucketNumber(xoredID)

	rt.KBuckets[bucketNumber].mutex.Lock()
	defer rt.KBuckets[bucketNumber].mutex.Unlock()

	value, exists := rt.KBuckets[bucketNumber].Contacts.Peek(otherID) // Peek or Get ?
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
