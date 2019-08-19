package router

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	ussd "github.com/samora/ussd-go"
)

var startRegex *regexp.Regexp = regexp.MustCompile(`^\*\d+\*(\d)[\*|#]`)
var DB *sql.DB

type Router struct {
}

func (r *Router) Dispatch(c *ussd.Context) ussd.Response {
	var callback string
	if startRegex.Match([]byte(c.Request.Message)) {
		code := getCode(c)
		DB, err := sql.Open("mysql", "root:root@/ussd")
		defer DB.Close()
		if err != nil {
			return c.Err(err)
		}

		statement := `SELECT callback FROM callbacks where service_code =?`
		row := DB.QueryRow(statement)
		err = row.Scan(&callback)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Err(fmt.Errorf("No service %d", code))
			}
			return c.Err(err)
		}
		key := getSessionKey(c)
		_ = c.DataBag.Set(key, callback)
	} else {
		key := getSessionKey(c)
		url, err := c.DataBag.Get(key)
		if err != nil {
			return c.Err(err)
		}
		callback = url
	}
	body := prepareBody(c)
	resp, err := http.Post(callback, "application/json", bytes.NewBuffer(body)) //timeouts!!!
	if err != nil {
		c.Err(err)
	}
	return parseBody(resp, c)
}

func getCode(c *ussd.Context) int {

	matches := startRegex.FindStringSubmatch(c.Request.Message)
	Icode, err := strconv.Atoi(matches[1])
	if err != nil {
		Icode = 0
	}
	return Icode
}

func prepareBody(c *ussd.Context) []byte {
	var buf []byte
	buf, err := json.Marshal(c.Request)
	if err != nil {
		return []byte(`{"error":Service Unavailable}`)
	}
	return buf

}

func getSessionKey(c *ussd.Context) string {
	key := c.Request.Mobile + "_sessionkey"
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

func parseBody(r *http.Response, c *ussd.Context) ussd.Response {
	var m map[string]interface{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	if err != nil {
		return c.Err(err)
	}
	message, _ := m["Message"].(string)
	release, _ := m["Release"].(bool)
	if release {
		return c.Release(message)
	}
	return c.Render(message, "router", "Dispatch")

}
