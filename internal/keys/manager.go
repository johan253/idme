package keys

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/johan253/idme/internal/db"
)

// Store is the subset of *db.Queries the manager needs. Kept as an interface
// so tests can fake it.
type Store interface {
	ListSigningKeys(ctx context.Context) ([]db.SigningKey, error)
}

// Manager holds the current signing key and all publishable verification
// keys in memory, refreshing periodically from the database.
type Manager struct {
	store   Store
	cipher  Cipher
	refresh time.Duration

	mu     sync.RWMutex
	active *KeyPair
	byKid  map[string]*KeyPair
}

func NewManager(store Store, cipher Cipher, refresh time.Duration) *Manager {
	return &Manager{
		store:   store,
		cipher:  cipher,
		refresh: refresh,
		byKid:   map[string]*KeyPair{},
	}
}

// Start does an initial load and then refreshes on an interval until ctx is
// cancelled. Returns an error if the first load fails so callers can fail
// fast at startup.
func (m *Manager) Start(ctx context.Context) error {
	if err := m.Reload(ctx); err != nil {
		return err
	}
	go m.loop(ctx)
	return nil
}

func (m *Manager) loop(ctx context.Context) {
	t := time.NewTicker(m.refresh)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			_ = m.Reload(ctx)
		}
	}
}

// Reload rebuilds the in-memory view from the database.
func (m *Manager) Reload(ctx context.Context) error {
	start := time.Now()
	rows, err := m.store.ListSigningKeys(ctx)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	log.Printf("manager: found %v keys duration=%v", len(rows), duration)

	var active *KeyPair
	byKid := make(map[string]*KeyPair, len(rows))

	for _, r := range rows {
		pub, err := decodePublicPEM(r.PublicPem)
		if err != nil {
			return err
		}
		kp := &KeyPair{Kid: r.Kid, Public: pub}
		if r.IsActive {
			privPEM, err := m.cipher.Open(r.PrivateEnc)
			if err != nil {
				return err
			}
			priv, err := decodePrivatePEM(privPEM)
			if err != nil {
				return err
			}
			kp.Private = priv
			active = kp
		}
		byKid[r.Kid] = kp
	}

	m.mu.Lock()
	m.active = active
	m.byKid = byKid
	m.mu.Unlock()
	return nil
}

var ErrNoActiveKey = errors.New("no active signing key")

// Active returns the current signing keypair.
func (m *Manager) Active() (*KeyPair, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.active == nil {
		return nil, ErrNoActiveKey
	}
	return m.active, nil
}

// ByKid returns a verification key by kid, or nil if unknown.
func (m *Manager) ByKid(kid string) *KeyPair {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byKid[kid]
}

// All returns every currently-published key (for JWKS).
func (m *Manager) All() []*KeyPair {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*KeyPair, 0, len(m.byKid))
	for _, k := range m.byKid {
		out = append(out, k)
	}
	return out
}
