package krpc

import (
    "encoding/hex"
    "fmt"
    "math/rand"
)

// GENERIC RECEIVE
type Response struct {
    Values []string `bencode:"values"`
    Id     []byte   `bencode:"id"`
    Nodes  []byte   `bencode:"nodes"`
    Nodes6 []byte   `bencode:"nodes6"`
    Token  string   `bencode:"token"`
}

type Answer struct {
    Id          []byte `bencode:"id"`
    Target      []byte `bencode:"target"`
    InfoHash    []byte `bencode:"info_hash"`
    Port        int    `bencode:"port"`
    Token       string `bencode:"token"`
    ImpliedPort int    `bencode:"implied_port"`
}

type GenericMessage struct {
    T string   `bencode:"t"`
    Y string   `bencode:"y"`
    Q string   `bencode:"q"`
    R Response `bencode:"r"`
    A Answer   `bencode:"a"`
    E []string `bencode:"e"`
}

// GENERIC SEND
type RequestMessage struct {
    T string            `bencode:"t"`
    Y string                 `bencode:"y"`
    Q string                 `bencode:"q"`
    A map[string]interface{} `bencode:"a"`
}

type ResponseMessage struct {
    T string            `bencode:"t"`
    Y string                 `bencode:"y"`
    R map[string]interface{} `bencode:"r"`
}

type RandomBytes []byte

func (t RandomBytes) String() string {
    return hex.EncodeToString([]byte(t))
}

func NewRandomBytes(n int) RandomBytes {
    token := make([]byte, n)
    if _, err := rand.Read(token); err != nil {
        fmt.Printf("Failed to generate NewRandomBytes: %v", err)
    }
    return token
}

func NewRandomBytesFromString(token string) RandomBytes {
    t := RandomBytes{}
    t, err := hex.DecodeString(token)
    if err != nil {
        fmt.Printf("Error when decoding from string: %v", err)
    }
    return t
}
