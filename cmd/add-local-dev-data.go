package main

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type AddLocalDevData struct{}

func (m *AddLocalDevData) mutate(ctx context.Context, c *dgo.Dgraph) error {
	log.Debug("Searching previous execution of mutation AddLocalDevData")
	t := c.NewTxn()
	defer t.Discard(ctx)

	resp, err := t.Query(ctx, `{
		mutation(func: eq(internals.migrations.executed, "AddLocalDevData")) {
			uid
		}
	}`)

	if err != nil {
		return errors.Wrap(err, "Unable to retrieve previous run informations")
	}

	var executed struct {
		Mutation []struct{ Uid string }
	}
	if err := json.Unmarshal(resp.GetJson(), &executed); err != nil {
		return errors.Wrap(err, "Unmarshal of previous runs query failed")
	} else if len(executed.Mutation) > 0 {
		log.Debugf("Mutation already executed at uid %s", executed.Mutation[0].Uid)
		return nil
	}

	log.Debug("Starting mutation AddLocalDevData")
	if err := c.Alter(ctx, &api.Operation{
		Schema: `
			name: string @index(hash) @upsert .
			attacker: uid @reverse .
			defender: uid @reverse .
			mixed: uid @reverse .
			won_match: uid @count @reverse .
		`,
	}); err != nil {
		return errors.Wrap(err, "Mutation AddLocalDevData failed")
	}

	type saveExecution struct {
		Executions []string `json:"internals.migrations.executed"`
	}
	s, errMarshal := json.Marshal(&saveExecution{
		Executions: []string{"AddLocalDevData"},
	})
	if errMarshal != nil {
		return errors.Wrap(errMarshal, "Unable to marshal mutation AddLocalDevData for commit")
	}

	_, errCommit := t.Mutate(ctx, &api.Mutation{
		SetJson: s,
	})
	if errCommit != nil {
		return errors.Wrap(errCommit, "Unable to prepare mutation AddLocalDevData for commit")
	}

	errEndTx := t.Commit(ctx)
	return errors.Wrap(errEndTx, "Unable to commit mutation AddLocalDevData")
}
