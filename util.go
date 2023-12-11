package slick

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// ParseRequestBody parses the request body based on the content type
// and returns the user given type T.
// Currently JSON and from-urlencoded are supported.
func ParseRequestBody[T any](c *Context) (T, error) {
	var t T
	ctype := c.Request.Header.Get("Content-Type")
	switch {
	case ctype == "application/json":
		defer c.Request.Body.Close()
		err := json.NewDecoder(c.Request.Body).Decode(&t)
		return t, err
	case ctype == "application/x-www-form-urlencoded":
		t, err := parseMultipartFormData[T](c)
		return t, err
	default:
		return t, fmt.Errorf("cannot parse request with content type (%s)", ctype)
	}
}

func parseMultipartFormData[T any](c *Context) (T, error) {
	var (
		result T
		v      = reflect.ValueOf(&result).Elem() // Get the reflect.Value of result
		t      = v.Type()                        // Get the reflect.Type of result
	)
	if err := c.Request.ParseForm(); err != nil {
		return result, err
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.CanSet() {
			formKey := lowerCaseFirst(t.Field(i).Name)
			if values, ok := c.Request.Form[formKey]; ok && len(values) > 0 {
				fieldValue := values[0]
				switch field.Kind() {
				case reflect.String:
					field.SetString(fieldValue)
				case reflect.Int:
					intResult, _ := strconv.Atoi(fieldValue)
					field.SetInt(int64(intResult))
				case reflect.Struct:
					if field.Type().PkgPath() == "time" && field.Type().Name() == "Time" {
						// Assume the fieldValue is in a specific layout, e.g., "2006-01-02"
						timeResult, err := time.Parse("2006-01-02", fieldValue)
						if err != nil {
							// Handle error
						} else {
							field.Set(reflect.ValueOf(timeResult))
						}
					}
				}
			}
		}
	}
	return result, nil
}

func isMultipartFormData(header string) bool {
	return strings.HasPrefix(header, "multipart/form-data")
}

func uppercaseFirst(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

func lowerCaseFirst(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}
