package test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	ds "kademlia/datastructure"
)

// 0b1000010 = 0x42
// 0b1101101 = 0x6D

func fakeRandomNodeID() ds.NodeId {
	var n ds.NodeId
	for i := 0; i < ds.BytesInNodeID; i++ {
		n[i] = uint8(rand.Intn(255))
	}
	return n
}

func newRandomContact() ds.Contact {
	return ds.Contact{
		IP:     "12.34.56.78",
		Port:   1337,
		NodeID: fakeRandomNodeID(),
	}
}

func TestNewKBucket(t *testing.T) {
	k := 10
	kb := ds.NewKBucket(k)

	if kb.K != k {
		t.Error("K is incorrect.")
	}

	if kb.Nodes.Len() != 0 {
		t.Error("Nodes' length must be 0.")
	}

}

func TestNewKBucket_Fails(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		ds.NewKBucket(0)
		return
	}

	// Start the actual test in a different subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestNewKBucket_Fails")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	stdout, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Check that the log fatal message is what we expected
	gotBytes, _ := ioutil.ReadAll(stdout)
	got := string(gotBytes)
	expected := "Must provide a positive size"
	if !strings.HasSuffix(got[:len(got)-1], expected) {
		t.Fatalf("Unexpected log message. Got %s but should be %s", got, expected)
	}

	// Check that the program exited
	err := cmd.Wait()
	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatalf("Process ran with err %v, want exit status 1", err)
	}

}

func TestInsertKB_full(t *testing.T) {
	k := 2
	kb := ds.NewKBucket(k)
	c1 := newRandomContact()
	c2 := newRandomContact()
	c3 := newRandomContact()

	kb.Insert(c1)
	kb.Insert(c2)
	kb.Insert(c3)

	if kb.Nodes.Len() != k {
		t.Error("Nodes length must be ", k)
	}

	// assert node2 & node3 are the only one in the list
	nodeIDs := kb.Nodes.Keys()
	node2 := nodeIDs[0].(ds.NodeId)
	node3 := nodeIDs[1].(ds.NodeId)

	if c2.NodeID != node2 {
		t.Error("c2.nodeID != node2", c2.NodeID, node2)
	}

	if c3.NodeID != node3 {
		t.Error("c3.nodeID != node3")
	}
}

func TestInsertKB_notFull(t *testing.T) {
	k := 2
	kb := ds.NewKBucket(k)
	c1 := newRandomContact()
	c2 := newRandomContact()

	if kb.Nodes.Len() != 0 {
		t.Error("Nodes length must be 0.")
	}

	kb.Insert(c1)
	if kb.Nodes.Len() != 1 {
		t.Error("Nodes length must be 1.")
	}

	kb.Insert(c2)
	if kb.Nodes.Len() != 2 {
		t.Error("Nodes length must be 2.")
	}

	// test If c1 & C2 correctly inserted
	nodes := kb.Nodes.Keys()
	node1 := nodes[0].(ds.NodeId)
	node2 := nodes[1].(ds.NodeId)

	if c1.NodeID != node1 {
		t.Error("c1.NodeId != node1")
	}

	if c2.NodeID != node2 {
		t.Error("c2.NodeId != node2")
	}
}

func TestInsertKB_update(t *testing.T) {
	k := 3
	kb := ds.NewKBucket(k)
	c1 := newRandomContact()
	c2 := newRandomContact()

	if kb.Nodes.Len() != 0 {
		t.Error("Nodes length must be 0.")
	}

	kb.Insert(c1)
	kb.Insert(c2)
	kb.Insert(c1) // reinsert c1

	if kb.Nodes.Len() != 2 {
		t.Error("Nodes length must be 2.")
	}

	nodes := kb.Nodes.Keys()
	node1 := nodes[1].(ds.NodeId)
	node2 := nodes[0].(ds.NodeId)

	if c1.NodeID != node1 {
		t.Error("c1.NodeId != node1")
	}

	if c2.NodeID != node2 {
		t.Error("c2.NodeId != node2")
	}
}
