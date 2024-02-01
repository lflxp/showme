package sqlite

// 设置数据库名称
func SetDbName(name string) {
	dbname = name
}

// 设置数据库驱动类型
// eg: sqlite3 | mysql | postgres | mssql | oracle
func SetDriverName(name string) {
	driverName = name
}

// 设置数据库连接地址
func SetDataResourceName(url string) {
	dataResourceName = url
}
