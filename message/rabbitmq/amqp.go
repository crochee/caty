// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/3

package rabbitmq

import "github.com/streadway/amqp"

func Setup() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return
	}
	defer conn.Close()

}
