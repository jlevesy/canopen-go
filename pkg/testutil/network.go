package testutil

import (
	"os"
	"testing"

	"github.com/jlevesy/canopen-go"
	"github.com/stretchr/testify/require"
	"go.einride.tech/can/pkg/socketcan"
)

func SetupNetwork(t *testing.T) (*canopen.Network, func()) {
	t.Helper()

	conn, err := socketcan.Dial("can", getTestDevice())
	require.NoError(t, err)

	return canopen.NewNetwork(conn), func() { _ = conn.Close() }
}

// This helper sets up two different networks (which means it opens two different socket over the same device.
// It is required to do so in testing because the socket that emits a frame on the bus doesn't receive it.
func SetupSenderReceiverNetworks(t *testing.T) (*canopen.Network, *canopen.Network, func()) {
	senderNet, doneSender := SetupNetwork(t)
	receiverNet, doneReceiver := SetupNetwork(t)

	return senderNet, receiverNet, func() {
		doneReceiver()
		doneSender()
	}
}

func getTestDevice() string {
	if dev := os.Getenv("CANOPEN_TEST_DEVICE"); dev != "" {
		return dev
	}

	return "vcan0"
}
