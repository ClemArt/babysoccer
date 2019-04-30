package main

import (
	"context"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"encoding/json"
)
// GraphMutation represents a single mutation to apply to the database on application startup
// Each GraphMutation will be run exactly once and should not be modified after this first run
type GraphMutation interface {
	mutate(context.Context, *dgo.Dgraph) error
	label() string
}

func RunMutation(c *dgo.Dgraph, gm GraphMutation) error {
	ctx := context.Background()
	t := c.NewTxn()
	defer t.Discard(ctx)

	if ran, err := checkPreviousRun(ctx, t, gm); err != nil {
		return errors.Wrapf(err, "Search for previous run of %s", gm.label())
	} else if ran {
		return nil
	}

	log.Debugf("Starting mutation %s", gm.label())
	if err := gm.mutate(ctx, c); err != nil {
		return errors.Wrapf(err, "Mutation %s", gm.label())
	}

	if err := saveExecution(ctx, t, gm); err != nil {
		return errors.Wrapf(err, "Saving execution of %s", gm.label())
	}

	return errors.Wrapf(t.Commit(ctx), "Commit mutation %s", gm.label())
}

func checkPreviousRun(ctx context.Context, t *dgo.Txn, gm GraphMutation) (bool, error) {
	log.Debugf("Searching previous execution of mutation %s", gm.label())

	var executed struct {
		Mutation []struct{ Uid string }
	}
	variables := make(map[string]string)
	variables["$a"] = gm.label()
	if resp, err := t.QueryWithVars(ctx, `query Mutation($a: string){
		mutation(func: eq(internals.migrations.executed, $a)) {
			uid
		}
	}`, variables); err != nil {
		return false, errors.Wrap(err, "Query")
	} else {
		if err := json.Unmarshal(resp.GetJson(), &executed); err != nil {
			return false, errors.Wrap(err, "Unmarshal")
		} else if len(executed.Mutation) > 0 {
			log.Debugf("Mutation already executed at uid %s", executed.Mutation[0].Uid)
			return true, nil
		}
	}
	return false, nil
}

func saveExecution(ctx context.Context, t *dgo.Txn, gm GraphMutation) error {
	type saveExecution struct {
		Executions []string `json:"internals.migrations.executed"`
	}
	s, errMarshal := json.Marshal(&saveExecution{
		Executions: []string{gm.label()},
	})
	if errMarshal != nil {
		return errors.Wrap(errMarshal, "Marshal")
	}

	_, errCommit := t.Mutate(ctx, &api.Mutation{
		SetJson: s,
	})
	if errCommit != nil {
		return errors.Wrap(errCommit, "Save to database")
	}
	return nil
}
