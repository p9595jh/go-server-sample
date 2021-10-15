package controller

import (
	"fmt"
	"net/http"

	"sample.com/book/model"
	"sample.com/book/service"
	"sample.com/book/util"
)

func GetAllBooks(rw http.ResponseWriter, r *http.Request) {
	/*
		request examples:
			/api/books
			/api/books?field=id
			/api/books?field=name&as=desc
	*/

	field := r.URL.Query().Get("field")
	as := r.URL.Query().Get("as")

	var align service.Align
	align.Field = util.FirstUpperElseLower(field)
	align.As = as

	books, err := service.FindAllBooks(&align)
	if err != nil {
		util.EndError(rw, r, err)
		return
	}

	util.EndNormal(rw, http.StatusOK, books)
}

func GetBook(rw http.ResponseWriter, r *http.Request) {
	id, err := util.IdFromUrl(r.URL.Path)
	if err != nil {
		util.EndError(rw, r, err)
		return
	}

	book, err := service.FindBook(id)
	if err != nil {
		util.EndError(rw, r, err)
		// util.EndErrorMessage(rw, r, http.StatusNotFound, fmt.Sprintf("book '%d' not found", id))
		return
	} else {
		util.EndNormal(rw, http.StatusOK, book)
	}
}

func CreateBook(rw http.ResponseWriter, r *http.Request) {
	m, err := util.ParseBody(r)
	if err != nil {
		util.EndError(rw, r, err)
		return
	}
	var book model.Book
	book.FromMap(m)
	if err := service.CreateBook(&book); err != nil {
		util.EndError(rw, r, err)
		return
	} else {
		util.EndNormal(rw, http.StatusCreated, book)
		return
	}
}

func ModifyBook(rw http.ResponseWriter, r *http.Request) {
	id, err := util.IdFromUrl(r.URL.Path)
	if err != nil {
		util.EndError(rw, r, err)
		return
	}

	book, err := service.FindBook(id)
	if err != nil {
		util.EndErrorMessage(rw, r, http.StatusNotFound, fmt.Sprintf("book '%d' not found", id))
		return
	}

	m, err := util.ParseBody(r)
	if err != nil {
		util.EndError(rw, r, err)
		return
	}

	for k, v := range m {
		switch k {
		case "name":
			book.Name = v.(string)
		case "author":
			book.Author = v.(string)
		case "publisher":
			book.Publisher = v.(string)
		case "genre":
			book.Genre = v.(string)
		}
	}

	service.ModifyBook(&book)
	util.EndNormal(rw, http.StatusOK, book)
}

func RemoveBook(rw http.ResponseWriter, r *http.Request) {
	id, err := util.IdFromUrl(r.URL.Path)
	if err != nil {
		util.EndError(rw, r, err)
		return
	}

	err = service.RemoveBook(id)
	if err != nil {
		util.EndError(rw, r, err)
		return
	}

	util.EndNormal(rw, http.StatusNoContent, nil)
}
