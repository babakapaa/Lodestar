package live

import "sync"

const subscriberBuffer = 64

type subscriber struct {
	events chan Event
	lagged chan struct{}
}

type Hub struct {
	mu     sync.Mutex
	subs   map[string]map[*subscriber]struct{}
	done   chan struct{}
	closed bool
}

func NewHub() *Hub {
	return &Hub{subs: map[string]map[*subscriber]struct{}{}, done: make(chan struct{})}
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.closed {
		h.closed = true
		close(h.done)
	}
}

func (h *Hub) subscribe(tenantID, userID string) *subscriber {
	s := &subscriber{events: make(chan Event, subscriberBuffer), lagged: make(chan struct{})}
	key := tenantID + ":" + userID
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.subs[key] == nil {
		h.subs[key] = map[*subscriber]struct{}{}
	}
	h.subs[key][s] = struct{}{}
	return s
}

func (h *Hub) unsubscribe(tenantID, userID string, s *subscriber) {
	key := tenantID + ":" + userID
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subs[key], s)
	if len(h.subs[key]) == 0 {
		delete(h.subs, key)
	}
}

func (h *Hub) watching(tenantID, userID string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subs[tenantID+":"+userID]) > 0
}

func (h *Hub) dispatch(tenantID, userID string, e Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for s := range h.subs[tenantID+":"+userID] {
		select {
		case s.events <- e:
		default:
			select {
			case <-s.lagged:
			default:
				close(s.lagged)
			}
		}
	}
}
