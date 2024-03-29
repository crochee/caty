// Author: crochee
// Date: 2021/9/6

// Package message
package message

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/crochee/lirity"
	"github.com/crochee/lirity/mq"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var router *message.Router

// Setup run message pub/sub or not
func Setup(ctx context.Context) error {
	if !viper.GetBool("rabbitmq.enable") {
		return nil
	}
	logger := watermill.NewStdLogger(viper.GetBool("debug"), gin.Mode() == gin.DebugMode)
	var err error
	router, err = message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return err
	}
	conn, err := mq.New(mq.WithURI(viper.GetString("rabbitmq.consumer.resource.URI")))
	if err != nil {
		return err
	}
	defer conn.Close()
	return router.Run(ctx)
}

func Close() {
	if router == nil {
		return
	}
	lirity.Close(router)
}
