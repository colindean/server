// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/yaml"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/pipelines/{org}/{repo}/{pipeline}/expand pipelines ExpandPipeline
//
// Get and expand a pipeline configuration from the database
//
// ---
// produces:
// - application/x-yaml
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
//     description: Successfully retrieved and expanded the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/PipelineBuild"
//   '400':
//     description: Unable to expand the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration
//     schema:
//       "$ref": "#/definitions/Error"

// ExpandPipeline represents the API handler to capture and
// expand a pipeline configuration.
func ExpandPipeline(c *gin.Context) {
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
	}).Infof("expanding templates for pipeline %s", entry)

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

	// expand the parsed pipeline configuration
	err = expandPipeline(compiler, pipeline, entry, false)
	if err != nil {
		util.HandleError(c, http.StatusBadRequest, err)

		return
	}

	writeOutput(c, pipeline)
}

// expandPipeline is a helper function to expand the stages or steps in the provided
// pipeline with the provided compiler. A boolean is provided to control if the
// environment variables will be substituted in those stages or steps.
func expandPipeline(c compiler.Engine, p *yaml.Build, entry string, substituteEnv bool) error {
	var err error

	// create map of templates for easy lookup
	templates := p.Templates.Map()

	if len(p.Stages) > 0 {
		// inject the templates into the stages
		p.Stages, p.Secrets, p.Services, p.Environment, err = c.ExpandStages(p, templates)
		if err != nil {
			return fmt.Errorf("unable to expand stages for pipeline %s: %w", entry, err)
		}

		if substituteEnv {
			// inject the substituted environment variables into the stages
			p.Stages, err = c.SubstituteStages(p.Stages)
			if err != nil {
				return fmt.Errorf("unable to substitute stages for pipeline %s: %w", entry, err)
			}
		}
	} else {
		// inject the templates into the steps
		p.Steps, p.Secrets, p.Services, p.Environment, err = c.ExpandSteps(p, templates)
		if err != nil {
			return fmt.Errorf("unable to expand steps for pipeline %s: %w", entry, err)
		}

		if substituteEnv {
			// inject the substituted environment variables into the steps
			p.Steps, err = c.SubstituteSteps(p.Steps)
			if err != nil {
				return fmt.Errorf("unable to substitute steps for pipeline %s: %w", entry, err)
			}
		}
	}

	return nil
}
