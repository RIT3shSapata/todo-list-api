package couchbase

import "github.com/couchbase/gocb/v2"

type collection struct {
	collection *gocb.Collection
}

var _ Collection = &collection{}

func (col *collection) Get(id string, opts *gocb.GetOptions) (*gocb.GetResult, error) {
	return col.collection.Get(id, opts)
}

func (col *collection) Upsert(id string, val interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	return col.collection.Upsert(id, val, opts)
}

func (col *collection) Remove(id string, opts *gocb.RemoveOptions) (*gocb.MutationResult, error) {
	return col.collection.Remove(id, opts)
}
