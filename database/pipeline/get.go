// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"errors"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm"
)

// GetPipeline gets a pipeline by ref and repo ID from the database.
func (e *engine) GetPipeline(ref string, r *library.Repo) (*library.Pipeline, error) {
	// TODO: figure this out
	//c.Logger.WithFields(logrus.Fields{
	//	"pipeline": ref,
	//	"org":      r.GetOrg(),
	//	"repo":     r.GetName(),
	//}).Tracef("getting pipeline %s/%s from the database", r.GetFullName(), ref)

	// variable to store query results
	p := new(database.Pipeline)

	// send query to the database and store result in variable
	result := e.client.
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Where("ref = ?", ref).
		Scan(p)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// decompress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Pipeline.Decompress
	err := p.Decompress()
	if err != nil {
		return nil, err
	}

	// return the decompressed pipeline
	return p.ToLibrary(), result.Error
}
