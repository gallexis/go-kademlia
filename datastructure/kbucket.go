package datastructure

import (
	"github.com/hashicorp/golang-lru"
	"log"
	"sync"
)

type KBucket struct {
	Contacts *lru.Cache
	K        int
	mutex    *sync.Mutex
}

func NewKBucket(k int) KBucket {
	contacts, err := lru.New(k)
	if err != nil {
		log.Fatalln(err)
	}
	return KBucket{
		Contacts: contacts,
		K:        k,
		mutex:    &sync.Mutex{},
	}
}

func (kb *KBucket) Insert(newContact Contact) {
	if !kb.isInBucket(newContact.NodeID) && kb.freeSpaceLeft() { // Insert
		kb.Contacts.Add(newContact.NodeID, newContact)

	} else if kb.isInBucket(newContact.NodeID) { // Update
		kb.Contacts.Remove(newContact.NodeID)
		kb.Contacts.Add(newContact.NodeID, newContact)

	} else { // Insert when full KB
		keys := kb.Contacts.Keys()
		oldestNodeID := keys[0].(NodeID)
		oldestContact, _ := kb.Contacts.Peek(oldestNodeID)

		if true { // TODO : if oldestContact doesn't reply
			kb.Contacts.Add(newContact.NodeID, newContact) //Add will remove oldestContact then add newContact
		} else {
			kb.Contacts.Add(oldestNodeID, oldestContact) // if the oldest answers, put it back to the tail
		}
	}
}

func (kb KBucket) freeSpaceLeft() bool {
	return kb.Contacts.Len() < kb.K
}

func (kb KBucket) isInBucket(nid NodeID) bool {
	return kb.Contacts.Contains(nid)
}
