package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"net/http"
	"net/url"
)

type ROUTES int

const (
	PadSize   = 5
	UrlPrefix = "/api/books"
)

const DataPath = "./data/"

func PaddedFilename(i int) string {
	return fmt.Sprintf("%0*d.txt", PadSize, i)
}

func IdFromFilename(filename string) int {
	// parse the id from the filename
	// filename format example: 00023.txt
	id, _ := strconv.Atoi(filename[:strings.LastIndex(filename, ".")])
	return id
}

func IdFromUrl(url string) (int, error) {
	return strconv.Atoi(url[strings.LastIndex(url, "/")+1:])
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func Includes(array []string, value string) bool {
	for _, s := range array {
		if s == value {
			return true
		}
	}
	return false
}

func Trim(s string) string {
	return strings.Trim(s, " \t\n")
}

func FirstUpperElseLower(s string) string {
	var res []rune = make([]rune, 0)
	defer func() {
		recover()
	}()

	if 'a' <= s[0] && s[0] <= 'z' {
		res = append(res, rune(s[0]-32))
	}
	for i := 1; i < len(s); i++ {
		if 'A' <= s[i] && s[i] <= 'Z' {
			res = append(res, rune(s[i]+32))
		} else {
			res = append(res, rune(s[i]))
		}
	}

	return string(res)
}

type ErrorDetail struct {
	path    string
	method  string
	message error
}

func EndError(rw http.ResponseWriter, r *http.Request, err error) {
	detail := ErrorDetail{
		path:    r.URL.Path,
		method:  r.Method,
		message: err,
	}
	log.Fatalln(err)
	bs, _ := json.Marshal(detail)
	const statusCode = http.StatusInternalServerError
	rw.WriteHeader(statusCode)
	rw.Write(bs)
}

func EndErrorMessage(rw http.ResponseWriter, r *http.Request, statusCode int, message string) {
	rw.WriteHeader(statusCode)
	rw.Write([]byte(message))
}

func EndNormal(rw http.ResponseWriter, statusCode int, data interface{}) {
	bs, _ := json.Marshal(data)
	rw.WriteHeader(statusCode)
	rw.Write(bs)
}

func ParseBody(r *http.Request) (map[string]interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	s := string(b)
	data := map[string]interface{}{}

	if strings.HasPrefix(s, "----------------------------") {
		lines := strings.Split(s, "\n")
		for i := 0; i < len(lines); i += 4 {
			q := strings.Index(lines[i+1], `"`) + 1
			key := Trim(lines[i+1][q:len(lines[i+1])])
			value := Trim(lines[i+3])
			data[key] = value
		}
		return data, nil
	} else {
		if err := json.Unmarshal(b, &data); err == nil {
			return data, nil
		} else {
			m, _ := url.ParseQuery(s)
			for k, v := range m {
				if i, err := strconv.Atoi(v[0]); err == nil {
					data[k] = i
				} else {
					data[k] = v[0]
				}
			}
			return data, nil
		}
	}
}
