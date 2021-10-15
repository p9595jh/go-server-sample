package model

import (
	"bytes"
	"reflect"
	"strings"

	"sample.com/book/util"
)

/*

each book is saved as a text file (.txt)
[todo: set its encoding to utf8]

shape of the book file is like below:
name
author
publisher
genre

the items are seperated with a new line (\n)
the filename is composed with zero-padded number (like 00023.txt) and this number means an id of the book

** every book name cannot be duplicated

*/

type Book struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Genre     string `json:"genre"`
}

// get all columns of Book structure using reflect
var BookColumns = func() []string {
	columns := make([]string, 0)
	book := Book{}
	e := reflect.ValueOf(&book).Elem()

	for i := 0; i < e.NumField(); i++ {
		columns = append(columns, e.Type().Field(i).Name)
	}
	return columns
}()

func (b *Book) String() string {
	var buf bytes.Buffer
	buf.WriteString(b.Name)
	buf.WriteString("\n")
	buf.WriteString(b.Author)
	buf.WriteString("\n")
	buf.WriteString(b.Publisher)
	buf.WriteString("\n")
	buf.WriteString(b.Genre)
	return buf.String()
}

func (b *Book) FromString(s, filename string) {
	b.Id = util.IdFromFilename(filename)

	//  parse the content of the txt file and init
	cut := strings.Split(strings.Trim(s, " \t\n"), "\n")
	b.Name = cut[0]
	b.Author = cut[1]
	b.Publisher = cut[2]
	b.Genre = cut[3]
}

func (b *Book) FromIdAndString(s string, id int) {
	b.Id = id

	cut := strings.Split(strings.Trim(s, " \t\n"), "\n")
	b.Name = cut[0]
	b.Author = cut[1]
	b.Publisher = cut[2]
	b.Genre = cut[3]
}

func (b *Book) FromMap(m map[string]interface{}) {
	b.Name = m["name"].(string)
	b.Author = m["author"].(string)
	b.Publisher = m["publisher"].(string)
	b.Genre = m["genre"].(string)
}
