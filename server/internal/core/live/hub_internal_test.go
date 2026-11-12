package live

import "testing"

func TestHubMarksSlowSubscriberAsLagged(t *testing.T) {
	h := NewHub()
	s := h.subscribe("t", "u")
	for i := range subscriberBuffer + 2 {
		h.dispatch("t", "u", Event{ID: int64(i)})
	}
	select {
	case <-s.lagged:
	default:
		t.Fatal("overflowing subscriber is not marked lagged")
	}
	h.unsubscribe("t", "u", s)
	if h.watching("t", "u") {
		t.Fatal("hub still watches after unsubscribe")
	}
}

func TestSignalRoundTrip(t *testing.T) {
	s := signal{tenantID: "a", userID: "b", eventID: 42}
	got, err := decodeSignal(s.encode())
	if err != nil || got != s {
		t.Fatalf("got %+v, %v", got, err)
	}
}
