// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// DeletePipeline deletes an existing pipeline from the database.
func (e *engine) DeletePipeline(p *library.Pipeline) error {
	// TODO: figure this out
	//c.Logger.WithFields(logrus.Fields{
	//	"pipeline": p.GetRef(),
	//}).Tracef("deleting pipeline %s in the database", p.GetRef())

	// cast to database type
	pipeline := database.PipelineFromLibrary(p)

	// send query to the database
	return e.client.
		Table(constants.TablePipeline).
		Delete(pipeline).
		Error
}
