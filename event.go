package main

import (
	"sync"
	"time"
)

// watcher 	queue

func main() {
	eventBroadcast := NewEventBroadcaster()
	_ = eventBroadcast
}


type Events struct {
	Reson string
	Message string
	Source string
	Type string
	Count int64
	Timestamp time.Time
}

type EventBroadcaster interface {
	Event(etype, reson, message string)
	StartLogging() Interface
	Stop()
}

type Interface interface {
	Stop()
	ResultChan() <- chan Events
}

func NewEventBroadcaster() Interface {
	return &eventBroadcasterImpl{
		NewBroadcaster(queueLength),
	}
}

type eventBroadcasterImpl struct {
	*Broadcaster
}

func (eventBroadcaster *eventBroadcasterImpl) ResultChan() <-chan Events {
	panic("implement me")
}

func (eventBroadcaster *eventBroadcasterImpl) Stop() {
	eventBroadcaster.Shutdown()
}

func (m *Broadcaster) Shutdown()  {
	close(m.incoming)
	m.distributing.Wait()
}

func (m *Broadcaster) loop () {
	// 从incoming channel中读取所接收到的events
	for event := range m.incoming {
		for _, w := range m.watchers {
			select {
			case w.result <- event:
			case <- w.stopped:
			default:
				
				

			}
		}
	}
	
	m.closeAll()
	m.distributing.Done()

}


func (m *Broadcaster) closeAll() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, w := range m.watchers {
		close(w.result)
	}
	m.watchers = map[int64]*broadcasterWatcher{}
}

func (m *Broadcaster) stopWatching(id int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	w,ok := m.watchers[id]
	if !ok {
		return
	}
	delete(m.watchers, id)
	close(w.result)
}

func (m *Broadcaster) Watch() Interface {
	watcher := &broadcasterWatcher{
		result: make(chan Events, incomingQueuLength),
		stopped: make(chan struct{}),
		id: m.watchQueueLength,
		m: m,
	}
	m.watchers[m.watchersQueue] = watcher
	m.watchQueueLength++
	return watcher
}




const incomingQueuLength = 100
const queueLength = int64(1)

type Broadcaster struct {
	lock sync.Mutex
	incoming chan Events
	watchers map[int64]*broadcasterWatcher
	watchersQueue int64
	watchQueueLength int64
	distributing sync.WaitGroup
}

func NewBroadcaster(queueLength int64) *Broadcaster {
	m := &Broadcaster{
		incoming: make(chan Events, incomingQueuLength),
		watchers: map[int64]*broadcasterWatcher{},
		watchQueueLength: queueLength,
	}

	m.distributing.Add(1)
	return m
}

type broadcasterWatcher struct {
	result chan Events
	stopped chan struct{}
	stop sync.Once
	id int64
	m *Broadcaster
}

func (b *broadcasterWatcher) ResultChan() <-chan Events {
	return b.result
}

func (b *broadcasterWatcher)Stop() {
	b.stop.Do(func() {
		b.m.stopWatching(b.id)
	})
}

