// +build integration

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/onsi/gocleanup"
)

func init() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGUSR1)
		<-c
		gocleanup.Exit(166)
	}()
}
