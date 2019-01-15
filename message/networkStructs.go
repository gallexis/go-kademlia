package message

import (
    "bytes"
    "encoding/hex"
    log "github.com/sirupsen/logrus"
    "github.com/zeebo/bencode"
    "math/rand"
)

// GENERIC RECEIVE
type Response struct {
    Values []string `bencode:"values"`
    Id     []byte   `bencode:"id"`
    Nodes  []byte   `bencode:"nodes"`
    Nodes6 []byte   `bencode:"nodes6"`
    Token  []byte   `bencode:"token"`
}

type Answer struct {
    Id          []byte `bencode:"id"`
    Target      []byte `bencode:"target"`
    InfoHash    []byte `bencode:"info_hash"`
    Port        int    `bencode:"port"`
    Token       []byte `bencode:"token"`
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
    T string                 `bencode:"t"`
    Y string                 `bencode:"y"`
    Q string                 `bencode:"q"`
    A map[string]interface{} `bencode:"a"`
}

type ResponseMessage struct {
    T string                 `bencode:"t"`
    Y string                 `bencode:"y"`
    R map[string]interface{} `bencode:"r"`
}

func MessageToBytes(message interface{}) []byte {
    buffer, err := bencode.EncodeBytes(message)
    if err != nil {
        log.Panic(err)
    }

    return buffer
}

func BytesToMessage(data []byte) (g GenericMessage, ok bool) {
    decoder := bencode.NewDecoder(bytes.NewBuffer(data))
    if err := decoder.Decode(&g); err != nil {
        log.Error(err)
        return
    }
    ok = true
    return
}

type TransactionId []byte

func (t TransactionId) String() string {
    return hex.EncodeToString([]byte(t))
}

func NewTransactionId() TransactionId {
    tx := make([]byte, 2)
    if _, err := rand.Read(tx); err != nil {
        log.Panicf("Failed to generate NewTransactionId: %v", err)
    }
    return tx
}

func NewTransactionIdFromString(tx string) TransactionId {
    t := TransactionId{}
    t, err := hex.DecodeString(tx)
    if err != nil {
        log.Errorf("Error when decoding from string: %v | %v | %v", err, tx, []byte(tx))
        return []byte(tx)
    }
    return t
}

type Token = TransactionId


type Message interface {
    Encode() []byte
    Decode(GenericMessage)
}
