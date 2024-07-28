package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDatabase() {

	var err error

	// DBConnectionString := "user=postgres.hhqemokocrfgbnaierzf password=Jkhqdlkhqlkwjd235hlqdw9adf! host=aws-0-ap-southeast-1.pooler.supabase.com? port=6543 dbname=postgres"
	DBConnectionString := "postgresql://postgres.hhqemokocrfgbnaierzf:Jkhqdlkhqlkwjd235hlqdw9adf!@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?pgbouncer=true"

	DB, err = gorm.Open(postgres.Open(DBConnectionString), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Database connected!")

}

type BLT struct {
	ID               uint `gorm:"primaryKey"`
	Nama             string
	Pekerjaan        string
	Penghasilan      string
	KepemilikanRumah string
	Aset             string
	JumlahTanggungan string
}

func (BLT) TableName() string {
	return "BLT"
}
