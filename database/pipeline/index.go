package pipeline

const (
	// CreateRepoIDIndex represents a query to create an
	// index on the pipelines table for the repo_id column.
	CreateRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
pipelines_repo_id
ON pipelines (repo_id);
`
)

// CreateIndexes does stuff...
func (e *engine) CreateIndexes() error {
	// TODO: figure this out
	//c.Logger.Tracef("creating pipelines table in the database")

	// create the repo_id column index for the pipelines table
	return e.client.Exec(CreateRepoIDIndex).Error
}
