package handler

import (
	"fmt"
	"strings"
	"testing"
)

func Test_CacheReader(t *testing.T) {
	r := strings.NewReader("hello world")
	c := CacheReader{
		reader: r,
	}
	a := make([]byte, 7)
	n, err := c.Read(a) // 返回读取到字节数
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(n, a)
	fmt.Println("read  ===>", string(a))
	fmt.Println("cache ===>", string(c.Cache()))
}
