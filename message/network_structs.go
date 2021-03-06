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
	Values []interface{} `bencode:"values"`
	Id     []byte        `bencode:"id"`
	Nodes  []byte        `bencode:"nodes"`
	Nodes6 []byte        `bencode:"nodes6"`
	Token  []byte        `bencode:"token"`
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
		log.Panic("MessageToBytes ->", err)
	}

	return buffer
}

func BytesToMessage(data []byte) (g GenericMessage, ok bool) {
	decoder := bencode.NewDecoder(bytes.NewBuffer(data))
	if err := decoder.Decode(&g); err != nil {
		log.Error("BytesToMessage ", err)
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
	tx := make([]byte, 10)
	if _, err := rand.Read(tx); err != nil {
		log.Panicf("Failed to generate NewTransactionId: %v", err)
	}
	return tx
}

func NewTransactionIdFromString(tx string) TransactionId {
	t := TransactionId{}
	t, err := hex.DecodeString(tx)
	if err != nil {
		return TransactionId(tx)
	}
	return t
}

type Token []byte

func (t Token) Equals(other Token) bool {
	if len(t) != len(other) {
		return false
	}

	for i := range other {
		if t[i] != other[i] {
			return false
		}
	}
	return true
}

func (t Token) String() string {
	return hex.EncodeToString([]byte(t))
}

func NewToken() Token {
	tx := make([]byte, 10)
	if _, err := rand.Read(tx); err != nil {
		log.Panicf("Failed to generate NewTransactionId: %v", err)
	}
	return tx
}

type Message interface {
	Encode() []byte
	Decode(GenericMessage)
}
