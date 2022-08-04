package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/canhlinh/cointracker/backend"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
)

var formDecoder = form.NewDecoder()

type Error struct {
	Message string `json:"message"`
}

type ResponseMeta struct {
	Total int64 `json:"total"`
}

type Response struct {
	Data interface{}   `json:"data"`
	Meta *ResponseMeta `json:"meta"`
}

func NewResponse(data interface{}) *Response {
	return &Response{Data: data}
}

func (r *Response) SetTotal(total int64) *Response {
	r.Meta = &ResponseMeta{Total: total}
	return r
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	App     *backend.App
}

func NewContext(c *backend.App, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		App:     c,
		Writer:  w,
		Request: r,
	}
}

func (c *Context) Jsonify(res interface{}, err error) {
	c.Writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		c.Writer.WriteHeader(400)
		json.NewEncoder(c.Writer).Encode(NewResponse(&Error{
			Message: err.Error(),
		}))
		return
	}

	json.NewEncoder(c.Writer).Encode(res)
}

func (c *Context) BindJSON(obj interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(obj)
}

func (c *Context) ParamInt(key string) int64 {
	rawValue := chi.URLParam(c.Request, key)
	value, _ := strconv.ParseInt(rawValue, 10, 64)
	return value
}

func (c *Context) Param(key string) string {
	return chi.URLParam(c.Request, key)
}

func (c *Context) ParamStr(key string) string {
	return chi.URLParam(c.Request, key)
}

func (c *Context) QueryInt(key string) int64 {
	rawValue := c.Request.URL.Query().Get(key)
	value, _ := strconv.ParseInt(rawValue, 10, 64)
	return value
}

func (c *Context) BindQuery(obj interface{}) error {
	return formDecoder.Decode(obj, c.Request.URL.Query())
}
