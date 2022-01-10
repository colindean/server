// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/{pipeline}/validate pipelines ValidatePipeline
//
// Get, expand and validate a pipeline configuration from the database
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: query
//   name: ref
//   description: Ref for retrieving pipeline configuration file
//   type: string
// - in: query
//   name: output
//   description: Output string for specifying output format
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved, expanded and validated the pipeline
//     schema:
//       type: string
//   '400':
//     description: Unable to validate the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"

// ValidatePipeline represents the API handler to capture,
// expand and validate a pipeline configuration.
func ValidatePipeline(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*types.Metadata)
	o := org.Retrieve(c)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%s", r.GetFullName(), p.GetRef())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"pipeline": p.GetRef(),
		"repo":     r.GetName(),
		"user":     u.GetName(),
	}).Infof("validating pipeline %s", entry)

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithMetadata(m).WithRepo(r).WithUser(u)

	// parse the pipeline configuration
	pipeline, err := compiler.Parse(p.GetData())
	if err != nil {
		util.HandleError(c, http.StatusBadRequest,
			fmt.Errorf("unable to parse pipeline %s: %v", entry, err),
		)

		return
	}

	// check optional template query parameter
	if ok, _ := strconv.ParseBool(c.DefaultQuery("template", "true")); ok {
		// expand the parsed pipeline configuration
		err = expandPipeline(compiler, pipeline, entry, false)
		if err != nil {
			util.HandleError(c, http.StatusBadRequest, err)

			return
		}
	}

	// validate the parsed and possibly expanded pipeline configuration
	err = compiler.Validate(pipeline)
	if err != nil {
		util.HandleError(c, http.StatusBadRequest,
			fmt.Errorf("unable to validate pipeline %s: %v", entry, err),
		)

		return
	}

	writeOutput(c, pipeline)
}
