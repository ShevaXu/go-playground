package grace

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// EnvMark is a temporary environment variable
// to mark a server restart.
const EnvMark = "GRACE_RESTART"

// GraceServer wraps a http.Server to
// utilize the GraceListener and handle signal & restart.
type GraceServer struct {
	Addr    string
	Handler http.Handler
	srv     *http.Server
	l       *GraceListener
}

// Restart forks a new process with current listener's os.File
// and set environment variable to mark a restart.
func (s *GraceServer) Restart() error {
	os.Setenv(EnvMark, "1")
	f, err := s.l.File()
	if err != nil {
		return err
	}
	err = Fork([]*os.File{f})
	if err != nil {
		return err
	}
	return nil
}

// Shutdown breaks the serve() loop by closing listener.
func (s *GraceServer) Shutdown() {
	err := s.l.Close()
	if err != nil {
		log.Printf("Error shutdown %d - %v\n", os.Getpid(), err)
	} else {
		log.Printf("Server %d shutdown\n", os.Getpid())
	}
}

// handleSignal intercepts system signals to
// either shutdown or restart server.
// The restart signal is only handled once.
func (s *GraceServer) handleSignal() {
	var sig os.Signal
	var err error

	ch := make(chan os.Signal)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	for {
		sig = <-ch
		log.Println(os.Getpid(), " receive signal", sig)
		switch sig {
		case syscall.SIGTERM:
			s.Shutdown()
		case syscall.SIGINT:
			err = s.Restart()
			if err != nil {
				log.Printf("Error restart %d - %v\n", os.Getpid(), err)
			} else {
				log.Printf("Server %d forked\n", os.Getpid())
				// Ignore() still intercept the signal,
				// Reset() will let the signal kill the process.
				signal.Ignore(sig)
				s.Shutdown()
			}
		default:
		}
	}
}

// Serve makes the underlying server listen to connections,
// dispatches a goroutine to handle signal and
// blocks after serve() returns to make sure
// all ongoing connections handled or closed.
func (s *GraceServer) Serve(l *GraceListener) error {
	s.l = l
	go s.handleSignal()
	s.srv = &http.Server{Addr: s.Addr, Handler: s.Handler}
	err := s.srv.Serve(l)
	l.Wait()
	return err
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

	return s.Serve(NewGraceListener(ln.(*net.TCPListener)))
}

// ListenAndServe provides shorthand to init and serve a GraceServer,
// as a drop in replacement of http.ListenAndServe.
func ListenAndServe(addr string, handler http.Handler) error {
	srv := GraceServer{Addr: addr, Handler: handler}
	return srv.ListenAndServe()
}
