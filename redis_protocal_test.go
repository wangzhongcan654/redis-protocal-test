package redis_client

import (
	"testing"
	"fmt"
)

func TestRedisClient(t *testing.T){
	/*
	str:="*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	n,m:=RedisProCovertArgs(str)
	fmt.Println(n,":",m)
	fmt.Println("==============")
	*/
	str:="+OK\r\n"
	n,_:=RedisProCovertRes(str)
	fmt.Println(n.(string))
	str=":1000\r\n"
	n,_=RedisProCovertRes(str)
	fmt.Println(n)
}