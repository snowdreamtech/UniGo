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

// TestDatabaseIntegration tests the complete database workflow
func TestDatabaseIntegration(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "integration.db")

	// Requirement 2.1: Initialize SQLite database on first run
	t.Run("initialize database on first run", func(t *testing.T) {
		db, err := Open(ctx, Config{
			Path:    dbPath,
			WALMode: true,
		})
		require.NoError(t, err)
		defer db.Close()

		// Verify database was created
		assert.FileExists(t, dbPath)

		// Verify schema version is set
		version, err := db.GetSchemaVersion(ctx)
		require.NoError(t, err)
		assert.Equal(t, CurrentSchemaVersion, version)
	})

	// Requirement 2.6: Perform automatic migrations when schema changes
	t.Run("automatic migrations on schema changes", func(t *testing.T) {
		db, err := Open(ctx, Config{
			Path:    dbPath,
			WALMode: true,
		})
		require.NoError(t, err)
		defer db.Close()

		// Verify migrations were applied
		var count int
		err = db.Conn().QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, len(migrations), count)

		// Verify all tables exist
		tables := []string{"example_items"}
		for _, table := range tables {
			var tableCount int
			query := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`
			err := db.Conn().QueryRowContext(ctx, query, table).Scan(&tableCount)
			require.NoError(t, err)
			assert.Equal(t, 1, tableCount, "table %s should exist", table)
		}
	})

	// Test all table operations
	t.Run("example_items table operations", func(t *testing.T) {
		db, err := Open(ctx, Config{
			Path:    dbPath,
			WALMode: true,
		})
		require.NoError(t, err)
		defer db.Close()

		// Insert item
		_, err = db.Conn().ExecContext(ctx, `
			INSERT INTO example_items (key, value)
			VALUES (?, ?)
		`, "test-key", "test-value")
		require.NoError(t, err)

		// Query item
		var key, value string
		err = db.Conn().QueryRowContext(ctx, `
			SELECT key, value FROM example_items WHERE key = ?
		`, "test-key").Scan(&key, &value)
		require.NoError(t, err)
		assert.Equal(t, "test-key", key)
		assert.Equal(t, "test-value", value)
	})

	// Test indexes are working
	t.Run("verify indexes improve query performance", func(t *testing.T) {
		db, err := Open(ctx, Config{
			Path:    dbPath,
			WALMode: true,
		})
		require.NoError(t, err)
		defer db.Close()

		// Query using indexed column should use the index
		rows, err := db.Conn().QueryContext(ctx, `
			EXPLAIN QUERY PLAN
			SELECT * FROM example_items WHERE key = 'test-key'
		`)
		require.NoError(t, err)
		defer rows.Close()

		// The query plan should mention the index
		foundIndex := false
		for rows.Next() {
			var id, parent, notused int
			var detail string
			err := rows.Scan(&id, &parent, &notused, &detail)
			require.NoError(t, err)
			if contains(detail, "sqlite_autoindex") {
				foundIndex = true
				break
			}
		}
		assert.True(t, foundIndex, "query should use sqlite_autoindex index")
	})
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
