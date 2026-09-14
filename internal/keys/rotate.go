package keys

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/johan253/idme/internal/db"
)

// Rotator generates and promotes new signing keys. Runs out-of-band (CLI or
// admin endpoint), not from the request path.
type Rotator struct {
	pool   PgxPool
	cipher Cipher
}

// PgxPool is the subset of *pgxpool.Pool needed to run rotation transactions.
type PgxPool interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}

func NewRotator(pool PgxPool, cipher Cipher) *Rotator {
	return &Rotator{pool: pool, cipher: cipher}
}

// Generate creates a new inactive signing key. Pods pick it up on their next
// reload; publish it as a verification key first, then Promote once every
// pod has seen it.
func (r *Rotator) Generate(ctx context.Context) (string, error) {
	kp, err := GenerateEd25519()
	if err != nil {
		return "", err
	}
	pubPEM, err := encodePublicPEM(kp.Public)
	if err != nil {
		return "", err
	}
	privPEM, err := encodePrivatePEM(kp.Private)
	if err != nil {
		return "", err
	}
	sealed, err := r.cipher.Seal(privPEM)
	if err != nil {
		return "", err
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)
	if _, err := q.InsertSigningKey(ctx, db.InsertSigningKeyParams{
		Kid:        kp.Kid,
		PublicPem:  pubPEM,
		PrivateEnc: sealed,
		IsActive:   false,
	}); err != nil {
		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return kp.Kid, nil
}

// Promote makes kid the sole active signing key. The partial unique index
// on signing_keys(is_active) makes this safe under concurrent rotations —
// one caller wins, the other errors.
func (r *Rotator) Promote(ctx context.Context, kid string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := db.New(tx)
	if err := q.DeactivateAllSigningKeys(ctx); err != nil {
		return err
	}
	if err := q.ActivateSigningKey(ctx, kid); err != nil {
		return fmt.Errorf("activate %s: %w", kid, err)
	}
	return tx.Commit(ctx)
}
