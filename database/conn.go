package database

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToMySQL(dc *DatabaseConfig) *sql.DB {
	return ConnectToMySQLWithDb(dc, "")
}

func ConnectToMySQLWithDb(dc *DatabaseConfig, dbName string) *sql.DB {
	// 拼接数据库连接
	var connectLinkBuilder strings.Builder
	connectLinkBuilder.Grow(10)
	connectLinkBuilder.WriteString(dc.Username)
	connectLinkBuilder.WriteString(":")
	connectLinkBuilder.WriteString(dc.Password)
	connectLinkBuilder.WriteString("@tcp(")
	connectLinkBuilder.WriteString(dc.Host)
	connectLinkBuilder.WriteString(":")
	connectLinkBuilder.WriteString(dc.Port)
	connectLinkBuilder.WriteString(")/")
	connectLinkBuilder.WriteString(dbName)
	if argCount := len(dc.Args); argCount > 0 {
		var argsBuilder strings.Builder
		argsBuilder.Grow(argCount * 4)
		argsBuilder.WriteString("?")
		for k, v := range dc.Args {
			argCount--
			argsBuilder.WriteString(k)
			argsBuilder.WriteString("=")
			argsBuilder.WriteString(v)
			if argCount > 0 {
				argsBuilder.WriteString("&")
			}
		}
		connectLinkBuilder.WriteString(argsBuilder.String())
	}

	db, err := sql.Open(dc.Driver, connectLinkBuilder.String())
	if err != nil {
		logger.Fatal("Connect to ", dbName, " failed:", err)
		return nil
	}

	return db
}
