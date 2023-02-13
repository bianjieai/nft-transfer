package ics721testing

import (
	"testing"
	"time"

	ibctesting "github.com/cosmos/ibc-go/v7/testing"
)

type Coordinator struct {
	*ibctesting.Coordinator
}

// NewCoordinator initializes Coordinator with N TestChain's
func NewCoordinator(t *testing.T, n int) *Coordinator {
	chains := make(map[string]*ibctesting.TestChain)
	coord := &ibctesting.Coordinator{
		T:           t,
		CurrentTime: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	for i := 1; i <= n; i++ {
		chainID := ibctesting.GetChainID(i)
		chains[chainID] = NewTestChain(t, coord, chainID)
	}
	coord.Chains = chains

	return &Coordinator{coord}
}
