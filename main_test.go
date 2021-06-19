package main

import (
	"context"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"xorm.io/xorm"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	db "github.com/upper/db/v4"
	uppermy "github.com/upper/db/v4/adapter/mysql"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        string    `gorm:"column:ID" xorm:"ID"`
	Name      string    `gorm:"column:Name" xorm:"Name"`
	Age       int       `gorm:"column:Age" xorm:"Age"`
	Status    string    `gorm:"column:Status" xorm:"Status"`
	CreatedAt time.Time `gorm:"column:CreatedAt" xorm:"CreatedAt"`
}

var (
	ctx = context.Background()
	db1 *gorm.DB
	db2 *sqlike.Database
	db3 *xorm.Engine
	db4 db.Session
	db5 *sqlx.DB
)

func init() {
	var err error
	dsn := "root:abcd1234@tcp(127.0.0.1:3306)/sqlike?charset=utf8mb4&parseTime=True&loc=Local"
	db1, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		panic("failed to connect database")
	}

	ctx := context.Background()
	db2 = sqlike.MustConnect(
		ctx,
		"mysql",
		options.Connect().
			SetUsername("root").
			SetPassword("abcd1234").
			SetHost("localhost").
			SetPort("3306"),
	).Database("sqlike")

	db3, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		panic(err)
	}
	table := db2.Table("users")
	table.MustMigrate(ctx, User{})
	table.Truncate(ctx)

	var settings = uppermy.ConnectionURL{
		Database: "sqlike",
		Host:     "127.0.0.1",
		User:     "root",
		Password: "abcd1234",
	}

	db4, err = uppermy.Open(settings)
	if err != nil {
		panic(err)
	}

	db5, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func newUser() (u *User) {
	u = new(User)
	u.ID = ksuid.New().String()
	u.Name = "name"
	u.Status = "status"
	u.CreatedAt = time.Now()
	return
}

func BenchmarkTestUpperDBSingle_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := db4.Collection("users").Insert(
			newUser(),
		); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkTestGormSingle_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := db1.Create(newUser()).Error; err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkTestGormMultiple_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if err := db1.Create([]*User{
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
		}).Error; err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkTestSqlikeSingle_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := db2.Table("users").InsertOne(
			ctx,
			newUser(),
		); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkTestSqlikeMultiple_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := db2.Table("users").Insert(
			ctx,
			&[]*User{
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
			}); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkTestXormSingle_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := db3.Table("users").InsertOne(
			newUser(),
		); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkTestXormMultiple_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := db3.Table("users").InsertMulti(
			[]*User{
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
				newUser(), newUser(), newUser(), newUser(), newUser(),
			},
		); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkTestSqlxSingle_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := newUser()
		if _, err := db5.NamedExec(`
		INSERT INTO users 
		(ID, Name, Age, Status, CreatedAt) 
		VALUES 
		(:id, :name, :age, :status, :createdAt)`,
			map[string]interface{}{
				"id":        user.ID,
				"name":      user.Name,
				"age":       user.Age,
				"status":    user.Status,
				"createdAt": user.CreatedAt.Format("2006-01-02 15:04:05"),
			},
		); err != nil {
			log.Println(err)
			b.FailNow()
		}
	}
}

func BenchmarkTestSqlxMultiple_Insert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		users := []*User{
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
			newUser(), newUser(), newUser(), newUser(), newUser(),
		}
		datas := make([]map[string]interface{}, 0)
		for _, user := range users {
			datas = append(datas, map[string]interface{}{
				"id":        user.ID,
				"name":      user.Name,
				"age":       user.Age,
				"status":    user.Status,
				"createdAt": user.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		if _, err := db5.NamedExec(`
		INSERT INTO users 
		(ID, Name, Age, Status, CreatedAt) 
		VALUES 
		(:id, :name, :age, :status, :createdAt)`,
			datas,
		); err != nil {
			log.Println(err)
			b.FailNow()
		}
	}
}
