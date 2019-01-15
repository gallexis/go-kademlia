package datastructure

type PeerStore map[string][]Node

func (p *PeerStore) Add(infoHash InfoHash, nodes []Node) {
    v, exists := (*p)[infoHash.String()]
    if exists {
        v = append(v, nodes...)
        (*p)[infoHash.String()] = v
    } else {
        (*p)[infoHash.String()] = nodes
    }
}

func (p *PeerStore) Contains(infoHash InfoHash) bool {
    _, exists := (*p)[infoHash.String()]
    return exists
}
