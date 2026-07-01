// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package database

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// simulateItemInsert simulates an upsert pattern on example_items
func simulateItemInsert(ctx context.Context, db *DB, key, value string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO example_items (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = CURRENT_TIMESTAMP
	`, key, value)
	if err != nil {
		return fmt.Errorf("upsert example_items: %w", err)
	}

	return tx.Commit()
}

func TestConcurrentWrites_8Workers(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, Config{Path: filepath.Join(t.TempDir(), "concurrent_8.db")})
	require.NoError(t, err)
	defer db.Close()

	keys := []string{"k1", "k2", "k3", "k4", "k5", "k6", "k7", "k8"}

	var wg sync.WaitGroup
	var errCount atomic.Int64

	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			if err := simulateItemInsert(ctx, db, key, "value"); err != nil {
				t.Logf("ERROR inserting %s: %v", key, err)
				errCount.Add(1)
			}
		}(key)
	}

	wg.Wait()
	assert.Equal(t, int64(0), errCount.Load(), "expected 0 errors from 8 concurrent inserts")

	var count int
	err = db.Conn().QueryRowContext(ctx, `SELECT COUNT(*) FROM example_items`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 8, count, "all 8 items should be persisted")
}

func TestConcurrentWrites_16Workers(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, Config{Path: filepath.Join(t.TempDir(), "concurrent_16.db")})
	require.NoError(t, err)
	defer db.Close()

	const workers = 16
	var wg sync.WaitGroup
	var errCount atomic.Int64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := simulateItemInsert(ctx, db, fmt.Sprintf("k%d", i), "value"); err != nil {
				errCount.Add(1)
			}
		}(i)
	}

	wg.Wait()
	assert.Equal(t, int64(0), errCount.Load(), "expected 0 errors")

	var count int
	err = db.Conn().QueryRowContext(ctx, `SELECT COUNT(*) FROM example_items`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, workers, count)
}

func TestConcurrentUpsert_Idempotency(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, Config{Path: filepath.Join(t.TempDir(), "upsert_idempotent.db")})
	require.NoError(t, err)
	defer db.Close()

	const workers = 16
	var wg sync.WaitGroup
	var errCount atomic.Int64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := simulateItemInsert(ctx, db, "same-key", "value"); err != nil {
				errCount.Add(1)
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(0), errCount.Load(), "concurrent upserts of same key should never error")

	var count int
	err = db.Conn().QueryRowContext(ctx, `SELECT COUNT(*) FROM example_items WHERE key='same-key'`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "idempotent upsert should yield exactly 1 row")
}

func TestConcurrentMixedReadWrite(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, Config{Path: filepath.Join(t.TempDir(), "mixed_rw.db")})
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, simulateItemInsert(ctx, db, "seed-key", "value"))

	const workers = 8
	var wg sync.WaitGroup
	var writeErr, readErr atomic.Int64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := simulateItemInsert(ctx, db, fmt.Sprintf("write-key-%d", i), "value"); err != nil {
				writeErr.Add(1)
			}
		}(i)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var count int
			if err := db.Conn().QueryRowContext(ctx, `SELECT COUNT(*) FROM example_items`).Scan(&count); err != nil {
				readErr.Add(1)
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(0), writeErr.Load())
	assert.Equal(t, int64(0), readErr.Load())
}

func TestBusyTimeout_ExceedGracefully(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, Config{Path: filepath.Join(t.TempDir(), "busy_timeout.db")})
	require.NoError(t, err)
	defer db.Close()

	holdDone := make(chan struct{})
	go func() {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			close(holdDone)
			return
		}
		time.Sleep(200 * time.Millisecond)
		_ = tx.Rollback()
		close(holdDone)
	}()

	time.Sleep(10 * time.Millisecond)

	err = simulateItemInsert(ctx, db, "queued-key", "value")
	assert.NoError(t, err, "write should succeed after waiting for lock release")

	<-holdDone
}
