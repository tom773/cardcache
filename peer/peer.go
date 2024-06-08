package peer

import (
	"fmt"
	"net"
	"sync"

	"github.com/tom773/cardcache/pkg"
)

type Peer struct {
	Conn  net.Conn
	MsgCh chan Message
	OutCh chan []byte
}

type Message struct {
	Peer    *Peer
	Content []byte
}

func NewPeer(Conn net.Conn, MsgCh chan Message) *Peer {
	return &Peer{Conn: Conn, MsgCh: MsgCh, OutCh: make(chan []byte)}
}

func (p *Peer) ReadLoop() error {
	for {
		buf := make([]byte, 1024)
		n, err := p.Conn.Read(buf)
		if err != nil {
			return err
		}

		msg := Message{
			Peer:    p,
			Content: make([]byte, n),
		}

		copy(msg.Content, buf[:n])
		p.MsgCh <- msg
	}
}

func (p *Peer) WriteLoop() error {
	for {
		select {
		case msg := <-p.OutCh:
			finalmsg := append(msg, []byte("\n")...)
			_, err := p.Conn.Write(finalmsg)
			if err != nil {
				return err
			}
		}
	}
}

// Implementing Pub/Sub here
type Subscriber struct {
	key string
	p   *Peer
}

type PubSub struct {
	subscribers map[string][]Subscriber
	mu          sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]Subscriber),
	}
}

func (ps *PubSub) Subscribe(peer *Peer, key string) chan []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := Subscriber{
		key: key,
		p:   peer,
	}
	ps.subscribers[key] = append(ps.subscribers[key], ch)
	return ch.p.OutCh
}

func (ps *PubSub) Publish(key, value string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	msg := fmt.Sprintf("Key %s%s%s, the Value is now: %s%s%s", pkg.Colors.Cyan, key, pkg.Colors.Reset, pkg.Colors.Green, value, pkg.Colors.Reset)
	for _, ch := range ps.subscribers[key] {
		go func(ch chan []byte) {
			ch <- []byte(msg)
		}(ch.p.OutCh)
	}
}
