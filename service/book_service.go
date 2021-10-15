package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"

	e "sample.com/book/error"
	"sample.com/book/model"
	"sample.com/book/util"
)

// closure function to represent true:1 and false:-1
// due to use in sorting of []Book
var boolToInteger = func() func(bool) int {
	innerMap := map[bool]int{
		true:  1,
		false: -1,
	}
	return func(b bool) int {
		return innerMap[b]
	}
}()

type Align struct {
	As    string // asc (default), desc
	Field string // field of Book
}

func FindAllBooks(align *Align) ([]model.Book, error) {
	books := make([]model.Book, 0)
	txts, err := ioutil.ReadDir(util.DataPath)
	if err != nil {
		return books, err
	}

	for _, txt := range txts {
		bs, _ := ioutil.ReadFile(util.DataPath + txt.Name())
		bookData := string(bs)
		book := model.Book{}
		book.FromString(bookData, txt.Name())
		books = append(books, book)
	}

	// if align == nil {
	// 	return books, nil
	// }
	if align.Field == "" {
		return books, nil
	}

	if !util.Includes(model.BookColumns, align.Field) {
		return books, nil
	}

	var c int
	switch align.As {
	case "desc":
		c = -1
	default:
		c = 1
	}

	switch align.Field {
	case "Id":
		sort.Slice(books, func(i, j int) bool {
			return books[i].Id*c > books[j].Id*c
		})
	default:
		sort.Slice(books, func(i, j int) bool {
			fi := reflect.Indirect(reflect.ValueOf(books[i])).FieldByName(align.Field).String()
			fj := reflect.Indirect(reflect.ValueOf(books[j])).FieldByName(align.Field).String()
			return boolToInteger(fi > fj)*c > 0
		})
	}

	return books, nil
}

type FileNotFoundError struct {
	Message string
}

func (e *FileNotFoundError) Error() string {
	return e.Message
}

func fileExists(filename string) bool {
	// defer func() {
	// 	recover()
	// }()
	// info, err := os.Stat(filename)
	// if os.IsNotExist(err) {
	// 	return false
	// }
	// return !info.IsDir()
	_, err := os.Open(filename)
	return err != nil
}

func FindBook(id int) (model.Book, error) {
	filename := util.PaddedFilename(id)
	path := util.DataPath + filename
	book := model.Book{}
	var e error
	// if !fileExists(path) {
	// 	e := FileNotFoundError{
	// 		Message: fmt.Sprintf("'%s' not exists", path),
	// 	}
	// 	return book, &e
	// }
	defer func() {
		recover()
		e = &FileNotFoundError{
			Message: path + " not found",
		}
	}()
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return book, err
	}
	book.FromIdAndString(string(bs), id)
	return book, e
}

// todo: apply goroutine
// implementation signleton to handle this function, not using file i/o
func FindBookByName(name string, except int) (model.Book, error) {
	txts, err := ioutil.ReadDir(util.DataPath)
	book := model.Book{}
	if err != nil {
		return book, err
	}
	for _, txt := range txts {
		bs, _ := ioutil.ReadFile(util.DataPath + txt.Name())
		s := util.Trim(string(bs))

		// get the name of this book
		// first line means name
		firstLine := s[:strings.Index(s, "\n")]
		if firstLine == name {
			id := util.IdFromFilename(txt.Name())
			if id != except {
				book.FromIdAndString(s, id)
				return book, nil
			}
		}
	}
	return book, &e.NameNotFoundError{Name: name}
}

// todo: apply utf8
func CreateBook(book *model.Book) error {
	// get file list
	txts, _ := ioutil.ReadDir(util.DataPath)

	var newId int
	switch len(txts) {
	case 0:
		// set first id to 1
		newId = 1
	default:
		// add 1 to the latest filename
		newId = util.IdFromFilename(txts[len(txts)-1].Name()) + 1
	}

	// create new filename
	filename := util.PaddedFilename(newId)

	f, err := os.Create(util.DataPath + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprint(f, book.String())
	book.Id = newId
	return nil
}

func ModifyBook(book *model.Book) error {
	filename := util.PaddedFilename(book.Id)
	f, err := os.Open(util.DataPath + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprint(f, book.String())
	return nil
}

func RemoveBook(id int) error {
	filename := util.PaddedFilename(id)
	return os.Remove(util.DataPath + filename)
}
