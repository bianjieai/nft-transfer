package types

import (
	fmt "fmt"

	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
)

const (

	// NativePortID corresponds to the cosmos native module
	NativePortID = "nft-transfer"
	// ERC721PortID corresponds to the Ethereum erc721 protocol
	ERC721PortID = "erc721-transfer"
)

// The router is a map from port name to the NFTKeeper
// which is to support A single module can bind to multiple ports at once
type Router struct {
	sealed bool
	routes map[string]NFTKeeper
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]NFTKeeper),
	}
}

// AddRoute adds NFTKeeper for a given port. It returns the Router
// so AddRoute calls can be linked. It will panic if the Router is sealed.
func (r *Router) AddRoute(port string, nftKeeper NFTKeeper) *Router {
	if r.sealed {
		panic("router already sealed")
	}

	if r.HasRoute(port) {
		panic(fmt.Sprintf("route %s has already been registered", port))
	}

	if err := host.PortIdentifierValidator(port); err != nil {
		panic(err)
	}
	r.routes[port] = nftKeeper
	return r
}

// GetRoute returns a NFTKeeper for a given port.
func (r *Router) GetRoute(port string) (NFTKeeper, bool) {
	if !r.HasRoute(port) {
		return nil, false
	}
	return r.routes[port], true
}

// HasRoute returns true if the Router has registered or false otherwise.
func (r *Router) HasRoute(port string) bool {
	_, ok := r.routes[port]
	return ok
}

// Seal prevents the Router from any subsequent route handlers to be registered.
// Seal will panic if called more than once.
func (r *Router) Seal() {
	if r.sealed {
		panic("router already sealed")
	}
	r.sealed = true
}
