// Package power detects and reports the power source of a laptop.
package power

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// TODO: support linux

// TODO: make a cgo version of this using https://developer.apple.com/documentation/iokit/iopowersources_h

// TODO: allow reading and notifying on battery percentage

// TODO: test on other macOS versions (tested only on Mojave)

// Current returns the current power source.
func Current() (Source, error) {
	cmd := exec.Command("pmset", "-g", "ps")
	out, err := cmd.Output()
	if err != nil {
		return Unknown, err
	}
	i := bytes.IndexByte(out, '\n')
	if i == -1 {
		return Unknown, fmt.Errorf("could not parse pmset output: %q", out)
	}
	out = out[:i]
	src, ok := parse(out)
	if !ok {
		return Unknown, fmt.Errorf("could not parse pmset line: %q", out)
	}
	return src, nil
}

// Wait blocks until ctx is done or s is the current power source.
func Wait(ctx context.Context, s Source) error {
	// Create our own context so that we don't leak resources when Wait completes without error.
	nctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c, err := Notify(nctx)
	if err != nil {
		return err
	}
	for {
		select {
		case cur := <-c:
			if cur == s {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// C provides updates when the power source changes,
// including an initial update to provide the current power source.
// It stops providing updates and releases resources when ctx is done.
func Notify(ctx context.Context) (<-chan Source, error) {
	cmd := exec.CommandContext(ctx, "pmset", "-g", "pslog")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	c := make(chan Source)
	go func() {
		s := bufio.NewScanner(out)
	scanloop:
		for s.Scan() {
			src, ok := parse(s.Bytes())
			if !ok {
				continue
			}
			select {
			case <-ctx.Done():
				break scanloop
			case c <- src:
			}
		}
		// Ignore s.Error(); there's nothing we can do with it.
	}()
	return c, nil
}

// parse parses a single line from pmset output.
func parse(b []byte) (Source, bool) {
	const prefix = "Now drawing from '"
	if !bytes.HasPrefix(b, []byte(prefix)) {
		return Unknown, false
	}
	b = b[len(prefix):]
	if len(b) == 0 || b[len(b)-1] != '\'' {
		return Unknown, false
	}
	b = b[:len(b)-1]

	switch string(b) {
	case "AC Power":
		return AC, true
	case "Battery Power":
		return Battery, true
	case "UPS Power":
		return UPS, true
	}
	return Unknown, true
}
