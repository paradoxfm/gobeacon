package watch

import (
	"bufio"
	"gobeacon/service"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Addr         string
	IdleTimeout  time.Duration
	MaxReadBytes int64

	listener   net.Listener
	conns      map[*conn]struct{}
	mu         sync.Mutex
	inShutdown bool
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("starting server on %v\n", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	srv.listener = listener
	for {
		// should be guarded by mu
		if srv.inShutdown {
			break
		}
		newConn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		log.Printf("%s accepted connection from %v", time.Now(), newConn.RemoteAddr())
		conn := &conn{
			Conn:          newConn,
			IdleTimeout:   srv.IdleTimeout,
			MaxReadBuffer: srv.MaxReadBytes,
		}
		srv.trackConn(conn)
		conn.SetDeadline(time.Now().Add(conn.IdleTimeout))
		go srv.handle(conn)
	}
	return nil
}

func (srv *Server) trackConn(c *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	if srv.conns == nil {
		srv.conns = make(map[*conn]struct{})
	}
	srv.conns[c] = struct{}{}
}

func (srv *Server) handle(conn *conn) error {
	defer func() {
		log.Printf("%s closing connection from %v", time.Now(), conn.RemoteAddr())
		conn.Close()
		srv.deleteConn(conn)
	}()
	r := bufio.NewReader(conn)
	scanr := bufio.NewScanner(r)

	sc := make(chan bool)
	deadline := time.After(conn.IdleTimeout)
	for {
		go func(s chan bool) {
			s <- scanr.Scan()
		}(sc)
		select {
		case <-deadline:
			return nil
		case scanned := <-sc:
			if !scanned {
				if err := scanr.Err(); err != nil {
					return err
				}
				return nil
			}
			log.Printf("accept string %s", scanr.Text())
			//respMsg := service.WatchHandleMessage(scanr.Bytes())
			respMsg := service.WatchHandleMessage2(scanr.Text())
			if respMsg != nil {
				w := bufio.NewWriter(conn)
				log.Printf("write string %s", respMsg)
				w.WriteString(respMsg.(string))
				//w.WriteString(respMsg + "\n")
				w.Flush()
			}
			deadline = time.After(conn.IdleTimeout)
		}
	}
	return nil
}

func (srv *Server) deleteConn(conn *conn) {
	defer srv.mu.Unlock()
	srv.mu.Lock()
	delete(srv.conns, conn)
}

func (srv *Server) Shutdown() {
	// should be guarded by mu
	srv.inShutdown = true
	log.Println("shutting down...")
	srv.listener.Close()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("waiting on %v connections", len(srv.conns))
		}
		if len(srv.conns) == 0 {
			return
		}
	}
}
