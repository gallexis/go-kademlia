package main

type RoutingTable struct {
    Kbuckets []KBucket
    node     Node
}

func InitRoutingTable(bucketNumber int, K int) RoutingTable{
    var kBuckets []KBucket

    for i:=0; i<bucketNumber;i++{
        kBuckets = append(kBuckets,InitKBucket(K))
    }

    return RoutingTable{Kbuckets: kBuckets, node:NewNode()}
}

func (rt *RoutingTable) Insert(c Contact) bool {
    i := rt.node.Xor(c.Node).getPosition()
    return rt.Kbuckets[i].Insert(c)
}

func (rt *RoutingTable) Update(c Contact) bool {
    i := rt.node.Xor(c.Node).getPosition()
    return rt.Kbuckets[i].Update(c)
}

func (rt RoutingTable) Display(){
    for _, k := range rt.Kbuckets {
        k.Display()
    }
}