package couchbase

import (
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
)

type cluster struct {
	cluster *gocb.Cluster
}

var _ Cluster = &cluster{}

func NewCluster(cfg CouchbaseConfig) (*cluster, error) {
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: cfg.User,
			Password: cfg.Password,
		},
	}

	couchbaseCluster, err := gocb.Connect(cfg.Host, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create couchbase client: %w", err)
	}

	err = couchbaseCluster.WaitUntilReady(time.Second*10, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to couchbase cluster: %w", err)
	}

	c := &cluster{
		cluster: couchbaseCluster,
	}

	return c, nil
}

func (cl *cluster) BucketDefaultCol(bucketname string) Collection {
	bucket := cl.cluster.Bucket(bucketname)
	defaultCollection := bucket.DefaultCollection()
	couchbaseCollection := &collection{
		collection: defaultCollection,
	}
	return couchbaseCollection
}

func (cl *cluster) Query(statement string, opts *gocb.QueryOptions) (*gocb.QueryResult, error) {
	queryResult, err := cl.cluster.Query(statement, opts)
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}
