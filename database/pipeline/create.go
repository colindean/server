// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// CreatePipeline creates a new pipeline in the database.
func (e *engine) CreatePipeline(p *library.Pipeline) error {
	// TODO: figure this out
	//c.Logger.WithFields(logrus.Fields{
	//	"pipeline": p.GetRef(),
	//}).Tracef("creating pipeline %s in the database", p.GetRef())

	// cast to database type
	pipeline := database.PipelineFromLibrary(p)

	// validate the necessary fields are populated
	err := pipeline.Validate()
	if err != nil {
		return err
	}

	// compress data for the pipeline
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Log.Compress
	err = pipeline.Compress(e.compressionLevel)
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TablePipeline).
		Create(pipeline).
		Error
}
