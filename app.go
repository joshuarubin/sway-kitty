package main

import (
	"context"
	"flag"

	"github.com/joshuarubin/go-sway"
)

type app struct {
	socketPath string
	client     sway.Client
}

func newApp() *app {
	var a app
	flag.StringVar(&a.socketPath, "socketpath", "", "Use the specified socket path")
	return &a
}

func (a *app) getFocused(ctx context.Context) (appID string, pid uint32, err error) {
	n, err := a.client.GetTree(ctx)
	if err != nil {
		return
	}

	node := n.FocusedNode()

	if node == nil {
		return
	}

	if node.AppID != nil {
		appID = *node.AppID
	}

	if node.PID != nil {
		pid = *node.PID
	}

	return
}

func (a *app) run(ctx context.Context) error {
	flag.Parse()

	var err error
	a.client, err = sway.New(ctx, sway.WithSocketPath(a.socketPath))
	if err != nil {
		return err
	}

	var k kitty
	if k.appID, k.pid, err = a.getFocused(ctx); err != nil {
		return err
	}

	return k.Run(ctx)
}
