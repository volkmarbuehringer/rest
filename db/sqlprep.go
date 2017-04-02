package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Pager struct {
	Length int
	Offset int
	ID     int
	Where  string
}

var SQLTable = map[string]*SQLStatement{}

type SQLStatement struct {
	Table          string
	ColumnsAll     []string
	ColumnsTrans   []string
	ColumnsAllType []string
	ColumnsInsert  []string
	BindsInsert    []string
	ColumnsUpdate  []string
	BindsUpdate    []string
	Elemente       []string
	Hierarchie     []int
	Parent         []int
	PK             []string
	FK             []string
	sqlID          *sql.Stmt
	sqlAll         map[string]*sql.Stmt
	sqlInsert      *sql.Stmt
	sqlUpdate      *sql.Stmt
}

func (row *SQLStatement) addColumns(rows *sql.Rows) {

	row.Elemente = make([]string, 1)
	row.FK = make([]string, 1)
	row.Parent = make([]int, 10)
	var pk string
	if len(row.PK) > 0 {
		pk = row.PK[0]
	}

	for inser, upder := 0, 0; rows.Next(); {

		var column, columnType string
		if err := rows.Scan(&column, &columnType); err != nil {
			log.Fatal(err)
		}
		if column[0:1] == "$" {
			column = column[1:]
			parent := -1
			//			parent, _ := strconv.Atoi(column[0:1])
			//			column = column[1:]
			pos := strings.Index(column, "$")
			if pos > 0 {

				row.Elemente = append(row.Elemente, column[:pos])

				column = column[pos+1:]
				pos := strings.Index(column, "$")
				if pos > 0 {
					pk = column[:pos]
					row.PK = append(row.PK, pk)
					column = column[pos+1:]
					pos := strings.Index(column, "$")
					if pos > 0 {
						row.FK = append(row.FK, column[:pos])
						column = column[pos+1:]

						for i, r := range row.PK {
							if r == column {
								parent = i
							}
						}

					} else {
						log.Fatal("pos3", column, row)
					}

				} else {
					log.Fatal("pos2")
				}

				row.Parent[len(row.Elemente)-1] = parent
			} else {
				log.Fatal("pos1")
			}

		} else {
			row.Hierarchie = append(row.Hierarchie, len(row.Elemente)-1)
			row.ColumnsAll = append(row.ColumnsAll, column)
			row.ColumnsAllType = append(row.ColumnsAllType, columnType)

			switch {
			case column == pk:
				if g := dbSequenzer(row.Table); len(g) > 0 {
					row.ColumnsInsert = append(row.ColumnsInsert, column)
					row.BindsInsert = append(row.BindsInsert, g)
				}
			case strings.Contains(column, "_cr_date"):
				row.ColumnsInsert = append(row.ColumnsInsert, column)
				row.BindsInsert = append(row.BindsInsert, dbTimestamp)
			case strings.Contains(column, "_upd_date"):
				row.ColumnsUpdate = append(row.ColumnsUpdate, column)
				row.BindsUpdate = append(row.BindsUpdate, column+"="+dbTimestamp)
			case strings.Contains(column, "_upd_uid"):
				row.ColumnsUpdate = append(row.ColumnsUpdate, column)
				row.BindsUpdate = append(row.BindsUpdate, column+"='webSrv'")
			case strings.Contains(column, "_cr_uid"):
				row.ColumnsInsert = append(row.ColumnsInsert, column)
				row.BindsInsert = append(row.BindsInsert, "'webSrv'")
			default:
				inser++
				upder++
				row.ColumnsInsert = append(row.ColumnsInsert, column)
				row.ColumnsUpdate = append(row.ColumnsUpdate, column)
				row.BindsInsert = append(row.BindsInsert, BindVar+strconv.Itoa(inser))
				row.BindsUpdate = append(row.BindsUpdate, column+EqBindVar+strconv.Itoa(upder))
			}

		}

	}
}

func initStore() {

	rows, err := DB.Query(sqlalltabs)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var row SQLStatement
		var pk string

		if err := rows.Scan(&row.Table, &pk); err != nil {
			log.Fatal(err)
		}
		if len(pk) > 0 {
			row.PK = append(row.PK, pk)
		}

		if rows, err := DB.Query(sqlallcols, row.Table); err != nil {
			log.Fatal(err)
		} else {

			row.addColumns(rows)
			fmt.Println("lese", row)
			SQLTable[row.Table] = &row
			fmt.Println(row)
		}

	}

}
