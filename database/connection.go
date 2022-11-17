package database

import (
	"librarymanagement/entities"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbInstance *gorm.DB
var connectionString="root:root@tcp(localhost:3306)/librarymanagement"
var err error

func Connect(){
	DbInstance,err =gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err!=nil{
		log.Fatal(err)
		panic("cannot connect to Db")
	}
	log.Println("conected to db")
}

func Migrate(){
	DbInstance.AutoMigrate(&entities.BookList{},&entities.BorrowRecord{},&entities.LendingRecord{},&entities.BorrowUpdate{})
	log.Println("migration completed")
}