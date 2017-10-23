package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// EnvMark is a temporary environment variable
// to mark a server restart.
const EnvMark = "GRACE_RESTART"

// Fork starts a new process executing the same command
// as current process with same arguments,
// stdout, stderr and extraFiles are inherited by the new process.
func Fork(extraFiles []*os.File) error {
	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if len(extraFiles) > 0 {
		cmd.ExtraFiles = extraFiles
	}

	return cmd.Start()
}

// GraceServer wraps a http.Server to
// utilize the GraceListener and handle signal & restart.
type GraceServer struct {
	Addr    string
	Handler http.Handler
	srv     *http.Server
}

// Restart forks a new process with a provided *os.File
// and set environment variable to mark a restart.
func (s *GraceServer) Restart(f *os.File) error {
	os.Setenv(EnvMark, "1")
	err := Fork([]*os.File{f})
	if err != nil {
		return err
	}
	return nil
}

// HammerTime is the maximal duration given the server to shutdown.
const HammerTime = 30 * time.Second

// Shutdown calls Shutdown() of the underlying http.Server
// with a timeout context and log.
func (s *GraceServer) Shutdown() {
	ctx, _ := context.WithTimeout(context.Background(), HammerTime)
	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutdown %d - %v\n", os.Getpid(), err)
	} else {
		log.Printf("Server %d shutdown\n", os.Getpid())
	}
}

// Serve makes the underlying server listen to connections
// in a separate goroutine and blocks to wait for
// system signals to either Restart or Shutdown.
func (s *GraceServer) Serve(l net.Listener) error {
	s.srv = &http.Server{Addr: s.Addr, Handler: s.Handler}

	if ln, ok := l.(*net.TCPListener); ok {
		var sig os.Signal

		f, err := ln.File()
		if err != nil {
			return err
		}

		ch := make(chan os.Signal)
		signal.Notify(
			ch,
			syscall.SIGTERM,
			syscall.SIGINT,
		)

		go func() {
			if err := s.srv.Serve(l); err != nil {
				log.Fatalln("Fatal error from Serve():", err)
			}
		}()

		sig = <-ch
		log.Println(os.Getpid(), " receive signal", sig)
		switch sig {
		case syscall.SIGTERM:
			s.Shutdown()
		case syscall.SIGINT:
			err = s.Restart(f)
			if err != nil {
				log.Printf("Error restart %d - %v\n", os.Getpid(), err)
			} else {
				log.Printf("Server %d forked\n", os.Getpid())
				// Ignore() still intercept the signal,
				// Reset() will let the signal kill the process.
				signal.Ignore(sig)
				s.Shutdown()
			}
		}

		return nil
	} else {
		return s.srv.Serve(l)
	}
}

// ListenAndServe sets up the listener properly for restart.
func (s *GraceServer) ListenAndServe() error {
	var ln net.Listener
	var err error

	if s.Addr == "" {
		s.Addr = ":http"
	}

	if os.Getenv(EnvMark) != "" {
		f := os.NewFile(3, "")
		ln, err = net.FileListener(f)
	} else {
		// fresh start
		ln, err = net.Listen("tcp", s.Addr)
	}
	if err != nil {
		return err
	}

	return s.Serve(ln)
}

// ListenAndServe provides shorthand to init and serve a GraceServer,
// as a drop in replacement of http.ListenAndServe.
func ListenAndServe(addr string, handler http.Handler) error {
	srv := GraceServer{Addr: addr, Handler: handler}
	return srv.ListenAndServe()
}
