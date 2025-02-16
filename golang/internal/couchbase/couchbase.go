package couchbase

type CouchbaseConfig struct {
	Host     string `env:"COUCHBASE_HOST"`
	Bucket   string `env:"COUCHBASE_BUCKET"`
	User     string `env:"COUCHBASE_API_USER"`
	Password string `env:"COUCHBASE_API_PASS"`
}
