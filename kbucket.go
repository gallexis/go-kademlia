package main

import (
    "sync"
    "fmt"
)

type KBucket struct{
    Contacts []Contact
    K        int
    mutex    sync.Mutex
}

func InitKBucket(k int) KBucket{
    return KBucket{K:k}
}

func (kb KBucket) indexInBucketOf(contact Contact) int{
    for i, c := range kb.Contacts {
        if c.Node.Equals(contact.Node){
            return i
        }
    }
    return -1
}

func (kb KBucket) isFull() bool{
    return len(kb.Contacts) >= kb.K
}

func (kb *KBucket) removeAt(i int){
    kb.Contacts = append(kb.Contacts[:i], kb.Contacts[i+1:]...)
}

func (kb KBucket) Display(){
    fmt.Println(kb.isFull(), len(kb.Contacts))
}

func (kb *KBucket) Insert(contact Contact) bool{
    kb.mutex.Lock()
    defer kb.mutex.Unlock()

    if kb.isFull(){
        return false
    }

    kb.Contacts = append(kb.Contacts, contact)
    return true
}

func (kb *KBucket) Update(contact Contact) bool{
    kb.mutex.Lock()

    if i := kb.indexInBucketOf(contact); i > -1{
        c := kb.Contacts[i]
        kb.removeAt(i)

        kb.mutex.Unlock()
        return kb.Insert(c)
    }

    kb.mutex.Unlock()
    return kb.Insert(contact)
}
