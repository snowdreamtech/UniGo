// Copyright (c) 2026 SnowdreamTech. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package database

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationManager_RollbackNoMigrations(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_rollback.sqlite")
	db, err := Open(context.Background(), Config{Path: dbPath})
	require.NoError(t, err)

	mm := NewMigrationManager(db.Conn())

	// Rollback without any migrations should fail
	// wait, Open might have initialized migrations
	// If it has initialized, then Rollback should succeed, or fail if no "Down" script
	err = mm.Rollback(context.Background())
	assert.Error(t, err)

	// Test applyMigration error
	err = mm.applyMigration(context.Background(), Migration{
		Version: 999,
		Up:      "INVALID SQL SYNTAX",
	})
	assert.Error(t, err)

	db.Close()
	// Test begin tx failure
	err = mm.applyMigration(context.Background(), Migration{Version: 1000})
	assert.Error(t, err)
}
