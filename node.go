package main

import (
    "math/rand"
)

const LengthNode = 20
type Node [LengthNode]byte

// Define enum for equality
type equality int

const (
    equals equality = iota
    greater
    lower
)

func (n Node) compare(other Node) equality {
    for i := 0; i < LengthNode; i++{
        if n[i] > other[i]{
            return greater

        } else if n[i] < other[i]{
            return lower
        }
    }
    return equals
}

func (n Node) GreaterThan(other Node) bool{
    return n.compare(other) == greater
}

func (n Node) LowerThan(other Node) bool{
    return n.compare(other) == lower
}

func (n Node) Equals(other Node) bool{
    return n.compare(other) == equals
}

func (n Node) Xor(other Node) (newNode Node) {
    for i := 0; i < LengthNode; i++ {
        newNode[i] = n[i] ^ other[i]
    }
    return newNode
}

func (n Node) getPosition() int{
    for i:=0; i<LengthNode; i++ {
        for j:=0; j<LengthNode; j++ {
            if((n[i] >> uint8(7-j)) & 0x1) == 1{
                return i * 8 + j
            }
        }
    }
    return LengthNode * 8 - 1
}

func NewNode() (n Node){
    for i:=0; i<LengthNode; i++{
        n[i] = uint8(rand.Intn(256))
    }
    return n

    /*
    b := make([]byte, 20)
    if _, err := rand.Read(b); err != nil {
        log.Fatalln("nodeId rand:", err)
    }
    n = b
    */
}