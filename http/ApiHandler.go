package http

import (
	"../processor"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func DemoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write([]byte("Gorilla\n" + vars["city"] + vars["budget"]))
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]
	days := vars["days"]
	budget := vars["budget"]

	i_days, e := strconv.Atoi(days)
	if e != nil {
		panic("天数只能是数字")
	}

	f_budget, fe := strconv.ParseFloat(budget, 64)

	if fe != nil {
		panic("预算必须是浮点数")
	}

	path := processor.GenTripPath(city, i_days, f_budget)
	bytes, ji := json.Marshal(path)

	if ji != nil {
		panic("序列化发生错误")
	}

	w.Write([]byte(bytes))
}

func Startup() {
	r := mux.NewRouter()
	r.HandleFunc("/path/{city}/{days}/{budget}", pathHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}
