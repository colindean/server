// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"gorm.io/gorm"
)

// engine represents the pipeline functionality that implements the PipelineService interface.
type engine struct {
	client *gorm.DB

	// specifies the level of compression to use for the Data field in a Pipeline.
	compressionLevel int
}

// New creates and returns a Vela service for integrating with pipelines in the database.
//
// nolint: revive // ignore returning unexported engine
func New(client *gorm.DB, level int) *engine {
	return &engine{
		client:           client,
		compressionLevel: level,
	}
}
