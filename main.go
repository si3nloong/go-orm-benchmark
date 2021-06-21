package main

import (
	"context"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/si3nloong/sqlike"
	"github.com/si3nloong/sqlike/actions"
	"github.com/si3nloong/sqlike/options"
	"github.com/si3nloong/sqlike/sql/expr"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
)

type User struct {
	ID        string    `gorm:"column:ID" xorm:"ID"`
	Name      string    `gorm:"column:Name" xorm:"Name"`
	Age       int       `gorm:"column:Age" xorm:"Age"`
	Status    string    `gorm:"column:Status" xorm:"Status"`
	CreatedAt time.Time `gorm:"column:CreatedAt" xorm:"CreatedAt"`
}

type Logger struct {
}

func (l Logger) Debug(stmt *sqlstmt.Statement) {
	// log.Printf("%v", stmt)
	log.Printf("%+v", stmt)
}

func main() {

	query := "INSERT INTO `users` (`ID`, `Name`, `Age`, `Status`, `CreatedAt`) VALUES " + strings.Repeat(",(?,?,?,?,?)", 25)[1:] + ";"
	log.Println(query)
	ctx := context.Background()
	client := sqlike.MustConnect(
		ctx,
		"mysql",
		options.Connect().
			SetUsername("root").
			SetPassword("abcd1234").
			SetHost("localhost").
			SetPort("3306"),
	)

	client.SetPrimaryKey("ID")
	client.SetLogger(&Logger{})
	table := client.Database("sqlike").Table("users")

	datas := []User{}
	result, err := table.Paginate(
		ctx,
		actions.Paginate().
			OrderBy(
				expr.Desc("CreatedAt"),
			),
		options.Paginate().
			SetDebug(true),
	)
	if err != nil {
		panic(err)
	}

	result.NextCursor(ctx, "123")

	result.All(&datas)
	log.Println(result)
	log.Println(err)
}
