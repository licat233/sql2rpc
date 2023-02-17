/*
 * @Author: licat
 * @Date: 2023-02-07 14:54:44
 * @LastEditors: licat
 * @LastEditTime: 2023-02-07 15:08:06
 * @Description: licat233@gmail.com
 */

package _api

import (
	"bytes"
	"fmt"
)

type ApiCollection []*Api

func (ac ApiCollection) Len() int {
	return len(ac)
}

func (ac ApiCollection) Less(i, j int) bool {
	return ac[i].Path < ac[j].Path
}

func (ac ApiCollection) Swap(i, j int) {
	ac[i], ac[j] = ac[j], ac[i]
}

func (ac ApiCollection) String() string {
	buf := new(bytes.Buffer)
	for _, a := range ac {
		buf.WriteString(fmt.Sprint(a))
	}
	return buf.String()
}
