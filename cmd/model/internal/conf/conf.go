package conf

import "github.com/licat233/sql2rpc/config"

var (
	baseIgnoreTables  = []string{} //当前服务必须忽略的表
	baseIgnoreColumns = []string{} //每个结构体必须忽略的列

	FileContent string //文件内容

	MoreIgnoreTables  = []string{}             //当前服务可能忽略的表
	MoreIgnoreColumns = []string{"id", "uuid"} //某个结构可能忽略的列

	IgnoreTables  = []string{} //当前解析需要忽略的表
	IgnoreColumns = []string{} //每个结构体需要忽略的列

	QueryRow  = "m.conn.QueryRowCtx"
	QueryRows = "m.conn.QueryRowsCtx"
)

// InitConfig It initializes the configuration.
func InitConfig() {
	IgnoreTables = append(config.IgnoreTables, baseIgnoreTables...)
	IgnoreColumns = append(config.IgnoreColumns, baseIgnoreColumns...)
	ChangeQueryString(config.C.ModelCache.GetBool())
}

func ChangeQueryString(isCache bool) {
	if isCache {
		QueryRow = "m.QueryRowNoCacheCtx"
		QueryRows = "m.QueryRowsNoCacheCtx"
	} else {
		QueryRow = "m.conn.QueryRowCtx"
		QueryRows = "m.conn.QueryRowsCtx"
	}
}
