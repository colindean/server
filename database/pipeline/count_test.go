// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPipeline_Engine_CountPipelines(t *testing.T) {
	// setup types
	_pipelineOne := testPipeline()
	_pipelineOne.SetID(1)
	_pipelineOne.SetRepoID(1)
	_pipelineOne.SetNumber(1)
	_pipelineOne.SetRef("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipelineOne.SetType("yaml")
	_pipelineOne.SetVersion("1")

	_pipelineTwo := testPipeline()
	_pipelineTwo.SetID(2)
	_pipelineTwo.SetRepoID(2)
	_pipelineTwo.SetNumber(1)
	_pipelineTwo.SetRef("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipelineTwo.SetType("yaml")
	_pipelineTwo.SetVersion("1")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "pipelines"`).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreatePipeline(_pipelineOne)
	if err != nil {
		t.Errorf("unable to create test pipeline for sqlite: %v", err)
	}

	err = _sqlite.CreatePipeline(_pipelineTwo)
	if err != nil {
		t.Errorf("unable to create test pipeline for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     2,
		},
		{
			failure:  false,
			name:     "sqlite",
			database: _sqlite,
			want:     2,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.database.CountPipelines()

		if test.failure {
			if err == nil {
				t.Errorf("CountPipelines for %s should have returned err", test.name)
			}

			continue
		}

		if err != nil {
			t.Errorf("CountPipelines for %s returned err: %v", test.name, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("CountPipelines for %s is %v, want %v", test.name, got, test.want)
		}
	}
}
