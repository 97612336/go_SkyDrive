package main

import (
	"time"
	"reflect"
	"fmt"
)

func main() {
	now := time.Now()

	now2:=now.Format("2006-01-02 15:04:05")

	type_str := reflect.TypeOf(now2)
	fmt.Println(type_str)
	fmt.Println(now)
	fmt.Println(now2)
}
