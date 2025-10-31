package handlers

import (
	"sync"
	"time"
)

// dedupeKey is a composite identifier for idempotency
type dedupeKey struct {
	Source     string
	DeliveryID string
	PayloadSum string
}

type entry struct {
	expiresAt time.Time
}

// Deduper is an in-memory TTL store for idempotent request detection.
// NOTE: This is a temporary fallback before DB-backed event_deliveries is wired.
type Deduper struct {
	mu    sync.Mutex
	ttl   time.Duration
	items map[dedupeKey]entry
}

func NewDeduper(ttl time.Duration) *Deduper {
	d := &Deduper{
		ttl:   ttl,
		items: make(map[dedupeKey]entry, 1024),
	}
	go d.gc()
	return d
}

// CheckAndMark returns true if the key is seen within TTL, else marks it and returns false.
func (d *Deduper) CheckAndMark(source, deliveryID, payloadSum string) bool {
	if d == nil {
		return false
	}
	k := dedupeKey{Source: source, DeliveryID: deliveryID, PayloadSum: payloadSum}
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()
	if e, ok := d.items[k]; ok && e.expiresAt.After(now) {
		return true
	}
	d.items[k] = entry{expiresAt: now.Add(d.ttl)}
	return false
}

func (d *Deduper) gc() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		d.mu.Lock()
		for k, e := range d.items {
			if e.expiresAt.Before(now) {
				delete(d.items, k)
			}
		}
		d.mu.Unlock()
	}
}
