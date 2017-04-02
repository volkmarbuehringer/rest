package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"rest/db"
)

func formReader(r *http.Request, name string, defaulter int) int {
	zu := r.FormValue(name)
	if len(zu) > 0 {
		la, err := strconv.Atoi(zu)
		if err == nil {
			return la
		}
	}
	return defaulter
}

func formReaderS(r *http.Request, name string, defaulter string) string {
	zu := r.FormValue(name)
	zu = strings.Replace(zu, "\"", "", 2)
	if len(zu) > 0 {
		return zu
	}
	return defaulter
}

func leser(w http.ResponseWriter, r *http.Request) (todo map[string]interface{}, err error) {

	todo = make(map[string]interface{})

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return
	}
	if err = r.Body.Close(); err != nil {
		return
	}
	if err = json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err = json.NewEncoder(w).Encode(err); err != nil {
			return
		}
	}
	return
}

func sender(w http.ResponseWriter, todos interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}

}

func senderErr(w http.ResponseWriter, err error) {
	type JSONErr struct {
		Error string
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err1 := json.NewEncoder(w).Encode(JSONErr{err.Error()}); err1 != nil {
		panic(err1)
	}

}

func getByIDHandler3(w http.ResponseWriter, r *http.Request) {

	ga := db.Pager{}
	ga.Length = formReader(r, "length", 100)
	ga.Offset = formReader(r, "offset", 0)
	ga.Where = formReaderS(r, "where", "")

	vars := mux.Vars(r)
	fmt.Println(ga)

	tab := vars["tab"]
	if m, ok := db.SQLTable[tab]; ok {
		rows, err := m.SelectAll(tab, &ga)

		if err != nil {
			senderErr(w, err)
			return
		}
		//	ctx.SetContentType("application/json")
		//	ctx.Writef("[")
		if _, err := m.Fetcher2(w, rows); err != nil {
			senderErr(w, err)
			return
		}
	} else {
		senderErr(w, fmt.Errorf("Tabelle nicht gefunden %s ", tab))
		return
	}
}

func getByIDHandler1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tab := vars["tab"]
	id, _ := strconv.Atoi(vars["id"])

	if m, ok := db.SQLTable[tab]; ok {
		rows, err := m.SelectID(tab, id)

		if err != nil {
			senderErr(w, err)
			return
		}

		gesamt, err := m.Fetcher(rows)
		if err != nil {
			senderErr(w, err)
			return
		}
		if len(gesamt) == 1 {
			sender(w, gesamt[0])
		} else {
			sender(w, gesamt)
			//			senderErr(w, fmt.Errorf("Daten nicht gefunden %s %d", tab, id))
			return
		}

	} else {
		senderErr(w, fmt.Errorf("Tabelle nicht gefunden: %s", tab))
		return
	}

}

func poster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tab := vars["tab"]
	if m, ok := db.SQLTable[tab]; ok {

		json, err := leser(w, r)
		if err != nil {
			senderErr(w, err)
			return
		}
		rows, err := m.RowInsert(&json, tab)
		//	fmt.Println(tab, m, s, json, r)
		//		rows, err := db.DB.Query(s, *r...)

		if err != nil {
			senderErr(w, err)
			return
		}

		gesamt, err := m.Fetcher(rows)
		if err != nil {
			senderErr(w, err)
			return
		}
		if len(gesamt) == 1 {
			sender(w, gesamt[0])
		} else {
			senderErr(w, fmt.Errorf("Daten nicht gefunden %s %d", tab, 0))
			return
		}

	} else {
		senderErr(w, fmt.Errorf("Tabelle nicht gefunden: %s", tab))
		return
	}

}

func putter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tab := vars["tab"]
	id, _ := strconv.Atoi(vars["id"])
	//checkErr(err)
	if m, ok := db.SQLTable[tab]; ok {
		json, err := leser(w, r)
		if err != nil {
			senderErr(w, err)
			return
		}
		rows, err := m.RowUpdate(&json, id, tab)

		//fmt.Println(tab, m, s, json, r)
		//rows, err := db.DB.Query(s, *r...)

		if err != nil {
			senderErr(w, err)
			return
		}

		gesamt, err := m.Fetcher(rows)
		if err != nil {
			senderErr(w, err)
			return
		}
		if len(gesamt) == 1 {
			sender(w, gesamt[0])
		} else {
			senderErr(w, fmt.Errorf("Daten nicht gefunden %s %d", tab, id))
		}

	} else {
		senderErr(w, fmt.Errorf("Tabelle nicht gefunden: %s", tab))
		return
	}
}
