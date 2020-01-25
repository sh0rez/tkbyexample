package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
)

func develop() error {
	if err := render(); err != nil {
		return err
	}

	// file watcher
	go watch()

	// invoke gatsby
	gatsby, killGatsby := gatsbyDevelop()
	go gatsby()
	defer killGatsby()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		return nil
	case err := <-e:
		return err
	}
}

func watch() {
	w := watcher.New()
	defer w.Close()

	w.SetMaxEvents(1)
	check(w.AddRecursive("./src/examples"))
	check(w.Add("./src/example.tmpl"))

	go func() {
		for {
			select {
			case <-w.Event:
				check(render())
			case err := <-w.Error:
				check(err)
			case <-w.Closed:
				return
			}
		}
	}()

	check(w.Start(time.Millisecond * 100))
}

func check(err error) {
	if err != nil {
		e <- err
	}
}

func gatsbyDevelop() (run func(), kill func()) {
	cmd := exec.Command("yarn", "dev")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	run = func() {
		check(cmd.Run())
	}

	kill = func() {
		check(syscall.Kill(-cmd.Process.Pid, syscall.SIGINT))
		time.Sleep(100 * time.Millisecond)
	}

	return run, kill
}
