package dbi

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

// DBI  database interface
type DBI struct {
	db *gorm.DB
}

// Radiologdata DB model structure
type Radiologdata struct {
	gorm.Model
	Label        string
	Description  string
	Address      uint
	Timestamp    time.Time
	Lqi          int
	Rssi         int
	Uptime       int
	Tempcpu      int
	Vrefcpu      int
	Ntc0         int
	Ntc1         int
	Photores     int
	Pressure     int
	Temppressure int
}

// RadiologDevice Status in table
type RadiologDevice struct {
	gorm.Model
	Node      string `gorm:"type:varchar(100);unique_index"`
	Data      string
	Count     int
	Timestamp time.Time
}

// Temperature of gived address
func (dbp *DBI) Temperature(address uint) {
	var d []Radiologdata
	dbp.db.Where(&Radiologdata{Address: address}).Find(&d)

	for i := 0; i < 10; i++ {
		fmt.Println(d[i].Ntc0)
	}
}

// UpdateNode ...
func (dbp *DBI) UpdateNode(node string, data string) {
	var n RadiologDevice
	if err := dbp.db.Where("node = ?", node).First(&n).Error; gorm.IsRecordNotFoundError(err) {
		dbp.db.Create(&RadiologDevice{Node: node, Data: "", Count: 1, Timestamp: time.Now()})
		log.Info("New node found ", node)
	} else {
		dbp.db.Model(&n).Select("count", "timestamp").Updates(map[string]interface{}{"count": n.Count + 1, "timestamp": time.Now()})
		log.Info("Update Node: ", node)
	}
}

// StatusNode ...
func (dbp *DBI) StatusNode() []RadiologDevice {
	var nl []RadiologDevice
	now := time.Now()
	before := now.Add(-time.Duration(10) * time.Second)
	dbp.db.Where("timestamp BETWEEN ? AND ?", before, now).Find(&nl)

	return nl
}

// Init init db module interface
func (dbp *DBI) Init() error {
	config := dbConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	var err error
	dbp.db, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Unable to connect to DB: ", err)
		return err
	}
	//defer dbp.db.Close()
	dbp.db.LogMode(false)
	log.Info("Successfully connected!")

	// Migrate the schema
	dbp.db.AutoMigrate(&Radiologdata{})
	dbp.db.AutoMigrate(&RadiologDevice{})
	return nil
}
