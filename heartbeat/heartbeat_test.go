package heartbeat_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/jlevesy/canopen-go/heartbeat"
	"github.com/jlevesy/canopen-go/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeReceiver(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	senderNet, receiverNet, done := testutil.SetupSenderReceiverNetworks(t)
	defer done()

	var (
		handler  = heartbeatHandler{beats: make(chan heartbeat.Heartbeat, 1)}
		sender   = heartbeat.NewSender(senderNet)
		receiver = heartbeat.NewNodeReceiver(10, &handler)
	)

	err := receiverNet.RegisterReceiver(receiver)
	require.NoError(t, err)

	senderNet.Run()
	receiverNet.Run()

	// Send two heartbeats from two different node IDs, the received should only care about one.
	err = sender.SendHeartbeat(ctx, heartbeat.Heartbeat{NodeID: 12, State: heartbeat.StateInitialising})
	require.NoError(t, err)

	err = sender.SendHeartbeat(ctx, heartbeat.Heartbeat{NodeID: 10, State: heartbeat.StateOperational})
	require.NoError(t, err)

	select {
	case <-ctx.Done():
		t.FailNow()
	case gotBeat := <-handler.beats:
		assert.Equal(t, uint8(10), gotBeat.NodeID)
		assert.Equal(t, heartbeat.StateOperational, gotBeat.State)
	}

	// Try a second time and read from the channel.
	// If we can it means that we handled a beat that we shouldn't have.
	select {
	case gotBeat := <-handler.beats:
		t.Fatal("Unexpected beat received", gotBeat)
	default:
	}
}

func TestGlobalReceiver(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	senderNet, receiverNet, done := testutil.SetupSenderReceiverNetworks(t)
	defer done()

	var (
		handler  = heartbeatHandler{beats: make(chan heartbeat.Heartbeat, 2)}
		sender   = heartbeat.NewSender(senderNet)
		receiver = heartbeat.NewGlobalReceiver(&handler)
	)

	err := receiverNet.RegisterReceiver(receiver)
	require.NoError(t, err)

	senderNet.Run()
	receiverNet.Run()

	var (
		gotBeats  []heartbeat.Heartbeat
		wantBeats = []heartbeat.Heartbeat{
			heartbeat.Heartbeat{NodeID: 10, State: heartbeat.StateOperational},
			heartbeat.Heartbeat{NodeID: 12, State: heartbeat.StateInitialising},
		}
	)

	for _, b := range wantBeats {
		err = sender.SendHeartbeat(ctx, b)
		require.NoError(t, err)
	}

	for i := 0; i < len(wantBeats); i++ {
		select {
		case <-ctx.Done():
			t.FailNow()
		case gotBeat := <-handler.beats:
			gotBeats = append(gotBeats, gotBeat)
		}
	}

	select {
	case gotBeat := <-handler.beats:
		t.Fatal("Unexpected beat received", gotBeat)
	default:
	}

	sort.Slice(gotBeats, func(i, j int) bool { return gotBeats[i].NodeID < gotBeats[j].NodeID })
	assert.Equal(t, wantBeats, gotBeats)
}

type heartbeatHandler struct {
	beats chan heartbeat.Heartbeat
}

func (r *heartbeatHandler) HandleHeartbeat(b heartbeat.Heartbeat) {
	r.beats <- b
}
