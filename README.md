<!--
 * @Author: licat
 * @Date: 2023-02-06 14:26:42
 * @LastEditors: licat
 * @LastEditTime: 2023-02-17 17:10:06
 * @Description: licat233@gmail.com
-->

# sql2rpc

Generates a protobuf file and a api file from your sql database

## Uses

Tips:  If your operating system is windows, the default encoding of windows command line is "GBK", you need to change it to "UTF-8", otherwise the generated file will be messed up.

### Install

```shell
go install github.com/licat233/sql2rpc@latest
```

### Use from the command line

```text
$ sql2rpc -h
Usage of sql2rpc:
  -api
        generate .api files
  -api_jwt string
        the api service jwt, example: Auth
  -api_middleware string
        the api service middleware,  split multiple value by ",", example: AuthMiddleware
  -api_multiple
        Generate multiple api files according to table
  -api_prefix string
        the api service route prefix, example: api
  -api_style string
        the struct json naming format: sql_rpc | SqlRpc | sqlRpc  (default "sqlRpc")
  -db_host string
        the database host (default "localhost")
  -db_password string
        the database password
  -db_port int
        the database port (default 3306)
  -db_schema string
        the database schema (required)
  -db_table string
        the database table, split multiple tables with "," (default "*")
  -db_type string
        the database type (default "mysql")
  -db_user string
        the database user (default "root")
  -dir string
        directory of generated files
  -filename string
        the generated file name, defaults to the service name
  -ignore_column string
        the column to ignore, split multiple value by ","
  -ignore_table string
        the table to ignore, split multiple value by ","
  -init
        Create default config file，priority is given to the data in this file
  -model
        generate extend model .go files
  -pb
        generate .proto files
  -pb_gopackage string
        the protocol buffer go_package, defaults to the service name
  -pb_multiple
        the generated in multiple rpc service mode
  -pb_package string
        the protocol buffer package, defaults to the service name
  -service_name string
        the service name, defaults to the database schema
  -up
        Upgrade sql2rpc to latest version
  -upgrade
        Upgrade sql2rpc to latest version
  -v    Current version
  -version
        Current version
```

### By cmd, generate the xxx.proto file and update it

```shell
sql2rpc -pb -db_schema="admin" -db_table="*" -service_name="Admin" -filename="admin.proto" -pb_package="admin_proto" -pb_gopackage="./admin_pb" -pb_multiple=false
```

### By cmd, generate the xxx.api file and update it

```shell
sql2rpc -api -db_schema="admin" -db_table="*" -service_name="admin-api" -filename="admin.api" -api_jwt="Auth" -api_middleware="AuthMiddleware" -api_prefix="api" -api_multiple=true
```

### By cmd, generate the xxxModel_extend.go file and update it

```shell
sql2rpc -model -db_schema="admin"
```

### Generate a configuration file

Config.yaml will be created in the current directory

```shell
sql2rpc -init
```

A sql2rpcConfig.yaml configuration file will be created. [Sample file](./exemples/sql2rpcConfig.yaml)

### By configuration file, generate the xxx.proto file and update it

Please ensure that sql2rpcConfig.yaml already exists in the current directory

```shell
sql2rpc
```

### Use [goZero](https://github.com/zeromicro/go-zero)'s goctl tool to generate service code

Generate api service code

```shell
goctl api go -api admin.api -dir ./api -style goZero
```

Generate rpc service code

```shell
goctl rpc protoc "admin.proto" --go_out="./rpc" --go-grpc_out="./rpc" --zrpc_out="./rpc"
```

Generate model service code

```shell
goctl model mysql ddl --src "admin.sql" -dir . --style goZero
```

### Upgrade sql2rpc to latest version

```shell
sql2rpc -upgrade
```

### Configuration description

If there is a yaml configuration file, the configuration data in the file will be used first

### Update file content description

The content in this block will not be updated

```protobuf
// The content in this block will not be updated
// 此区块内的内容不会被更新
//[custom messages start]

message Person {
  int64 id = 1;
  string name = 2;
  bool is_student = 3;
}

//[custom messages end]
```

### Exemples

[see](./exemples/)

### Thanks

+ [goZero](https://github.com/zeromicro/go-zero)
+ [Mikaelemmmm](https://github.com/Mikaelemmmm/sql2pb)
