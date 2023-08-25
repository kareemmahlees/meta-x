package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)


func ListDatabases(db *sqlx.DB)[]string {

	var dbs []string
	rows,_ := db.Queryx("SHOW DATABASES")
	for rows.Next(){
		db := struct{Database string `db:"Database"`}{}
		err := rows.StructScan(&db)
        if err != nil {
			log.Error(err)
        } 
		dbs = append(dbs, db.Database)
	}	
	return dbs
}

func CreateDatabase(db *sqlx.DB,dbName string)(int64,error){
	res,err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s",dbName))
	if err != nil {
		return 0,err
	}
	num,_ := res.RowsAffected()	
	return num,nil
}