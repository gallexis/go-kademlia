package krpc

import (
    "bytes"
    "github.com/zeebo/bencode"
    "log"
)

func MessageToBytes(message interface{}) []byte{
    buffer, err := bencode.EncodeBytes(message)
    if err != nil {
        panic(err)
    }

    if err != nil {
        log.Println(err.Error())
    }

    return buffer
}

func BytesToMessage(data []byte) (g GenericMessage){
    decoder := bencode.NewDecoder(bytes.NewBuffer(data))
    if err := decoder.Decode(&g); err != nil {
        panic(err)
    }
    return
}

