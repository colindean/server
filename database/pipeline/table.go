// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import "github.com/go-vela/types/constants"

const (
	// CreatePostgresTable represents a query to create the Postgres pipelines table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
pipelines (
	id        SERIAL PRIMARY KEY,
	repo_id   INTEGER,
	flavor    VARCHAR(100),
	platform  VARCHAR(100),
	ref       VARCHAR(500),
	version   VARCHAR(50),
	services  BOOLEAN,
	stages    BOOLEAN,
	steps     BOOLEAN,
	templates BOOLEAN,
	data      BYTEA,
	UNIQUE(repo_id, ref)
);
`

	// CreateSqliteTable represents a query to create the Sqlite pipelines table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
pipelines (
	id        INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id   INTEGER,
	flavor    TEXT,
	platform  TEXT,
	ref       TEXT,
	version   TEXT,
	services  BOOLEAN,
	stages    BOOLEAN,
	steps     BOOLEAN,
	templates BOOLEAN,
	data      BLOB,
	UNIQUE(repo_id, ref)
);
`
)

// CreateTable does stuff...
func (e *engine) CreateTable(driver string) error {
	// TODO: figure this out
	//c.Logger.Tracef("creating pipelines table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the pipelines table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the pipelines table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
