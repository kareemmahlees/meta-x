package internal

import (
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

func InitDBConn(provider, cfg string) (*sqlx.DB, error) {
	// cfg := mysql.Config{
	// 	User:   dbUsername,
	// 	Passwd: dbPassword,
	// 	DBName: dbName,
	// 	Net:    "tcp",

	// 	AllowNativePasswords: true,
	// }

	// // if dbPort is set it means we are in testing
	// if dbPort != "" {
	// 	cfg.Addr = fmt.Sprintf("127.0.0.1:%s", dbPort)
	// }

	db, err := sqlx.Connect(provider, cfg)

	if err != nil {
		return nil, err
	}
	return db, nil
}
