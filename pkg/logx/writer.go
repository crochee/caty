// Package logx
package logx

import (
	"context"
	"time"

	"github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLogger struct {
	client *mongo.Client
}

func (m *MongoLogger) Write(p []byte) (int, error) {
	cli := m.client.Database("logx").Collection(time.Now().Local().Format("20060102"))
	var data map[string]interface{}
	if err := jsoniter.ConfigFastest.Unmarshal(p, &data); err != nil {
		return 0, err
	}
	if _, err := cli.InsertOne(context.Background(), data); err != nil {
		return 0, err
	}
	return len(p), nil
}
