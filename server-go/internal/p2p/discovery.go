package p2p

import (
	"fmt"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
// 节点发现的通告结构体，继承 Notifee 接口
type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
// 继承函数，节点发现后的处理函数：自动链接节点
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.String())
	n.PeerChan <- pi
	//err := n.h.Connect(context.Background(), pi)
	//if err != nil {
	//	fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	//}
}

// Initialize the MDNS service
func InitMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}
