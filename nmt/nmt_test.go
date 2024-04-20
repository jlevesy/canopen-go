package nmt_test

import (
	"context"
	"testing"
	"time"

	"github.com/jlevesy/canopen-go/nmt"
	"github.com/jlevesy/canopen-go/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNMT(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	senderNet, receiverNet, done := testutil.SetupSenderReceiverNetworks(t)
	defer done()

	var (
		mgr = nmt.NewStateManager(senderNet)

		stateChanges = make(chan stateChange)
		rec          = nmt.NewStateReceiver(nmt.StateChangeHandlerFunc(func(nodeID uint8, state nmt.State) {
			stateChanges <- stateChange{nodeID: nodeID, state: state}
		}))
	)

	err := receiverNet.RegisterReceiver(rec)
	require.NoError(t, err)

	senderNet.Run()
	receiverNet.Run()

	err = mgr.SetNetworkState(ctx, 92, nmt.StatePreOperational)
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		t.FailNow()
	case gotChange := <-stateChanges:
		assert.Equal(t, uint8(92), gotChange.nodeID)
		assert.Equal(t, nmt.StatePreOperational, gotChange.state)
	}
}

type stateChange struct {
	nodeID uint8
	state  nmt.State
}
