package datastructure

import (
	"encoding/binary"
	"net"
)

type Contact struct {
	IP     net.IP
	Port   uint16
	NodeID NodeId
}

func (c *Contact) Decode(data []byte){
	c.NodeID.Decode(data[:20])
	c.IP = net.IP(data[20:24])
	c.Port = binary.BigEndian.Uint16(data[24:26])
}

func (c *Contact) Encode() []byte{
	b := c.NodeID.Encode()
	b = append(b, c.IP.To4()...)
	binary.BigEndian.PutUint16(b, c.Port)
	return b
}

type PeerContact struct {
	IP     net.IP
	Port   uint16
}

func (c *PeerContact) Decode(data []byte){
	c.IP = net.IP(data[:4])
	c.Port = binary.BigEndian.Uint16(data[4:])
}

func (c *PeerContact) Encode() (b []byte){
	b = append(b, c.IP.To4()...)
	binary.BigEndian.PutUint16(b, c.Port)
	return
}

