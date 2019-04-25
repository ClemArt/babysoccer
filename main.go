package main

import (
	"context"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	dgraphAddresses := []string{"192.168.99.100:30744"}

	clients := make([]api.DgraphClient, len(dgraphAddresses))

	for _, a := range dgraphAddresses {
		clients = append(clients, newClient(a))
	}

	clusterCLient := newClusterCLient(clients)

	if err := dropAllData(clusterCLient); err != nil {
		log.Fatal(err)
	}

	if err := setupSchema(clusterCLient); err != nil {
		log.Fatal(err)
	}
}

func newClient(address string) api.DgraphClient {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	log.Debugf("Initialize grpc connection to dgraph %s", address)
	d, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		errors.Wrap(err, "dgraph connection failed")
		log.Fatal(err)
	}

	return api.NewDgraphClient(d)
}

func newClusterCLient(clients []api.DgraphClient) *dgo.Dgraph {
	return dgo.NewDgraphClient(clients...)
}

func dropAllData(c *dgo.Dgraph) error {
	log.Warn("Data drop required !")
	err := c.Alter(context.Background(), &api.Operation{
		DropAll: true,
	})

	return errors.Wrap(err, "Can't drop data in the database")
}

func setupSchema(c *dgo.Dgraph) error {
	log.Info("Setup database schema")
	err := c.Alter(context.Background(), &api.Operation{
		Schema:`
			name: string @index(hash) @upsert .
			in_team: uid @reverse .
			won_match: uid @count @reverse .
		`,
	})
	return errors.Wrap(err, "Schema setup failed")
}
