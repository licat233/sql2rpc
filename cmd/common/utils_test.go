/*
 * @Author: licat
 * @Date: 2023-02-17 23:07:59
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 23:16:39
 * @Description: licat233@gmail.com
 */
package common

import (
	"fmt"
	"testing"
)

func Test_TrimSpace(t *testing.T) {
	str := "   This is a line  \n with trailing spaces.  \t \n"
	fmt.Println(FormatContent(str, "  ")) // Output: "This is a line  "
}
