// Copyright 2021, crochee.All rights reserved.
// Author: crochee
// Date: 2021/8/13

// Package logx
package logx

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/natefinch/lumberjack.v2"

	"obs/pkg/model/mongox"
)

type MongoLogger struct {
	client *mongo.Client
}

func (m *MongoLogger) Write(p []byte) (int, error) {
	cli := m.client.Database("log").Collection(time.Now().Local().Format("20060102"))
	var data map[string]interface{}
	if err := jsoniter.ConfigFastest.Unmarshal(p, &data); err != nil {
		return 0, err
	}
	if _, err := cli.InsertOne(context.Background(), data); err != nil {
		return 0, err
	}
	return len(p), nil
}

// SetLoggerWriter return a io.Writer
func SetLoggerWriter(path string) io.Writer {
	if path == "" {
		return os.Stdout
	}
	if strings.HasPrefix(path, "mongodb://") {
		client, err := mongox.Setup(path)
		if err != nil {
			panic(err)
		}
		return &MongoLogger{client: client}
	}
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    DefaultLogSizeM, //单个日志文件最大MaxSize*M大小 // megabytes
		MaxAge:     MaxLogDays,      //days
		MaxBackups: DefaultMaxZip,   //备份数量
		Compress:   false,           //不压缩
		LocalTime:  true,            //备份名采用本地时间
	}
}
