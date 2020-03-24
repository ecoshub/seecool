package seecool

// Select("users_info", "c*first_name").Group("gender").Order("count").LimitOff("10", "0")
func Select(table string, columns ...string) *query {
	q := query{queryType: "SELECT", table: table, columns: columns}
	return &q
}

// SelectD("users_info", "c*first_name").Group("gender").Order("count").LimitOff("10", "0")
func SelectD(table string, columns ...string) *query {
	q := query{queryType: "SELECT DISTINCT", table: table, columns: columns}
	return &q
}

// Delete("users").Using("users_info").Cond("users_info.user_id",">","100")
func Delete(table string) *query {
	q := query{queryType: "DELETE", table: table}
	return &q
}

// Insert("env").Keys("key", "value").Values("key", "value")
func Insert(table string) *query {
	q := query{queryType: "INSERT", table: table}
	return &q
}

// Update("users").Keys("username","ip_address").Values("eco","192.168.1.108").Equal("user_id", "1000")
func Update(table string) *query {
	q := query{queryType: "UPDATE", table: table}
	return &q
}
