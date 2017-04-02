func getByIDHandler5(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tab := vars["tab"]
	id, _ := strconv.Atoi(vars["id"])

	/*
		tt := reflect.TypeOf(&rt).Elem()
		inter := make([]interface{}, tt.NumField())
		//	itemTyp := reflect.TypeOf(&r).Elem()
		//	item := reflect.New(itemTyp)

		for i := 0; i < tt.NumField(); i++ {

			gaga := reflect.ValueOf(&rt).Elem().Field(i)
			tag := string(tt.Field(i).Tag)
			//	name := gaga.Type().Name()
			fmt.Printf("%+v  %s\n", tt.Field(i), tag)
			switch gaga.Interface().(type) {
			case *int:
				lala := new(int)
				gaga.Set(reflect.ValueOf(lala))
				inter[i] = lala
				//item.Elem().FieldByName(name).Set(reflect.ValueOf(lala))
			case string:
				//lala := new(string)
				//gaga.Set(reflect.ValueOf(lala))
				if gaga.CanAddr() {
					inter[i] = gaga.Addr()
				}
				inter[i] = &rt.Fs

			case sql.NullFloat64:
				fmt.Println("hier aaaa", reflect.ValueOf(gaga))
				if gaga.CanAddr() {
					fmt.Printf("lulu %d %v %v\n", i, gaga.Kind(), gaga.Addr().Pointer())
					inter[i] = gaga.Addr()
					inter[i] = &rt.Gaga
				}

			case *float64:
				lala := new(float64)
				gaga.Set(reflect.ValueOf(lala))
				inter[i] = lala

			case *time.Time:
				lala := new(time.Time)
				gaga.Set(reflect.ValueOf(lala))
				inter[i] = lala
			default:
				fmt.Print("hier anders")
			}


		}
	*/

	if m, ok := db.SQLTable[tab]; ok {
		rows, err := m.SelectID(tab, id)

		if err != nil {
			senderErr(w, err)
			return
		}
		defer rows.Close()
		rt, inter := db.GetTable("webgaga")

		la, _ := rows.Columns()
		la1, _ := rows.ColumnTypes()
		fmt.Println("co", la, la1)
		for rows.Next() {

			if err = rows.Scan(inter...); err != nil {
				senderErr(w, err)
				return
			}

		}
		fmt.Println("hier ganz fertig", rt, inter)
		if err = rows.Err(); err != nil {
			senderErr(w, err)
			return
		}

		sender(w, rt)

	} else {
		senderErr(w, fmt.Errorf("Tabelle nicht gefunden: %s", tab))
		return
	}

}
