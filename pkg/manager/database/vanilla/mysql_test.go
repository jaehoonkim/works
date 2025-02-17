package vanilla_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jaehoonkim/sentinel/pkg/manager/database/vanilla/excute"
	mysql "github.com/jaehoonkim/sentinel/pkg/manager/database/vanilla/excute/dialects/mysql"
	"github.com/jaehoonkim/sentinel/pkg/manager/database/vanilla/stmt"
)

func TestGetCluster(t *testing.T) {
	const (
		dialect    = "mysql"
		connstring = "sentinel:sentinel@tcp(127.0.0.1:3306)/sentinel?charset=utf8mb4&parseTime=True&loc=Local"
	)
	var db *sql.DB
	db, err := sql.Open(dialect, connstring)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlexcute := excute.GetSqlExcutor(mysql.Dialect())

	var (
		tableName                     = "cluster"
		columns                       = []string{"uuid"}
		cond      stmt.ConditionStmt  = stmt.IsNull("deleted")
		order     stmt.OrderStmt      = nil
		page      stmt.PaginationStmt = nil
	)

	err = sqlexcute.QueryRows(tableName, columns, cond, order, page)(
		context.TODO(), db)(
		func(scan excute.Scanner, _ int) error {
			var uuid string
			if err := scan.Scan(&uuid); err != nil {
				return err
			}

			fmt.Println(uuid)

			return nil
		})

	if err != nil {
		t.Fatal(err)
	}

}
