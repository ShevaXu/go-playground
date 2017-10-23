package grace

import (
	"net"
	"os"
	"sync"
	"time"
)

// GraceConn synchronizes closing connections through WaitGroup.
type GraceConn struct {
	net.Conn // embedded
	wg       *sync.WaitGroup
}

// Close makes sure the WaitGroup counter decrements
// whether closing fails or not.
func (c GraceConn) Close() error {
	defer c.wg.Done()
	return c.Conn.Close()
}

// GraceListener accepts connections and wrap them in GraceConn
// to allow WaitGroup synchronization.
type GraceListener struct {
	tl *net.TCPListener
	wg *sync.WaitGroup
}

// Accept implements the Accept method in the Listener interface;
// it increases the WaitGroup counter.
func (l *GraceListener) Accept() (net.Conn, error) {
	tc, err := l.tl.AcceptTCP()
	if err != nil {
		return tc, err
	}

	tc.SetKeepAlive(true)
	// TODO: allow setting period?
	tc.SetKeepAlivePeriod(time.Minute)

	l.wg.Add(1)

	c := GraceConn{
		Conn: tc,
		wg:   l.wg,
	}

	return c, nil
}

// File returns a copy of the underlying TCPListener os.File.
func (l *GraceListener) File() (*os.File, error) {
	return l.tl.File()
}

// Wait blocks the goroutine using the WaitGroup
// since it is not exposed.
func (l *GraceListener) Wait() {
	l.wg.Wait()
}

// Close implements the Close method in the Listener interface.
func (l *GraceListener) Close() error {
	return l.tl.Close()
}

// Addr implements the Addr method in the Listener interface.
func (l *GraceListener) Addr() net.Addr {
	return l.tl.Addr()
}

// NewGraceListener exposes a method to
// make a GraceListener providing only a TCPListener;
// WaitGroup is initialized here.
func NewGraceListener(tl *net.TCPListener) *GraceListener {
	return &GraceListener{
		tl: tl,
		wg: &sync.WaitGroup{},
	}
}
