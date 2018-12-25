package datastructure


type PeerInfo struct {
    IP     string
    Port   string
}

type PeerStore map[string][]PeerInfo

func (p *PeerStore) Add(infoHash, ip, port string){
    //p[infoHash] = PeerInfo{}

}