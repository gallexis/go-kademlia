package datastructure

import (
    "github.com/hashicorp/golang-lru"
    log "github.com/sirupsen/logrus"
    "math/rand"
    "sync"
    "time"
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

func (kb *KBucket) Insert(newContact Contact, pingNode func(chan bool)) {
    if !kb.isInBucket(newContact.NodeID) && kb.freeSpaceLeft() { // Insert
        kb.Contacts.Add(newContact.NodeID, newContact)

    } else if kb.isInBucket(newContact.NodeID) { // Update
        kb.Contacts.Remove(newContact.NodeID)
        kb.Contacts.Add(newContact.NodeID, newContact)

    } else { // Insert when full KB
        keys := kb.Contacts.Keys()
        oldestContact, _ := kb.Contacts.Peek(keys[0].(NodeID))
        oldestNodeID := oldestContact.(Contact).NodeID
        tick := time.Tick(5 * time.Second)
        pingChan := make(chan bool)

        pingNode(pingChan)

        select {
        case <-tick:
            kb.Contacts.Add(newContact.NodeID, newContact) //Add will remove oldestContact then add newContact
            //log.Println("add new")
        case <-pingChan:
            kb.Contacts.Add(oldestNodeID, oldestContact) // if the oldest answers, put it back to the tail
            log.Info("add old")
        }
    }
}

func (kb KBucket) GetRandomContacts(alpha int) []Contact {
    keys := kb.Contacts.Keys()
    contactsLength := len(keys)
    var contacts []Contact

    for i, key := range rand.Perm(contactsLength) {
        if i >= alpha{
            break
        }
        contact, _ := kb.Contacts.Peek(keys[key].(NodeID))
        contacts = append(contacts, contact.(Contact))
    }
    return contacts
}

func (kb KBucket) freeSpaceLeft() bool {
    return kb.Contacts.Len() < kb.K
}

func (kb KBucket) isInBucket(nid NodeID) bool {
    return kb.Contacts.Contains(nid)
}
