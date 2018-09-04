package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DbPojo struct {
	Id          uint `gorm:"primary_key"`
	Name        string
	Address     string
	Longitude   float64
	Latitude    float64
	Score       float64
	category    string
	City        string
	Area        string
	FloorPrice  float64 `gorm:"column:floor_price"`
	MiddlePrice float64 `gorm:"column:middle_price"`
}

type dataSourceConfig struct {
	username    string
	password    string
	hostAndPort string
	database    string
}

func (config dataSourceConfig) getDns() string {
	return config.username + ":" + config.password + "@tcp(" + config.hostAndPort + ")/" + config.database + "?charset=utf8"
}

var config = dataSourceConfig{
	"root",
	"root",
	"localhost:3306",
	"fasttrip",
}

var db *gorm.DB

func init() {
	cdb, err := gorm.Open("mysql", config.getDns())
	if err != nil {
		panic("获取数据库连接发生错误" + err.Error())
	}

	db = cdb
}

func SelectRestaurantByCity(city string) []DbPojo {
	return selectByCity("restaurant", city)
}

func SelectHotelByCity(city string) []DbPojo {
	return selectByCity("hotel", city)
}

func SelectScenic(city string) []DbPojo {
	return selectByCity("scenic", city)
}

func selectByCity(tableName string, city string) []DbPojo {
	// db.AutoMigrate(&restaurant{})
	var pojos []DbPojo
	db.Table(tableName).Limit(20).Order("score desc").Find(&pojos, "city = ?", city)
	return pojos
}
