package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

// Create a p2p host that listen on the specified port
func MakeNewHost(ctx context.Context, port int) (host.Host, error) {
	// Generate a key pair for obtaining a valid host ID
	r := rand.Reader
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// Set options and create new host
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
		libp2p.Identity(priv),
	}
	host, err := libp2p.New(ctx, opts...)

	return host, err
}

func Connect(ctx context.Context, host host.Host, targetPeer string) peer.ID {
	// Extract target address from peer ID
	/*ipfsAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", targetPeer))
	if err != nil {
		log.Println("ipfsAddr", err)
	}

	targetAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", targetPeer))
	if err != nil {
		log.Println("targetAddr", err)
	}

	addr := ipfsAddr.Decapsulate(targetAddr)*/

	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/30000/ipfs/%s", targetPeer))
	fmt.Println(ipfsaddr)
	if err != nil {
		log.Println(err)
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	fmt.Println(pid)
	if err != nil {
		log.Println(err)
	}

	peerId, err := peer.Decode(pid)
	fmt.Println(peerId)
	if err != nil {
		log.Println(err)
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", pid))
	fmt.Println(targetPeerAddr)
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
	fmt.Println(targetAddr)

	// Add target address and peer id to the peerstore (like a phonebook)
	host.Peerstore().AddAddr(peerId, targetAddr, peerstore.PermanentAddrTTL)

	// Add target address and peer id to the peerstore (like a phonebook)
	host.Peerstore().AddAddr(peer.ID(targetPeer), targetAddr, peerstore.PermanentAddrTTL)

	return peerId
}
