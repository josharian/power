package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/josharian/power"
)

// TODO: use process groups to suspend more processes?

const usage = `usage: nobatt cmd [args]

nobatt runs cmd.

It suspends cmd when it detects that the laptop is running on battery power,
and resumes it when it detects that the laptop is using wall power.
If you want to override nobatt, send SIGUSR1 to force cmd to resume,
or SIGUSR2 to force cmd to suspend.

Note that nobatt only controls cmd, not any other processes that cmd starts.
`

func main() {
	log.SetFlags(0)
	if len(os.Args) == 1 {
		log.Print(usage)
		os.Exit(1)
	}
	c, err := power.Notify(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Fatal(err) // TODO: remove?
		}
		close(done)
	}()

	sigc := make(chan os.Signal)
	signal.Notify(sigc, syscall.SIGUSR1, syscall.SIGUSR2)

	running := true
	run := func(want bool) {
		// log.Printf("running=%v want=%v", running, want)
		if want == running {
			// already in correct state
			return
		}
		var sig os.Signal
		if want {
			sig = syscall.SIGCONT
		} else {
			sig = syscall.SIGSTOP
		}
		err := cmd.Process.Signal(sig)
		// log.Printf("sent signal %v to cmd, err=%v", sig, err)
		if err != nil {
			log.Fatal(err)
		}
		running = want
	}

loop:
	for {
		select {
		case src := <-c:
			run(src == power.AC)
		case sig := <-sigc:
			run(sig == syscall.SIGUSR1)
		case <-done:
			break loop
		}
	}
}
