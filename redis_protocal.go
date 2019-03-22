package redis_client

import (
	"strconv"
	"strings"
	"errors"
)
/*
用单行回复，回复的第一个字节将是“+”
错误消息，回复的第一个字节将是“-”
整型数字，回复的第一个字节将是“:”
批量回复，回复的第一个字节将是“$”
多个批量回复，回复的第一个字节将是“*”
 */

var(
	simpleStrPrefix = []byte{'+'}
	errPrefix       = []byte{'-'}
	bulkStrPrefix   = []byte{'$'}
	nilBulkString   = []byte("$-1\r\n")
	delimiter       = []byte{'\r','\n'}
)
var(
	LineRes = "+"
	ErrRES = "-"
	IntRes = ":"
	BatchRes="$"
	MoreBatchRes="*"
)
//args convert redis proto
func MarshalRESP(b []byte) []byte {
	scratch := make([]byte,0,64)
	scratch = append(scratch, bulkStrPrefix...)
	scratch = strconv.AppendInt(scratch, int64(len(b)), 10)
	scratch = append(scratch,delimiter...)
	scratch = append(scratch, b...)
	scratch = append(scratch, delimiter...)
	return scratch
}


func RedisProCovertRes(protocol string) (interface{}, error) {
	//协议拆解，\r\n
	parts := strings.Split(strings.Trim(protocol, " "), "\r\n")
	if len(parts) == 0 {
		return "", errors.New("string is empty")
	}
	var str string
	for _, v := range parts{
		if len(v) == 0 {
			continue
		}
		if v[0]==':'{
			tmpl, err := strconv.Atoi(v[1:])
			if err != nil {
				return "",err
			}
			return tmpl,nil
		} else if v[0] == '+'||v[0]=='-'{
			tmpl:=v[1:]
			return tmpl,nil
		}else if v[0]=='$'{
			str=str+v[1:]
			continue
		}else if v[0]=='*'{
			//
		} else{
			return nil,errors.New("Unknown character,Parsing failure")
		}
	}
	return str,nil
}

