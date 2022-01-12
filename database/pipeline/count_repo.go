// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// CountPipelinesForRepo gets the count of pipelines by repo ID from the database.
func (e *engine) CountPipelinesForRepo(r *library.Repo) (int64, error) {
	// TODO: figure this out
	//c.Logger.WithFields(logrus.Fields{
	//	"org":  r.GetOrg(),
	//	"repo": r.GetName(),
	//}).Tracef("getting count of pipelines for repo %s from the database", r.GetFullName())

	// variable to store query results
	var p int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TablePipeline).
		Where("repo_id = ?", r.GetID()).
		Count(&p).
		Error

	return p, err
}
