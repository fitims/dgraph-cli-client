package dgraph

import (
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
	"log"
)

var Client *dgo.Dgraph

var connection *grpc.ClientConn
var dgClient api.DgraphClient

// Open opens a separate connection for each url provided
func Open(url string) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Error connecting to DGraph : ", err)
		panic("Error connecting to DGraph")
	}

	connection = conn
	dgClient = api.NewDgraphClient(conn)
	Client = dgo.NewDgraphClient(dgClient)
}

// Close will close all the connections
func Close() {
	connection.Close()
}
