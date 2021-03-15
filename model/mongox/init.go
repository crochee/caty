// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/15

package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup init mongo.Client
func Setup(path string) (*mongo.Client, error) {
	// "mongodb://user:pass@host:port"
	return mongo.Connect(context.Background(), options.Client().ApplyURI(path).SetMaxPoolSize(100))
}
