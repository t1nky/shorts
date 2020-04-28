package helper

import (
	"errors"
	"reflect"
	"sort"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// ListOfErrors : Returns formatted list of errors
func ListOfErrors(validationObject interface{}, e error) []map[string]string {
	ve := e.(validator.ValidationErrors)
	InvalidFields := make([]map[string]string, 0)

	for _, e := range ve {
		errorsList := map[string]string{}

		field, isFound := reflect.TypeOf(validationObject).FieldByName(e.Field())
		if isFound {
			name := field.Tag.Get("json")
			errorsList[name] = e.Tag()
			InvalidFields = append(InvalidFields, errorsList)
		}
	}

	return InvalidFields
}

// ResponseOK : response without data or error fields
type ResponseOK struct {
	Result string `json:"result"`
}

// ResponseError : response with error
type ResponseError struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// ValidationError : response with error
type ValidationError struct {
	Result string              `json:"result"`
	Error  []map[string]string `json:"error"`
}

// ResponseData : response with data
type ResponseData struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}

// NewResponseOK : Returns default success response
func NewResponseOK() ResponseOK {
	return ResponseOK{Result: "ok"}
}

// NewResponseError : Returns response with error message
func NewResponseError(err error) ResponseError {
	return ResponseError{Result: "error", Error: err.Error()}
}

// NewValidationError : Returns response with error message
func NewValidationError(validationObject interface{}, err error) interface{} {
	errorMessage := err.Error()
	if errorMessage == "EOF" || errorMessage == "unexpected EOF" {
		return NewResponseError(errors.New("Invalid requrest body"))
	}

	switch err.(type) {
	default:
		return NewResponseError(err)
	case validator.ValidationErrors:
		return ValidationError{Result: "error", Error: ListOfErrors(validationObject, err)}
	}
}

// NewResponseOkWithData : Returns successfull response with data
func NewResponseOkWithData(data interface{}) ResponseData {
	return ResponseData{Result: "ok", Data: data}
}

// MakeShortlinkFromID returns id in Base36 format
func MakeShortlinkFromID(s uint64) string {
	return strconv.FormatUint(uint64(s), 36)
}

// NewAbsoluteLinksOnlyError returns error to indicate that full link is not absolute
func NewAbsoluteLinksOnlyError() error {
	return errors.New("Only absolule URLs are supported")
}

// NewPageNotFoundError returns error to indicate that route was not found
func NewPageNotFoundError() error {
	return errors.New("Page not found")
}

// A data structure to hold a key/value pair.
type pair struct {
	Key   string
	Value uint64
}

type pairList []pair

func (p pairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p pairList) Len() int      { return len(p) }
func (p pairList) Less(i, j int) bool {
	if p[i].Value == p[j].Value {
		return p[i].Key > p[j].Key
	}
	return p[i].Value > p[j].Value
}

func sortMapByValue(m map[string]uint64) pairList {
	p := make(pairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

// GetTopDomains : Returns map of top domains, amount of returned domains is determined by the topCount parameter
func GetTopDomains(m map[string]uint64, topCount int) (res map[string]uint64) {
	res = make(map[string]uint64)

	if topCount > len(m) {
		topCount = len(m)
	}
	sorted := sortMapByValue(m)[:topCount]
	for _, v := range sorted {
		res[v.Key] = v.Value
	}

	return
}
