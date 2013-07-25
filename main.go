package main

import (
	"bytes"
	"fmt"
	"github.com/ceh/gounqlite"
	"sync"
	"time"
)

var (
	db    *gounqlite.Handle
	mutex sync.Mutex
)

func init() {
	mutex.Lock()
	defer mutex.Unlock()

	if db != nil {
		return
	}

	// 如果open定义为":mem:"字符，数据存储在内存中;
	// 指定路径文件，数据会存储在文件中。
	udb, err := gounqlite.Open(":mem:")
	if err != nil {
		fmt.Println("Open: ", err.Error())
		return
	}

	db = udb
}

// 测试数据
var kvs = []struct {
	key   []byte
	value []byte
}{
	{[]byte("name"), []byte("viney")},
	{[]byte("国家"), []byte("中国")},
	{[]byte("email"), []byte("viney.chow@gmail.com")},
}

// 性能测试
func benchmark() {
	count := 10
	finish := make(chan bool)

	t := time.Now()
	for i := 0; i < count; i++ {
		go func(i int) {
			defer func() { finish <- true }()
			byt := []byte(fmt.Sprint(i))
			if err := db.Store(byt, byt); err != nil {
				fmt.Println("benchmark: ", err.Error())
				return
			}
		}(i)
	}

	for i := 0; i < count; i++ {
		<-finish
	}

	fmt.Println(time.Now().Sub(t).String())
}

func main() {
	// 性能测试
	benchmark()

	// 添加数据/修改数据
	// 如果key有数据会修改数据, 否则添加数据
	for _, v := range kvs {
		// 添加数据
		if err := db.Store(v.key, v.value); err != nil {
			fmt.Println("insert: ", err.Error(), string(v.key), string(v.value))
			return
		}

		// 查询
		if value, err := db.Fetch(v.key); err != nil {
			fmt.Println("Fetch: ", err.Error())
			return
		} else if !bytes.Equal(v.value, value) {
			fmt.Println("Equal: ", string(v.value), string(value))
			return
		}

		// 修改数据
		var hello []byte = []byte("hello")
		if err := db.Store(v.key, hello); err != nil {
			fmt.Println("update: ", err.Error())
			return
		}

		// 查询修改之后的数据
		if value, err := db.Fetch(v.key); err != nil {
			fmt.Println("Fetch: ", err.Error())
			return
		} else if !bytes.Equal(hello, value) {
			fmt.Println("Equal: ", string(hello), string(value))
			return
		}

		// 追加数据
		var world []byte = []byte("world")
		if err := db.Append(v.key, world); err != nil {
			fmt.Println("Append: ", err.Error())
			return
		}

		// 查询追加之后的数据
		if value, err := db.Fetch(v.key); err != nil {
			fmt.Println("Fetch: ", err.Error())
			return
		} else if value != nil {
			hello = append(hello, world...)
			if !bytes.Equal(value, hello) {
				fmt.Println("Equal: ", string(value), string(hello))
				return
			}
		}

		// 删除
		if err := db.Delete(v.key); err != nil {
			fmt.Println("Delete: ", err.Error())
			return
		}
	}

	fmt.Println("version: ", gounqlite.Version())

	// 关闭open
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Close: ", err.Error())
			return
		}
	}()
}
