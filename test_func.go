package gorm_helper

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getMockGormDB() (*gorm.DB, error) {
	//创建sqlmock
	var err error
	var mockDB *sql.DB
	mockDB, _, err = sqlmock.New()
	if err != nil {
		return nil, err
	}
	//结合gorm、sqlmock
	db, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      mockDB,
	}), &gorm.Config{})
	if nil != err {
		return nil, err
	}
	return db, nil
}

//getRawSql 获取gorm预执行的SQL——此处不实际执行
//warning 基于`gorm.io/gorm v1.23.6`封装 仅单元测试使用，避免升级gorm版本造成报错
func getRawSql(db *gorm.DB, cb func(tx *gorm.DB) *gorm.DB) string {
	beforeDryRun := db.Statement.DB.DryRun
	db.Statement.DB.DryRun = true

	tx := cb(db)

	explainSql := tx.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)

	db.Statement.DB.DryRun = beforeDryRun
	if !beforeDryRun {
		db.Statement.SQL.Reset()
		db.Statement.Vars = nil
	}
	return explainSql
}
