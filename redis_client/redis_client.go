package redis_client

import (
	"net"
	"bufio"
	"os"
	"fmt"

	"io/ioutil"
	"redis-protocal-test/common"
)

var(
	conn net.Conn
)
func init(){
	conn,_=net.Dial("tcp","106.75.31.54:6379")
}

func Test(){
	input:=bufio.NewScanner(os.Stdin)

	for{
		input.Scan()
		fmt.Println(input.Text())
		a:=[]byte(input.Text())
		args:=common.MarshalRESP(a)

		n,err:=conn.Write(args)
		if err!=nil{
			fmt.Println("insert failed")
		}else{
			fmt.Println("insert success",n)
		}

		data,err:=ioutil.ReadAll(conn)
		if err!=nil{
			fmt.Println("read error",err)
			continue
		}
		fmt.Println(string(data))
		/*
		buf:=bufio.NewReader(conn)
		str,err:=common.UnmarshalRESP(buf)
		if err!=nil{
			fmt.Println(err)
			continue
		}
		fmt.Println(str)
		*/
	}
}