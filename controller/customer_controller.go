package controller

import (
	"net/http"

	"sample.com/book/util"
)

func end(rw http.ResponseWriter, r *http.Request) {
	util.EndNormal(rw, http.StatusServiceUnavailable, nil)
}

func GetAllCustomers(rw http.ResponseWriter, r *http.Request) {
	end(rw, r)
}

func GetCustomer(rw http.ResponseWriter, r *http.Request) {
	end(rw, r)
}

func CreateCustomer(rw http.ResponseWriter, r *http.Request) {
	end(rw, r)
}

func ModifyCustomer(rw http.ResponseWriter, r *http.Request) {
	end(rw, r)
}

func RemoveCustomer(rw http.ResponseWriter, r *http.Request) {
	end(rw, r)
}
