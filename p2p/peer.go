package p2p

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
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
