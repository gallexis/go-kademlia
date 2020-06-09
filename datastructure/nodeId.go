package datastructure

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
)

const BitsInNodeID = 160
const BytesInNodeID = BitsInNodeID / 8

type NodeId [BytesInNodeID]byte

type InfoHash = NodeId

type equality uint8

const (
	Equals equality = iota
	Greater
	Lower
)

func (n NodeId) compare(other NodeId) equality {
	for i := 0; i < BytesInNodeID; i++ {
		if n[i] > other[i] {
			return Greater

		} else if n[i] < other[i] {
			return Lower
		}
	}
	return Equals
}

func (n NodeId) IsGreaterThan(other NodeId) bool {
	return n.compare(other) == Greater
}

func (n NodeId) IsLowerThan(other NodeId) bool {
	return n.compare(other) == Lower
}

func (n NodeId) Equals(other NodeId) bool {
	return n.compare(other) == Equals
}

func (n NodeId) XOR(other NodeId) (newNodeID NodeId) {
	for i := 0; i < BytesInNodeID; i++ {
		newNodeID[i] = n[i] ^ other[i]
	}
	return
}

func NewNodeIdFromString(hash string) (n NodeId) {
	c, err := hex.DecodeString(hash)
	if err != nil {
		log.Fatal("NewNodeIdFromString", err.Error())
	}
	for i := 0; i < BytesInNodeID; i++ {
		n[i] = c[i]
	}
	return
}

func NewNodeID() (n NodeId) {
	for i := 0; i < BytesInNodeID; i++ {
		n[i] = byte(rand.Intn(256))
	}
	return
}

func FakeNodeID(id uint8) (n NodeId) {
	for i := 0; i < BytesInNodeID; i++ {
		n[i] = id
	}
	return
}

func (_ NodeId) GetBucketNumber(xoredID NodeId) BucketPosition {
	bytePosition := 0
	bitPosition := 0

	for bytePosition = 0; bytePosition < BytesInNodeID; bytePosition++ {
		if xoredID[bytePosition] != 0 {
			for bitPosition = 0; bitPosition < 8; bitPosition++ {
				if xoredID[bytePosition]&128 == 128 {
					return BucketPosition(bytePosition*8 + bitPosition)
				}
				xoredID[bytePosition] <<= 1
			}
		}
	}

	return BitsInNodeID - 1
}

func (n NodeId) String() string {
	return hex.EncodeToString(n[:])
}

func (n NodeId) Encode() (b []byte) {
	return n[:]
}

func (n *NodeId) Decode(data []byte) {
	if len(data) < len(n) {
		panic(fmt.Errorf("ERROR PING DECODE [%s]", data))
	}
	for i := range n {
		n[i] = data[i] // todo : Unsafe
	}
}
