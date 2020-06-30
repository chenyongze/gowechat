/*
 * @Author: your name
 * @Date: 2020-06-30 20:02:12
 * @LastEditTime: 2020-06-30 20:12:35
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /gowechat/util/string_test.go
 */
package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	//abc sig
	abc := "a100"
	if abc != ToStr(abc) {
		t.Error("test int to string Error")
	}

	fmt.Println(reflect.TypeOf(ToStr(abc)))
	fmt.Println(ToInt64("200"))
	fmt.Println(RandomStr(10))
}
