package datastructure

import (
	"fmt"
	"testing"
)

func FakeNodeID(id uint8) (n NodeID) {
	for i := 0; i < BytesInNodeID; i++ {
		n[i] = id
	}
	return
}

func TestNewNodeID(t *testing.T) {
	if NewNodeID() == NewNodeID() {
		t.Error("Very suspicious if the 2 nodes ID are identical")
	}
}

func TestNodeID_Compare(t *testing.T) {
	n1 := FakeNodeID(100)
	n2 := FakeNodeID(101)

	if !n1.IsLowerThan(n2) {
		t.Error("n1 is lower than n2")
	}

	if !n2.IsGreaterThan(n1) {
		t.Error("n2 is greater than n1")
	}

	n1 = FakeNodeID(100)
	n2 = FakeNodeID(100)

	if !n1.Equals(n2) {
		t.Error("n1 is equals to n2")
	}
}

func TestNodeID_XOR(t *testing.T) {
	n1 := FakeNodeID(100)
	n2 := FakeNodeID(100)
	xoredID := n1.XOR(n2)

	for i := 0; i < BytesInNodeID; i++ {
		if xoredID[i] != 0 {
			t.Error("xoredID[i] must be 0")
		}
	}
}

// 0b00001101 = 0x0d
// 0b10001101 = 0x8d
func TestGet_getBucketNumber(t *testing.T) {
	n1 := FakeNodeID(0x0d)
	pos := n1.getBucketNumber(n1)

	if pos != 4 {
		t.Error("should be 4 : ", pos)
	}

	n1 = FakeNodeID(0x8d)
	pos = n1.getBucketNumber(n1)

	if pos != 0 {
		t.Error("should be 0")
	}

	n1 = FakeNodeID(0)
	pos = n1.getBucketNumber(n1)

	if pos != -1 {
		t.Error("should be -1")
	}
}

func TestNodeID_String(t *testing.T) {
	n := FakeNodeID(0xff)
	n[0] = 0x00
	toString := fmt.Sprint(n)
	if toString != "00ffffffffffffffffffffffffffffffffffffff" {
		t.Error("Error when converting node ID to string")
	}
}
