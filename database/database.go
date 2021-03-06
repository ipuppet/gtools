package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ipuppet/gtools/cache"
	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/utils"
)

var (
	logger  *log.Logger
	dbCache *cache.Cache
)

func init() {
	logger = utils.Logger("database")
	dbCache = cache.NewCache()
}

func ConnectToMySQL(database string) *sql.DB {
	var dc map[string]string
	config.GetConfig("database.json", &dc)

	// 拼接数据库连接
	var connectLinkBuffer bytes.Buffer
	connectLinkBuffer.WriteString(dc["username"])
	connectLinkBuffer.WriteString(":")
	connectLinkBuffer.WriteString(dc["password"])
	connectLinkBuffer.WriteString("@tcp(")
	connectLinkBuffer.WriteString(dc["host"])
	connectLinkBuffer.WriteString(":")
	connectLinkBuffer.WriteString(dc["port"])
	connectLinkBuffer.WriteString(")/")
	connectLinkBuffer.WriteString(database)
	connectLinkBuffer.WriteString("?charset=")
	connectLinkBuffer.WriteString(dc["charset"])
	connectLinkBuffer.WriteString("&parseTime=True")

	db, err := sql.Open("mysql", connectLinkBuffer.String())
	if err != nil {
		logger.Fatal("Connect to ", database, " failed:", err)
		return nil
	}

	return db
}

func CleanCache() {
	dbCache.Clean()
}

func SQLQueryRetrieveMap(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	return _SQLQueryRetrieveMap(db, query, true, args...)
}

func SQLQueryRetrieveMapNoCache(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	return _SQLQueryRetrieveMap(db, query, false, args...)
}

func _SQLQueryRetrieveMap(db *sql.DB, query string, withCache bool, args ...interface{}) ([]map[string]interface{}, error) {
	// key 与 args 相关
	argsJson, _ := json.Marshal(args)
	keyStr := query + string(argsJson[:])
	cacheKey := utils.MD5(keyStr)

	// 读取缓存
	if withCache {
		cacheData := dbCache.Get(cacheKey)
		if cacheData != nil {
			return cacheData["results"].([]map[string]interface{}), nil
		}
	}

	// 准备查询语句
	stmt, err := db.Prepare(query)
	if err != nil {
		stmt.Close()
		return nil, err
	}
	defer stmt.Close()

	// 查询
	rows, err := stmt.Query(args...)
	if err != nil {
		rows.Close()
		return nil, err
	}
	defer rows.Close()

	// 数据列
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 列的个数
	count := len(columns)

	// 返回值 Map 切片
	results := make([]map[string]interface{}, 0)
	// 一条数据的各列的值（需要指定长度为列的个数，以便获取地址）
	values := make([]interface{}, count)
	// 一条数据的各列的值的地址
	valPointers := make([]interface{}, count)
	for rows.Next() {
		// 获取各列的值的地址
		for i := 0; i < count; i++ {
			valPointers[i] = &values[i]
		}

		// 获取各列的值，放到对应的地址中
		rows.Scan(valPointers...)

		// 一条数据的 Map (列名和值的键值对)
		entry := make(map[string]interface{})

		// Map 赋值
		for i, col := range columns {
			var v interface{}

			// 值复制给 val (所以 Scan 时指定的地址可重复使用)
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				// 字符切片转为字符串
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}

		results = append(results, entry)
	}

	if withCache {
		// 设置缓存
		dbCache.Set(cacheKey, map[string]interface{}{"results": results})
	}

	return results, nil
}

func SQLQueryRetrieveStruct(db *sql.DB, rowStruct interface{}, query string, args ...interface{}) error {
	row := db.QueryRow(query, args...)

	// 确定 Scan 函数的输入类型
	reflectStruct := reflect.ValueOf(rowStruct).Elem()
	params := make([]interface{}, reflectStruct.NumField())
	// 按顺序遍历结构体的每个元素，取其指针值
	for i := 0; i < reflectStruct.NumField(); i++ {
		params[i] = reflectStruct.Field(i).Addr().Interface()
	}

	if err := row.Scan(params...); err != nil {
		return err
	}

	return nil
}

func MustExec(result sql.Result, err error) (sql.Result, error) {
	if err != nil {
		return result, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return result, err
	}
	if rowsAffected == 0 {
		return result, errors.New("rowsAffected is 0")
	}

	return result, nil
}
