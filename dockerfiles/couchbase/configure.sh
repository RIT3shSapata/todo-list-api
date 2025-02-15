#!/bin/bash

# exit immediately if a command fails, or a command in a pipeline fails, or if there are unset
# variables
set -euo pipefail

# turn on bash's job control, used to bring couchbase-server back to the forground after the node
# is configured
set -m

/entrypoint.sh couchbase-server &

until curl -s http://localhost:8091/pools >/dev/null; do
  sleep 5
done

# check if cluster is already initialized
if ! couchbase-cli server-list -c localhost:8091 -u $COUCHBASE_ADMIN_USER -p $COUCHBASE_ADMIN_PASS >/dev/null; then
    couchbase-cli cluster-init -c localhost \
    --cluster-username $COUCHBASE_ADMIN_USER \
    --cluster-password $COUCHBASE_ADMIN_PASS \
    --services data,index,query \
    --cluster-ramsize 512 \
    --cluster-index-ramsize 256
fi

bucket_create() {
  if ! couchbase-cli bucket-list -c localhost:8091 -u $COUCHBASE_ADMIN_USER -p $COUCHBASE_ADMIN_PASS | grep $1 >/dev/null; then
    couchbase-cli bucket-create -c localhost:8091 \
        -u $COUCHBASE_ADMIN_USER \
        -p $COUCHBASE_ADMIN_PASS \
        --bucket $1 --bucket-type couchbase --bucket-ramsize 100 --bucket-replica 0 --enable-flush 1 --wait
  fi
}

scope_create() {
  if ! couchbase-cli collection-manage -c localhost:8091 -u $COUCHBASE_ADMIN_USER -p $COUCHBASE_ADMIN_PASS --bucket $1 --list-scopes | grep $2 >/dev/null; then
    couchbase-cli collection-manage -c localhost:8091 \
    -u $COUCHBASE_ADMIN_USER \
    -p $COUCHBASE_ADMIN_PASS \
    --bucket $1 --create-scope $2
  fi
}

collection_create() {
  if ! couchbase-cli collection-manage -c localhost:8091 -u $COUCHBASE_ADMIN_USER -p $COUCHBASE_ADMIN_PASS --bucket $1 --list-collections | grep $3 >/dev/null; then
    couchbase-cli collection-manage -c localhost:8091 \
      -u $COUCHBASE_ADMIN_USER \
      -p $COUCHBASE_ADMIN_PASS \
      --bucket $1 --create-collection "$2.$3"
  fi
}

bucket_create $COUCHBASE_BUCKET

fg 1