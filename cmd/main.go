package main

import (
	"context"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	dgraphAddresses := []string{"192.168.99.100:30980"}
	mutationsToRun := []GraphMutation{
		&InitSchemaMutation{},
	}

	client := openDgraphConnection(dgraphAddresses)

	// if err := dropAllData(client); err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.Background()
	client.Alter(ctx, &api.Operation{
		Schema: `
			internals.migrations.executed: string @index(hash) @upsert .
		`,
	})
	for _, m := range mutationsToRun {
		if e := RunMutation(client, m); e != nil {
			log.Fatal(e)
		}
	}
}

func openDgraphConnection(addresses []string) *dgo.Dgraph {
	clients := make([]api.DgraphClient, len(addresses))

	for _, a := range addresses {
		clients = append(clients, newClient(a))
	}

	return newClusterCLient(clients)
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
