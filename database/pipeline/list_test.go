// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestPipeline_Engine_ListPipelines(t *testing.T) {
	// setup types
	_pipelineOne := testPipeline()
	_pipelineOne.SetID(1)
	_pipelineOne.SetRepoID(1)
	_pipelineOne.SetNumber(1)
	_pipelineOne.SetRef("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipelineOne.SetType("yaml")
	_pipelineOne.SetVersion("1")
	_pipelineOne.SetData([]byte("foo"))

	_pipelineTwo := testPipeline()
	_pipelineTwo.SetID(2)
	_pipelineTwo.SetRepoID(1)
	_pipelineTwo.SetNumber(2)
	_pipelineTwo.SetRef("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipelineTwo.SetType("yaml")
	_pipelineTwo.SetVersion("1")
	_pipelineTwo.SetData([]byte("foo"))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "flavor", "platform", "ref", "type", "version", "services", "stages", "steps", "templates", "data"}).
		AddRow(1, 1, 1, "", "", "48afb5bdc41ad69bf22588491333f7cf71135163", "yaml", "1", false, false, false, false, []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}).
		AddRow(2, 1, 2, "", "", "48afb5bdc41ad69bf22588491333f7cf71135163", "yaml", "1", false, false, false, false, []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "pipelines"`).
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
		want     []*library.Pipeline
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Pipeline{_pipelineOne, _pipelineTwo},
		},
		{
			failure:  false,
			name:     "sqlite",
			database: _sqlite,
			want:     []*library.Pipeline{_pipelineOne, _pipelineTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.database.ListPipelines()

		if test.failure {
			if err == nil {
				t.Errorf("ListPipelines for %s should have returned err", test.name)
			}

			continue
		}

		if err != nil {
			t.Errorf("ListPipelines for %s returned err: %v", test.name, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ListPipelines for %s is %v, want %v", test.name, got, test.want)
		}
	}
}
