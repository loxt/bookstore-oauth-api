package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
)

var (
	cluster *gocql.ClusterConfig
)

func init() {
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	if err != nil {
		panic(err)
	}

	fmt.Println("cassandra connection successfully initialized")

	defer session.Close()
}

func GetSession() (*gocql.Session, error) {
	return cluster.CreateSession()
}
