package test

import (
	"fmt"
	"jam3.com/user/pgk/dao"
	"testing"
	"time"
)

func TestKfkInit(t *testing.T) {

	client, err := dao.NewClient([]string{"192.168.3.116:9092"}, "test", "goer")
	fmt.Println(err)

	err = client.SendSync("hello...." + time.Now().String())
	fmt.Println(err)
}
func TestKfkConsumer(t *testing.T) {
	dao.Consumer([]string{"192.168.3.116:9092"}, "test")
}
