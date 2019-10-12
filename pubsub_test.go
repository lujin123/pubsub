package pubsub

import (
	"fmt"
	"sync"
	"testing"
)

const (
	topicHello = "hello"
)

func TestPubSub(t *testing.T) {
	var wg sync.WaitGroup
	n := 4
	wg.Add(n + 1)

	pb := New()
	go publish(topicHello, pb)
	for i := 1; i <= n; i++ {
		go sub(i, pb, &wg)
	}
	go subOnce(pb, &wg)
	wg.Wait()
}

func publish(topic string, pb *PubSub) {
	for i := 1; i < 10; i++ {
		pb.Publish(fmt.Sprintf("hello world %d", i), topic)
	}
}

func sub(id int, pb *PubSub, wg *sync.WaitGroup) {
	defer wg.Done()
	ch := pb.Sub(topicHello)
	for i := 1; ; i++ {
		if i == 5 {
			go pb.Unsub(ch, topicHello)
		}
		msg, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("id=%d, msg=%+v, times=%d\n", id, msg, i)
	}
}

func subOnce(pb *PubSub, wg *sync.WaitGroup) {
	defer wg.Done()
	ch := pb.SubOnce(topicHello)
	for msg := range ch {
		fmt.Printf("id=subonece, msg=%+v\n", msg)
	}
}
