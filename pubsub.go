package pubsub

import "log"

const (
	subTypeNormal byte = iota
	subTypeOnce

	cmdTypeClose = "pb_close"
	cmdTypeSub   = "pb_sub"
	cmdTypePub   = "pb_pub"
	cmdTypeUnsub = "pb_unsub"
)

type command struct {
	typ     string
	topics  []string
	message interface{}
	client  *client
}

type client struct {
	typ byte
	ch  chan interface{}
}

func newClient(typ byte) *client {
	return &client{
		typ: typ,
		ch:  make(chan interface{}),
	}
}

type topic struct {
	name    string
	message interface{}
	subs    map[chan interface{}]*client
}

type PubSub struct {
	topics  map[string]*topic
	inputCh chan *command
}

func New() *PubSub {
	pb := &PubSub{
		topics:  make(map[string]*topic),
		inputCh: make(chan *command),
	}
	go pb.run()
	return pb
}

func (pb *PubSub) Publish(message interface{}, name ...string) {
	cmd := command{
		typ:     cmdTypePub,
		topics:  name,
		message: message,
	}

	pb.inputCh <- &cmd
}

func (pb *PubSub) Sub(name ...string) chan interface{} {
	cmd := command{
		typ:    cmdTypeSub,
		topics: name,
		client: newClient(subTypeNormal),
	}
	pb.inputCh <- &cmd
	return cmd.client.ch
}

func (pb *PubSub) SubOnce(name ...string) chan interface{} {
	cmd := command{
		typ:    cmdTypeSub,
		topics: name,
		client: newClient(subTypeOnce),
	}
	pb.inputCh <- &cmd
	return cmd.client.ch
}

func (pb *PubSub) Unsub(ch chan interface{}, name ...string) {
	cmd := command{
		typ:    cmdTypeUnsub,
		topics: name,
		client: &client{
			typ: 0,
			ch:  ch,
		},
	}
	pb.inputCh <- &cmd
}

func (pb *PubSub) Close() {
	cmd := command{
		typ: cmdTypeClose,
	}
	pb.inputCh <- &cmd
}

func (pb *PubSub) run() {
	for cmd := range pb.inputCh {
		log.Printf("pubsub recv: %+v\n", cmd)
		switch cmd.typ {
		case cmdTypeClose:
			close(pb.inputCh)
			return
		case cmdTypePub:
			pb.pub(cmd)
		case cmdTypeSub:
			pb.sub(cmd)
		case cmdTypeUnsub:
			pb.removeSub(cmd)
		}
	}
}

func (pb *PubSub) pub(cmd *command) {
	for _, name := range cmd.topics {
		tp := pb.topics[name]
		if tp != nil {
			for i := range tp.subs {
				tp.subs[i].ch <- cmd.message
				if tp.subs[i].typ == subTypeOnce {
					close(tp.subs[i].ch)
					delete(tp.subs, tp.subs[i].ch)
				}
			}
		}
	}
}

func (pb *PubSub) sub(cmd *command) {
	for _, name := range cmd.topics {
		tp, ok := pb.topics[name]
		if !ok {
			tp = &topic{
				name: name,
				subs: make(map[chan interface{}]*client),
			}
			pb.topics[name] = tp
		}
		tp.subs[cmd.client.ch] = cmd.client
	}
}

func (pb *PubSub) removeSub(cmd *command) {
	for _, name := range cmd.topics {
		tp := pb.topics[name]
		if tp != nil {
			close(cmd.client.ch)
			delete(tp.subs, cmd.client.ch)
		}
	}
}
