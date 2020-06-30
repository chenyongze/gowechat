/*
 * @Author: your name
 * @Date: 2020-06-30 20:14:18
 * @LastEditTime: 2020-06-30 20:31:26
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /gowechat/util/http_test.go
 */
package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Infox struct {
	Info string `json:"info"`
}

func TestGet(t *testing.T) {

	var response []byte
	response, _ = HTTPGet("http://go.dotalk.cn")

	fmt.Println(response)

	ticket := new(Infox)
	// var ticket Infox
	err := json.Unmarshal(response, &ticket)

	fmt.Println("err", err)
	fmt.Printf(" info=%s \n", ticket.Info)
}
