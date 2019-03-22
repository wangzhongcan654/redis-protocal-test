package common

import (
	"strconv"
	"fmt"
	"bufio"
	"bytes"
	"github.com/juju/errors"
	"strings"
)

/*
用单行回复，回复的第一个字节将是“+”
错误消息，回复的第一个字节将是“-”
整型数字，回复的第一个字节将是“:”
批量回复，回复的第一个字节将是“$”
多个批量回复，回复的第一个字节将是“*”
 */
/*
编写程序，要求能够正确解析redis协议的字符串，比如
当输入"+OK\r\n时"，应该解析成结果"OK"，当输入
":1000\r\n"应该解析成1000
简短的文档说明
 */
var(
	simpleStrPrefix = []byte{'+'}
	errPrefix       = []byte{'-'}
	bulkStrPrefix   = []byte{'$'}
	nilBulkString   = []byte("$-1\r\n")
	delimiter       = []byte{'\r','\n'}
)

const(
	defaultReadPacketSize =4*1024
	defaultWriterPacketSize = 64
)

//序列化string
func MarshalRESP(b []byte) []byte {
	scratch := make([]byte,0,64)
	scratch = append(scratch, bulkStrPrefix...)
	scratch = strconv.AppendInt(scratch, int64(len(b)), 10)
	scratch = append(scratch,delimiter...)
	scratch = append(scratch, b...)
	scratch = append(scratch, delimiter...)
	return scratch
}

func ParseUint(b []byte) (uint64, error) {
	if len(b) == 0 {

		return 0, errors.New("empty slice given to parseUint")
	}
	var n uint64

	for i, c := range b {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid character %c at position %d in parseUint", c, i)
		}
		n *= 10
		n += uint64(c - '0')
	}
	return n, nil
}
func ParseInt(b []byte) (int64, error) {
	if len(b) == 0 {
		return 0, errors.New("empty slice given to parseInt")
	}

	var neg bool
	if b[0] == '-' || b[0] == '+' {
		neg = b[0] == '-'
		b = b[1:]
	}

	n, err := ParseUint(b)
	if err != nil {
		return 0, err
	}

	if neg {
		return -int64(n), nil
	}

	return int64(n), nil
}
func BufferedPrefix(br *bufio.Reader, prefix []byte) error {
	b, err := br.Peek(len(prefix))
	if err != nil {
		return err
	} else if !bytes.Equal(b, prefix) {
		return fmt.Errorf("expected prefix %q, got %q", prefix, b)
	}
	_, err = br.Discard(len(prefix))
	return err
}

func BufferedBytesDelim(br *bufio.Reader) ([]byte, error) {
	b, err := br.ReadSlice('\n')
	if err != nil {
		return nil, err
	} else if len(b) < 2 || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("malformed resp %q", b)
	}
	return b[:len(b)-2], err
}



func UnmarshalRESP(br *bufio.Reader) (string,error) {
	if err := BufferedPrefix(br, simpleStrPrefix); err != nil {
		return "",err
	}
	b, err := BufferedBytesDelim(br)
	if err != nil {
		return "",err
	}

	str:= string(b)
	return str,nil
}
//proro covert args
func RedisProCovertArgs(protocol string) (argv []string, argc int) {
	//协议拆解，\r\n
	parts := strings.Split(strings.Trim(protocol, " "), "\r\n")
	if len(parts) == 0 {
		return nil, 0
	}
	argc, err := strconv.Atoi(parts[0][1:])
	if err != nil {
		return nil, 0
	}
	j := 0
	var vlen []int

	for _, v := range parts{
		if len(v) == 0 {
			continue
		}

		if v[0] == '$' {
			tmpl, err := strconv.Atoi(v[1:])
			if err == nil {
				vlen = append(vlen, tmpl)
			}
		} else {
			if j < len(vlen) && vlen[j] == len(v) {
				j++
				argv = append(argv, v)
			}
		}
	}
	return argv, argc
}