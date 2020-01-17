package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	//空白导入
	_ "github.com/denisenkom/go-mssqldb"
)

type msdb struct {
	server   string
	port     int
	user     string
	password string
	database string
	conn     *sql.DB
}

//NewDb 数据库配置
func NewDb() msdb {
	return msdb{
		server:   "127.0.0.1",
		port:     1433,
		user:     "sa",
		password: "sa",
		database: "LampWorld",
		conn:     nil,
	}
}

//NewdbConn 新建连接
func (d *msdb) NewConn() error {
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;;encrypt=disable", d.server, d.port, d.database, d.user, d.password)
	db, err := sql.Open("mssql", connString)
	if err != nil {
		return errors.New("连接数据库失败")
	}
	d.conn = db
	return nil
}

func (d msdb) Close() {
	d.conn.Close()
}

//Query mssql 语句查询数据库
func (d msdb) Query(str string) (res []map[string]string, err error) {
	var rows *sql.Rows
	rows, err = d.conn.Query(str)
	if err != nil {
		err = errors.New("查询数据库失败")
		return
	}
	defer rows.Close()
	res = rowToString(rows)
	return
}

//Query mssql 语句修改数据库
func (d msdb) Exec(str string) (num int64, err error) {
	var res sql.Result
	res, err = d.conn.Exec(str)
	if err != nil {
		err = errors.New("查询数据库失败")
		return
	}
	num, err = res.RowsAffected()
	if err != nil {
		err = errors.New("查询数据库失败")
		return
	}
	return
}

type users struct {
	Account string `json:"account"`
}

func rowToString(b interface{}) (res []map[string]string) {
	rows := b.(*sql.Rows)
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		record := make(map[string]string, len(columns))
		for i, val := range values {
			if val != nil {
				record[columns[i]] = convertRow(val)
			}
		}
		res = append(res, record)
	}
	return
}

func convertRow(row interface{}) string {
	switch b := row.(type) {
	case int:
		return strconv.FormatInt(int64(b), 10)
	case int32:
		return strconv.FormatInt(int64(b), 10)
	case int64:
		return strconv.FormatInt(b, 10)
	case float32:
		return strconv.FormatFloat(float64(b), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(b, 'f', -1, 64)
	case string:
		return b
	case []byte:
		return string(b)
	case time.Time:
		return b.Format("2006-01-02 15：04：05")
	case bool:
		return strconv.FormatBool(b)
	}
	return ""
}
