package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"sample.com/book/controller"
)

var route_map = map[string]map[int]map[string]map[string]func(http.ResponseWriter, *http.Request){
	"api": {
		2: {
			"books": {
				"GET":  controller.GetAllBooks,
				"POST": controller.CreateBook,
			},
			"customers": {
				"GET":  controller.GetAllCustomers,
				"POST": controller.CreateCustomer,
			},
		},
		3: {
			"books": {
				"GET":    controller.GetBook,
				"PUT":    controller.ModifyBook,
				"DELETE": controller.RemoveBook,
			},
			"customers": {
				"GET":    controller.GetCustomer,
				"PUT":    controller.ModifyCustomer,
				"DELETE": controller.RemoveCustomer,
			},
		},
	},
}

func handler(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if recover() != nil {
			// log.Printf("invalid access")
			rw.WriteHeader(404)
			rw.Write([]byte("invalid access"))
		}
	}()

	log.Println(r.Method, strings.TrimRight(r.URL.Path, "/"))

	path := strings.Trim(r.URL.Path, "/")
	routes := strings.Split(path, "/")

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	route_map[routes[0]][len(routes)][routes[1]][r.Method](rw, r)
}

func main() {
	const PORT = 3000

	fmt.Println()
	log.Println("server started at", PORT)
	fmt.Println()

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), http.HandlerFunc(handler))
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
