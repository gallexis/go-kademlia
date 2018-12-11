package test

import (
	ds "kademlia/datastructure"
	"testing"
)

// 0b10001101 = 0x8d
var selfID = ds.NodeID{
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
	0x8d,
}

// 0b10101100 = 0xac
// 0x8d ^ 0xac = 0b00100000 (0x20)

// 0b10101101 = 0xad
// 0x8d ^ 0xad = 0b00100001 (0x21)
func fakeContact(position int, value byte) ds.Contact {
	nid := selfID
	nid[position] = value

	return ds.Contact{
		IP:     "",
		Port:   0,
		NodeID: nid,
	}
}

func TestNewRoutingTable(t *testing.T) {
	rt := ds.NewRoutingTable(ds.FakeNodeID(0x8d))

	for i := 0; i < len(rt.KBuckets); i++ {
		if rt.KBuckets[i].K != ds.K || rt.KBuckets[i].Contacts.Len() != 0 {
			t.Error("problem with K or kbuckets' length")
		}
	}
	if rt.K != ds.K {
		t.Error("problem with K")
	}
	if rt.Alpha != ds.Alpha {
		t.Error("problem with Alpha")
	}
}

func TestRoutingTable_Insert(t *testing.T) {
	k := 2
	rt := ds.NewRoutingTableWithDetails(ds.FakeNodeID(0x8d), k, ds.Alpha)
	contact1 := fakeContact(2, 0xad)
	contact2 := fakeContact(2, 0xac)
	contact3 := fakeContact(2, 0xa0)
	expectedContactPositionInKBucket := 18

	rt.Insert(contact1)
	rt.Insert(contact2)
	rt.Insert(contact3)

	kb := rt.KBuckets[expectedContactPositionInKBucket].Contacts.Keys()

	if kb[0] != contact2.NodeID || kb[1] != contact3.NodeID || len(kb) != 2 {
		t.Error("problem in Insert")
	}
}

func TestRoutingTable_GetOne(t *testing.T) {
	k := 2
	rt := ds.NewRoutingTableWithDetails(ds.FakeNodeID(0x8d), k, ds.Alpha)
	contact1 := fakeContact(2, 0xad)
	contact2 := fakeContact(2, 0xac)

	rt.Insert(contact1)
	rt.Insert(contact2)

	contact, ok := rt.GetOne(contact1.NodeID)

	if !ok || contact.NodeID != contact1.NodeID {
		t.Error("problem in GetOne")
	}
}

func TestRoutingTable_GetOne_Fail(t *testing.T) {
	k := 2
	rt := ds.NewRoutingTableWithDetails(ds.FakeNodeID(0x8d), k, ds.Alpha)
	contact1 := fakeContact(2, 0xad)
	contact2 := fakeContact(2, 0xac)

	rt.Insert(contact1)

	contact, exists := rt.GetOne(contact2.NodeID)
	if exists || (contact != ds.Contact{}) {
		t.Error("problem in GetOne")
	}
}

func TestRoutingTable_Get_withFullKB(t *testing.T) {
	k := 2
	rt := ds.NewRoutingTableWithDetails(ds.FakeNodeID(0x8d), k, ds.Alpha)
	contact1 := fakeContact(2, 0xad)
	contact2 := fakeContact(2, 0xac)
	contact3 := fakeContact(0, 0xFF)
	contact4 := fakeContact(12, 0xFF)

	rt.Insert(contact1)
	rt.Insert(contact2)
	rt.Insert(contact3)
	rt.Insert(contact4)

	kb := rt.Get(contact1.NodeID)

	if len(kb) != k || kb[0] != contact1 || kb[1] != contact2 {
		t.Error("problem in Get")
	}
}

func TestRoutingTable_Get_withoutFullKB(t *testing.T) {
	k := 5
	rt := ds.NewRoutingTableWithDetails(ds.FakeNodeID(0x8d), k, ds.Alpha)
	contact1 := fakeContact(2, 0xad)
	contact2 := fakeContact(2, 0xac)
	contact3 := fakeContact(0, 0xFF)
	contact4 := fakeContact(12, 0xFF)
	contact5 := fakeContact(0, 0x00)

	rt.Insert(contact1)
	rt.Insert(contact2)
	rt.Insert(contact3)
	rt.Insert(contact4)
	rt.Insert(contact5)

	kb := rt.Get(contact1.NodeID)

	if len(kb) != k ||
		kb[0] != contact1 ||
		kb[1] != contact2 ||
		kb[2] != contact3 ||
		kb[3] != contact5 ||
		kb[4] != contact4 {
		t.Error("problem in Get")
	}
}

func TestRoutingTable_Get_FullTotalKB(t *testing.T) {
	k := 4
	rt := ds.NewRoutingTableWithDetails(ds.FakeNodeID(0x8d), k, ds.Alpha)
	contact1 := fakeContact(2, 0xad)
	contact2 := fakeContact(2, 0xac)
	contact3 := fakeContact(0, 0xFF)
	contact4 := fakeContact(12, 0xFF)
	contact5 := fakeContact(0, 0x00)

	rt.Insert(contact1)
	rt.Insert(contact2)
	rt.Insert(contact3)
	rt.Insert(contact4)
	rt.Insert(contact5)

	kb := rt.Get(contact1.NodeID)

	if len(kb) != k ||
		kb[0] != contact1 ||
		kb[1] != contact2 ||
		kb[2] != contact3 ||
		kb[3] != contact5 {
		t.Error("problem in Get")
	}
}
