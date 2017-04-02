package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var bindSlice []string

type RowMaps []*map[string]interface{}

func init() {
	bindSlice = make([]string, 100)
	for i := 0; i < 100; i++ {
		bindSlice[i] = BindVar + strconv.Itoa(i)
	}
}

func (t *SQLStatement) RowUpdate(store *map[string]interface{}, id int, tab string) (rows *sql.Rows, err error) {
	m := make([]interface{}, 0, len(t.ColumnsUpdate))

	for i, c := range t.ColumnsUpdate {
		if strings.Contains(t.BindsUpdate[i], EqBindVar) {
			if mm, ok := (*store)[c]; ok {
				m = append(m, mm)
			} else {
				m = append(m, nil)
			}
		}
	}
	m = append(m, id)

	if t.sqlUpdate == nil {
		sqls := fmt.Sprintf(sqlUpdate, tab, strings.Join(t.BindsUpdate, ","), t.PK[0]+EqBindVar+strconv.Itoa(len(m)), t.ReturnClause())
		if t.sqlUpdate, err = DB.Prepare(sqls); err != nil {
			return
		}
	}
	rows, err = t.sqlUpdate.Query(m...)
	return

}

func (t *SQLStatement) RowInsert(store *map[string]interface{}, tab string) (rows *sql.Rows, err error) {
	m := make([]interface{}, 0, len(t.ColumnsInsert))

	for i, c := range t.ColumnsInsert {
		if t.BindsInsert[i][0:1] == BindVar {
			if mm, ok := (*store)[c]; ok {
				m = append(m, mm)
			} else {
				m = append(m, nil)
			}
		}
	}
	if t.sqlInsert == nil {
		sqls := fmt.Sprintf(sqlInsert, tab, strings.Join(t.ColumnsInsert, ","), strings.Join(t.BindsInsert, ","), t.ReturnClause())
		if t.sqlInsert, err = DB.Prepare(sqls); err != nil {
			return
		}
	}
	rows, err = t.sqlInsert.Query(m...)
	return
}
func (m *SQLStatement) SelectAll(tab string, ga *Pager) (rows *sql.Rows, err error) {

	if m.sqlAll == nil {
		m.sqlAll = make(map[string]*sql.Stmt)
	}
	if _, ok := m.sqlAll[ga.Where]; !ok {

		where := ""
		if len(ga.Where) > 0 {
			where = " and " + ga.Where
		}
		sqls := fmt.Sprintf(sqlLimit, strings.Join(m.ColumnsAllType, ","), tab, where, strings.Join(m.PK, ","))

		fmt.Println("prep", sqls)
		if m.sqlAll[ga.Where], err = DB.Prepare(sqls); err != nil {
			return
		}

	}
	rows, err = m.sqlAll[ga.Where].Query(ga.Length, ga.Offset)

	return
}

func (m *SQLStatement) SelectID(tab string, id int) (rows *sql.Rows, err error) {
	if m.sqlID == nil {
		sqls := fmt.Sprintf(sqlID, strings.Join(m.ColumnsAllType, ","), tab, m.PK[0], EqBindVar+"1", strings.Join(m.PK, ","))
		//	fmt.Println("prep", sqls)
		if m.sqlID, err = DB.Prepare(sqls); err != nil {
			return
		}
	}
	rows, err = m.sqlID.Query(id)

	return

}

func (m *SQLStatement) help(w io.Writer, zahler *int) func(**[]RowMaps) (err error) {

	return func(store **[]RowMaps) (err error) {
		var e []byte
		if *zahler > 0 {
			w.Write([]byte(","))
		}
		(*zahler)++
		gesamt, err := m.Packer(*store)
		if err != nil {
			return
		}
		if len(gesamt) > 0 {
			e, err = json.Marshal(gesamt[0])
			if err != nil {
				return
			}
			w.Write(e)

		} else {
			w.Write([]byte("[]"))
		}
		m.InitStore(store)
		return
	}

}
func (m *SQLStatement) Fetcher2(w io.Writer, rows *sql.Rows) (zahler int, err error) {
	defer rows.Close()
	zahler = 0
	//var merker int64

	store := m.InitStore(nil)
	for rows.Next() {
		r, mm := m.RowStore1()

		if err = rows.Scan(*r...); err != nil {
			return
		}
		helper := m.help(w, &zahler)
		store = m.Storer(store, mm, helper)

		/*
			if ggg, ok := (*mm[0])[m.PK[0]]; ok {
				g := reflect.ValueOf(ggg)
				gg := g.Elem().Interface().(int64)

				if merker > 0 && merker != gg {
					err = helper()
					if err != nil {
						return
					}
					store = m.InitStore()
				}
				merker = gg
				store = m.Storer(store, mm, helper)

			} else {
				err = fmt.Errorf("column nicht gefunden %s ", m.PK[0])
				fmt.Println("gaga", ggg, (*mm[0]))
				return
			}
		*/
	}
	if err = rows.Err(); err != nil {
		return
	}
	helper := m.help(w, &zahler)
	err = helper(&store)

	return
}

/*
func (m *SQLStatement) Checker(rows RowMaps, index int) RowMaps {

	neu := make([]*map[string]interface{}, 0, len(rows))
	pk := m.PK[index]
	for i := range rows {
		mapper := *(rows)[i]

		if h, ok := mapper[pk]; ok {
			g := reflect.ValueOf(h)
			if !g.Elem().IsNil() {
				neu = append(neu, rows[i])
			}
		} else {
			log.Fatal("pk  nicht gefunden", i, pk)
		}
	}

	return neu
}
*/

func (m *SQLStatement) Packer(gesamt *[]RowMaps) (gesamte RowMaps, err error) {
	if len(*gesamt) == 1 {
		gesamte = (*gesamt)[0]
		return
	}

	for i := len(m.Elemente) - 1; i > 0; i-- {

		v := m.Elemente[i]
		x := m.Parent[i]
		f := m.FK[i]
		fpk := m.PK[x]
		//		fmt.Println(x, v, len((*gesamt)[i]))
		ziel := (*gesamt)[x]

		mapper := make(map[int64]RowMaps)
		for _, v := range (*gesamt)[i] {
			if vgl, ok := (*v)[f]; ok {
				g1 := reflect.ValueOf(vgl)
				gg1 := g1.Elem().Interface().(int64)
				if vv, ok := mapper[gg1]; ok {
					//			fmt.Println(f, gg1)
					mapper[gg1] = append(vv, v)
				} else {
					m := make(RowMaps, 1)
					m[0] = v
					mapper[gg1] = m
				}
			}

		}
		//	fmt.Println(v, mapper)

		for j := range ziel {

			if x == 0 && j > 0 {
				break
			}
			if vgl, ok := (*ziel[j])[fpk]; ok {
				g1 := reflect.ValueOf(vgl)
				gg1 := g1.Elem().Interface().(int64)
				if vz, ok := mapper[gg1]; ok {
					//		fmt.Println("gefunden", gg1)
					(*ziel[j])[v] = vz
				} else {
					(*ziel[j])[v] = make(RowMaps, 0)
				}

			}

			//(*ziel[j])[v] = (*gesamt)[i]
		}

		//(*gesamt)[x] = (*gesamt)[x][:1]
	}
	gesamte = (*gesamt)[0]

	return
}

func (m *SQLStatement) Fetcher(rows *sql.Rows) (gesamte RowMaps, err error) {
	defer rows.Close()
	gesamt := m.InitStore(nil)
	for rows.Next() {
		var r *[]interface{}
		//		r, gesamt = m.RowStore(gesamt)
		r, mm := m.RowStore1()
		if err = rows.Scan(*r...); err != nil {
			return
		}
		gesamt = m.Storer(gesamt, mm, nil)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return m.Packer(gesamt)

}

func (t SQLStatement) RowStore1() (*[]interface{}, []*map[string]interface{}) {

	m := make([]interface{}, len(t.ColumnsAll))
	s := make([]*map[string]interface{}, len(t.Elemente))
	for i := range t.Elemente {
		p := make(map[string]interface{}, len(t.ColumnsAll))
		s[i] = &p
	}

	for i, c := range t.ColumnsAll {
		m[i] = new(interface{})
		g := t.Hierarchie[i]
		(*s[g])[c] = m[i]
	}
	return &m, s
}

//

func (t SQLStatement) Storer(storer *[]RowMaps, s []*map[string]interface{}, cb func(**[]RowMaps) error) *[]RowMaps {

	store := storer
	for i := range *store {
		mapper := s[i]
		pk := t.PK[i]
		if key, ok := (*mapper)[pk]; ok {
			g := reflect.ValueOf(key)

			if !g.Elem().IsNil() {
				gg := g.Elem().Interface().(int64)
				if lener := len((*store)[i]); lener > 0 {
					zack := (*store)[i][lener-1]
					if vgl, ok := (*zack)[pk]; ok {
						g1 := reflect.ValueOf(vgl)
						gg1 := g1.Elem().Interface().(int64)

						if gg != gg1 {
							if cb != nil && i == 0 {
								//fmt.Println("gag", gg, gg1, len((*store)[i]))
								cb(&store)

							}
							(*store)[i] = append((*store)[i], mapper)
							//fmt.Println("gag1", gg, gg1, len((*store)[i]))

						}

					} else {
						log.Fatal("key1 nicht gefunde")
					}
				} else {
					(*store)[i] = append((*store)[i], mapper)
				}
			}

		} else {
			log.Fatal("key nicht gefunde")
		}

	}
	return store
}

func (t SQLStatement) InitStore(store **[]RowMaps) *[]RowMaps {

	tt := make([]RowMaps, len(t.Elemente))
	for i := range t.Elemente {
		tt[i] = make(RowMaps, 0)

	}

	if store != nil {
		*store = &tt
	}

	return &tt
}

/*
func (t SQLStatement) RowStore(store *[]RowMaps) (*[]interface{}, *[]RowMaps) {

	return t.RowStore1()

	//return m, t.Storer(store, s)

}
*/
