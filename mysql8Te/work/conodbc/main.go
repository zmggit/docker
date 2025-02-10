package main

import (
	"database/sql"
	"fmt"
	_ "github.com/alexbrainman/odbc"
)

//安装对应的包
//brew install unixodbc
//arch=arm64 brew install unixodbc
//查看当前架构
//brew upgrade composer
//brew upgrade freetds

//odbcinst -j 查看配置文件路径
//file /opt/homebrew/opt/unixodbc/lib/libodbc.dylib
// file /usr/local/lib/libodbc.dylib
// /opt/homebrew/bin/odbcinst -j

///*
//#cgo CFLAGS: -I/opt/homebrew/opt/unixodbc/include
//#cgo LDFLAGS: -L/opt/homebrew/opt/unixodbc/lib -lodbc
//*/
//import "C"

func main() {

	// 数据库连接字符串
	//dsn := "Driver=myodbc;Server=127.0.0.1;Port=3306;Database=test;UID=zmg;PWD=123456;Charset=utf8"
	dsn := "DSN=myodbc;UID=zmg;PWD=123456"
	// 打开数据库连接
	db, err := sql.Open("odbc", dsn)
	if err != nil {
		fmt.Printf("连接数据库失败：%v\n", err)
		return
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		fmt.Printf("Ping 数据库失败：%v\n", err)
		return
	}

	fmt.Println("连接成功！")

	// 执行查询
	rows, err := db.Query("SELECT id, name FROM user")
	if err != nil {
		fmt.Printf("查询失败：%v\n", err)
		return
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			fmt.Printf("获取数据失败：%v\n", err)
			return
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
}
