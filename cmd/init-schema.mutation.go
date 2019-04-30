package main

import (
	"context"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
)

type InitSchemaMutation struct{}

func (m *InitSchemaMutation) mutate(ctx context.Context, c *dgo.Dgraph) error {
	return c.Alter(ctx, &api.Operation{
		Schema: `
			name: string @index(hash) @upsert .
			attacker: uid @reverse .
			defender: uid @reverse .
			mixed: uid @reverse .
			won_match: uid @count @reverse .
		`,
	})
}

func (m *InitSchemaMutation) label() string {
	return "InitSchemaMutation"
}
