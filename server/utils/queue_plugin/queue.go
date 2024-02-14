package Queue

import (
	"errors"
	"sync"
	"time"
)

type Queue struct {
	exit     chan bool
	capacity int
	topics   map[string][]chan any
	sync.RWMutex
	//once         sync.Once
}

func NewQueue() *Queue {
	return &Queue{
		exit:   make(chan bool),
		topics: make(map[string][]chan any),
	}
}

func (q *Queue) ShowExit() chan bool {
	return q.exit

}

func (q *Queue) SetConditions(capacity int) {
	q.capacity = capacity
}

func (q *Queue) Start() {
	select {
	case <-q.exit:
		q.exit = make(chan bool)
	default:
		return
	}
}
func (q *Queue) Close() {
	select {
	case <-q.exit:
		return
	default:
		close(q.exit)
		q.Lock()
		q.topics = make(map[string][]chan any)
		q.Unlock()
	}
	return
}

func (q *Queue) Publish(topic string, pub any) error {
	select {
	case <-q.exit:
		return errors.New("Queue is closed")
	default:
	}
	q.RLock()
	subscribers, ok := q.topics[topic]
	q.RUnlock()
	if !ok {
		return nil
	}
	q.Broadcast(pub, subscribers)
	return nil
}

func (q *Queue) Broadcast(msg any, subscribers []chan any) {
	count := len(subscribers)
	concurrency := 1
	switch {
	case count > 1000:
		concurrency = 3
	case count > 100:
		concurrency = 2
	default:
		concurrency = 1
	}
	pub := func(start int) {
		idleDuration := 5 * time.Millisecond
		ticker := time.NewTicker(idleDuration)
		defer ticker.Stop()
		for j := start; j < count; j += concurrency {
			select {
			case subscribers[j] <- msg:
			case <-ticker.C:
			case <-q.exit:
				return
			}
		}
	}
	for i := 0; i < concurrency; i++ {
		go pub(i)
	}
}

func (q *Queue) Subscribe(topic string, capacity ...int) (<-chan any, error) {
	select {
	case <-q.exit:
		return nil, errors.New("Queue is closed")
	default:
	}
	if q.capacity == 0 {
		q.capacity = 100
	}
	var c = q.capacity
	if len(capacity) > 0 {
		if capacity[0] != 0 {
			c = capacity[0]
		}
	}
	ch := make(chan any, c)
	q.Lock()
	q.topics[topic] = append(q.topics[topic], ch)
	q.Unlock()
	return ch, nil
}

func (q *Queue) Unsubscribe(topic string, sub <-chan any) error {
	select {
	case <-q.exit:
		return errors.New("Queue is closed")
	default:
	}
	q.RLock()
	subscribers, ok := q.topics[topic]
	q.RUnlock()
	if !ok {
		return nil
	}
	q.Lock()
	var newSubs []chan any
	for _, subscriber := range subscribers {
		if subscriber == sub {
			continue
		}
		newSubs = append(newSubs, subscriber)
	}
	q.topics[topic] = newSubs
	q.Unlock()
	return nil
}

func (q *Queue) GetPayLoad(sub <-chan any) any {
	for val := range sub {
		if val != nil {
			return val
		}
	}
	return nil
}
