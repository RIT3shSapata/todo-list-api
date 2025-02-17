package couchbase

import (
	"github.com/couchbase/gocb/v2"
)

type CouchbaseConfig struct {
	Host     string `env:"COUCHBASE_HOST"`
	Bucket   string `env:"COUCHBASE_BUCKET"`
	User     string `env:"COUCHBASE_API_USER"`
	Password string `env:"COUCHBASE_API_PASS"`
}

type Cluster interface {
	Query(statement string, opts *gocb.QueryOptions) (*gocb.QueryResult, error)
	BucketDefaultCol(bucketname string) Collection
}

type Collection interface {
	Get(id string, opts *gocb.GetOptions) (*gocb.GetResult, error)
	Upsert(id string, val interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error)
	Remove(id string, opts *gocb.RemoveOptions) (*gocb.MutationResult, error)
}
