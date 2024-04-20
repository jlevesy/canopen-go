package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/jlevesy/canopen-go"
	"github.com/jlevesy/canopen-go/heartbeat"
	"github.com/jlevesy/canopen-go/nmt"
	"go.einride.tech/can/pkg/socketcan"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := socketcan.DialContext(ctx, "can", "vcan0")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	net := canopen.NewNetwork(conn)

	stateManager := nmt.NewStateManager(net)
	heartbeatSender := heartbeat.NewSender(net)

	net.Run()

	fmt.Println("Setting all nodes to PRE-OPERATIONAL")
	stateManager.SetNetworkState(ctx, nmt.AllNodes, nmt.StatePreOperational)

	defer func() {
		fmt.Println("Setting all nodes to OPERATIONAL")
		stateManager.SetNetworkState(context.WithoutCancel(ctx), nmt.AllNodes, nmt.StateOperational)
	}()

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("bye bye")
			return
		case <-t.C:
			if err := heartbeatSender.SendHeartbeat(
				ctx,
				heartbeat.Heartbeat{NodeID: 92, State: heartbeat.StateOperational},
			); err != nil {
				panic(err)
			}
		}

	}

}
