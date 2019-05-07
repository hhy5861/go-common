package mysql

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	Model struct {
		CreatedAt int64 `json:"createdAt" gorm:"column:created_at"`
		UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at"`
		Valid     int32 `json:"valid" gorm:"column:valid"`
	}

	MysqlConf struct {
		Dialect   string `yaml:"dialect"`
		Host      string `yaml:"host"`
		Port      int64  `yaml:"port"`
		DbName    string `yaml:"dbname"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		Charset   string `yaml:"charset"`
		ParseTime bool   `yaml:"parseTime"`
		MaxIdle   int    `yaml:"maxIdle"`
		MaxOpen   int    `yaml:"maxOpen"`
		Debug     bool   `yaml:"debug"`
	}
)

var (
	connMap sync.Map
	err     error
)

func NewMysql(conf *MysqlConf) *MysqlConf {
	return conf
}

const (
	connectionDefault = "default"
)

func NewDBClient(key ...string) *gorm.DB {
	conn, ok := connMap.Load(getConnKey(key))
	if !ok {
		log.Fatal("no mysql connection found")
		return nil
	}

	return conn.(*gorm.DB)
}

func (m *MysqlConf) Connection(key ...string) *gorm.DB {
	conn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.DbName,
		m.Charset,
		m.ParseTime)

	db, err := gorm.Open(m.Dialect, conn)
	if err != nil {
		log.Fatal(err)
	}

	db.Debug()
	db.DB().SetMaxIdleConns(m.MaxIdle)
	db.DB().SetMaxOpenConns(m.MaxOpen)
	db.DB().SetConnMaxLifetime(time.Second * 14400)

	db.SingularTable(true)
	db.LogMode(m.Debug)
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	connMap.LoadOrStore(getConnKey(key), db)
	return db
}

func (m *Model) BeforeSave() {
	m.UpdatedAt = tools.GetNowMillisecond()
}

func (m *Model) BeforeCreate() {
	m.CreatedAt = tools.GetNowMillisecond()
	m.UpdatedAt = tools.GetNowMillisecond()
}

func (m *Model) BeforeUpdate() {
	m.UpdatedAt = tools.GetNowMillisecond()
}

func (m *Model) BatchInsert(table string, fields, values []string) *gorm.DB {
	sql := fmt.Sprintf("INSERT INTO %s ( %s ) VALUES ", table, strings.Join(fields, ","))

	var i = 0
	var strSql = ""
	for _, v := range values {
		if i > 0 {
			strSql += " , "
		}

		strSql += fmt.Sprintf("( %s )", v)
		i++
	}

	sql = sql + strSql
	return NewDBClient().Exec(sql)
}

func getConnKey(key []string) string {
	if len(key) == 1 {
		return key[0]
	}

	return connectionDefault
}
