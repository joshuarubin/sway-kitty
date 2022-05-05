package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type kitty struct {
	appID string
	pid   uint32
}

func (k *kitty) Run(ctx context.Context) error {
	cmd, err := k.Cmd(ctx)
	if err != nil {
		return err
	}

	return cmd.Run()
}

func (k *kitty) Cmd(ctx context.Context) (*exec.Cmd, error) {
	cmd := exec.Command("kitty")

	var err error
	if k.appID == "kitty" || k.appID == "kitty_autostart" {
		err = k.WindowCmd(ctx, cmd)
	} else {
		err = k.NewCmd(ctx, cmd)
	}
	if err != nil {
		return nil, err
	}

	if len(os.Args) > 1 {
		cmd.Args = append(cmd.Args, "--")
		cmd.Args = append(cmd.Args, os.Args[1:]...)
	}

	return cmd, nil
}

func (k *kitty) NewCmd(ctx context.Context, cmd *exec.Cmd) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cmd.Args = append(cmd.Args,
		"--working-directory", home,
		"--detach",
	)

	return nil
}

func (k *kitty) WindowCmd(ctx context.Context, cmd *exec.Cmd) error {
	cwd, err := k.FocusedCWD(ctx)
	if err != nil {
		return err
	}

	cmd.Args = append(cmd.Args,
		"@",
		"--to", fmt.Sprintf("unix:@kitty-%d", k.pid),
		"launch",
		"--type", "window",
		"--cwd", cwd,
	)

	return nil
}

func (k *kitty) FocusedCWD(ctx context.Context) (string, error) {
	cmd := exec.Command("kitty",
		"@",
		"--to", fmt.Sprintf("unix:@kitty-%d", k.pid),
		"ls",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	var ws osWindows
	done := make(chan struct{})
	var decodeErr error
	go func() {
		defer close(done)
		decodeErr = json.NewDecoder(stdout).Decode(&ws)
	}()

	if err = cmd.Run(); err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-done:
		if decodeErr != nil {
			return "", decodeErr
		}
	}

	return ws.GetFocusedCWD(), nil
}
