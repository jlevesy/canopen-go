package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/jlevesy/canopen-go"
	"github.com/jlevesy/canopen-go/heartbeat"
	"github.com/jlevesy/canopen-go/nmt"
	"go.einride.tech/can/pkg/socketcan"
)

func main() {
	fmt.Println("Starting the listener")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := socketcan.DialContext(ctx, "can", "vcan0")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	net := canopen.NewNetwork(conn)

	stateReceiver := nmt.NewStateReceiver(nmt.StateChangeHandlerFunc(func(nodeID uint8, state nmt.State) {
		fmt.Println("Received a network state change", nodeID, state)
	}))

	if err := net.RegisterReceiver(stateReceiver); err != nil {
		panic(err)
	}

	heartBeatReceiver := heartbeat.NewGlobalReceiver(heartbeat.HandlerFunc(func(b heartbeat.Heartbeat) {
		fmt.Println("Received an heartbeat from node", b.NodeID, "Status is", b.State)
	}))

	if err := net.RegisterReceiver(heartBeatReceiver); err != nil {
		panic(err)
	}

	net.Run()

	fmt.Println("Listener is ready")

	<-ctx.Done()

}
