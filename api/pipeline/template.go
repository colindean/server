// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/registry/github"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/pipelines/{org}/{repo}/{pipeline}/templates pipelines GetTemplates
//
// Get a map of templates utilized by a pipeline configuration from the database
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
//     description: Successfully retrieved the map of pipeline templates
//     schema:
//       "$ref": "#/definitions/Template"
//   '400':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve the pipeline configuration templates
//     schema:
//       "$ref": "#/definitions/Error"

// GetTemplates represents the API handler to capture a
// map of templates utilized by a pipeline configuration.
func GetTemplates(c *gin.Context) {
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
	}).Infof("reading templates from pipeline %s", entry)

	// create the compiler object
	compiler := compiler.FromContext(c).Duplicate().WithMetadata(m).WithRepo(r).WithUser(u)

	// parse the pipeline configuration
	pipeline, _, err := compiler.Parse(p.GetData())
	if err != nil {
		util.HandleError(c, http.StatusBadRequest,
			fmt.Errorf("unable to parse pipeline %s: %v", entry, err),
		)

		return
	}

	// send API call to capture the repo owner
	user, err := database.FromContext(c).GetUser(r.GetUserID())
	if err != nil {
		util.HandleError(c, http.StatusBadRequest,
			fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err),
		)

		return
	}

	templates := make(map[string]*library.Template)
	for name, template := range pipeline.Templates.Map() {
		templates[name] = template.ToLibrary()

		// create a compiler registry client for parsing (no address or token needed for Parse)
		registry, err := github.New("", "")
		if err != nil {
			util.HandleError(c, http.StatusBadRequest,
				fmt.Errorf("unable to create compiler github client: %w", err),
			)

			return
		}

		// parse the source for the template using the compiler registry client
		src, err := registry.Parse(template.Source)
		if err != nil {
			util.HandleError(c, http.StatusBadRequest,
				fmt.Errorf("unable to parse source for %s: %w", template.Source, err),
			)

			return
		}

		// retrieve link to template file from github
		link, err := scm.FromContext(c).GetHTMLURL(user, src.Org, src.Repo, src.Name, src.Ref)
		if err != nil {
			util.HandleError(c, http.StatusBadRequest,
				fmt.Errorf("unable to get html url for %s/%s/%s/@%s: %w", src.Org, src.Repo, src.Name, src.Ref, err),
			)

			return
		}

		// set link to file for template
		templates[name].SetLink(link)
	}

	writeOutput(c, templates)
}
