package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func InitDb() (err error) {
	namingStrategy := schema.NamingStrategy{
		TablePrefix:   "tb_",
		SingularTable: true,
	}
	Db, err = gorm.Open(postgres.Open(os.Getenv("DBPATH")), &gorm.Config{AllowGlobalUpdate: false, NamingStrategy: &namingStrategy})
	return err
}
