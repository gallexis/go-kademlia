package datastructure

import (
	"encoding/hex"
	"math/rand"
)

const BitsInNodeID = 160
const BytesInNodeID = BitsInNodeID / 8

type NodeID [BytesInNodeID]byte

type equality uint8

const (
	Equals equality = iota
	Greater
	Lower
)

func (n NodeID) compare(other NodeID) equality {
	for i := 0; i < BytesInNodeID; i++ {
		if n[i] > other[i] {
			return Greater

		} else if n[i] < other[i] {
			return Lower
		}
	}
	return Equals
}

func (n NodeID) IsGreaterThan(other NodeID) bool {
	return n.compare(other) == Greater
}

func (n NodeID) IsLowerThan(other NodeID) bool {
	return n.compare(other) == Lower
}

func (n NodeID) Equals(other NodeID) bool {
	return n.compare(other) == Equals
}

func (n NodeID) XOR(other NodeID) (newNodeID NodeID) {
	for i := 0; i < BytesInNodeID; i++ {
		newNodeID[i] = n[i] ^ other[i]
	}
	return
}

func NewNodeID() (n NodeID) {
	for i := 0; i < BytesInNodeID; i++ {
		n[i] = byte(rand.Intn(256))
	}
	return
}

func (_ NodeID) GetBucketNumber(xoredID NodeID) int {
	bytePosition := 0
	bitPosition := 0

	for bytePosition = 0; bytePosition < BytesInNodeID; bytePosition++ {
		if xoredID[bytePosition] != 0 {
			for bitPosition = 0; bitPosition < 8; bitPosition++ {
				if xoredID[bytePosition]&128 == 128 {
					return bytePosition*8 + bitPosition
				}
				xoredID[bytePosition] <<= 1
			}
		}
	}

	return -1
}

func (n NodeID) String() string {
	return hex.EncodeToString(n[:])
}

func (n NodeID) Encode() (b []byte) {
	return n[:]
}

func (n *NodeID) Decode(data []byte){
	for i := range n {
		n[i] = data[i]
	}
}

func FakeNodeID(id uint8) (n NodeID) {
	for i := 0; i < BytesInNodeID; i++ {
		n[i] = id
	}
	return
}