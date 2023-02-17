/*
 * @Author: licat
 * @Date: 2023-02-06 21:49:31
 * @LastEditors: licat
 * @LastEditTime: 2023-02-07 11:42:04
 * @Description: licat233@gmail.com
 */
package _rpc

import (
	"bytes"
	"fmt"
)

type RpcCollection []*Rpc

func (rc RpcCollection) Len() int {
	return len(rc)
}

func (rc RpcCollection) Less(i, j int) bool {
	return rc[i].Name < rc[j].Name
}

func (rc RpcCollection) Swap(i, j int) {
	rc[i], rc[j] = rc[j], rc[i]
}

func (rc RpcCollection) String() string {
	var buf = new(bytes.Buffer)
	for _, rpc := range rc {
		buf.WriteString(fmt.Sprint(rpc))
	}
	return buf.String()
}
