package enum

import (
    "fmt"
    "slices"
    "strings"
)

// This file has been created automatically by `go-enum-generate`
// DO NOT MODIFY NOR EDIT THIS FILE DIRECTLY.
// To modify this enum, edit the enums.json or enums.yaml definition file
// To know more about `go-enum-generate`, see go to `https://github.com/debarbarinantoine/go-enum-generate`
// Generated at: 2025-08-21 20:58:55

type Tag uint

const (
    query Tag = iota
    urlParam
    header
    cookie
    json
    form
    multipartForm
)

var tagKeys = make(map[Tag]struct{}, 7)
var tagValues = make(map[string]Tag, 7)
var tagKeysArray = make([]Tag, 7)
var tagValuesArray = make([]string, 7)

func init() {
    tagKeys[query] = struct{}{}
    tagKeysArray[0] = query
    tagValues["query"] = query
    tagValuesArray[0] = "query"

    tagKeys[urlParam] = struct{}{}
    tagKeysArray[1] = urlParam
    tagValues["param"] = urlParam
    tagValuesArray[1] = "param"

    tagKeys[header] = struct{}{}
    tagKeysArray[2] = header
    tagValues["header"] = header
    tagValuesArray[2] = "header"

    tagKeys[cookie] = struct{}{}
    tagKeysArray[3] = cookie
    tagValues["cookie"] = cookie
    tagValuesArray[3] = "cookie"

    tagKeys[json] = struct{}{}
    tagKeysArray[4] = json
    tagValues["json"] = json
    tagValuesArray[4] = "json"

    tagKeys[form] = struct{}{}
    tagKeysArray[5] = form
    tagValues["form"] = form
    tagValuesArray[5] = "form"

    tagKeys[multipartForm] = struct{}{}
    tagKeysArray[6] = multipartForm
    tagValues["multipart"] = multipartForm
    tagValuesArray[6] = "multipart"
}

func (e Tag) String() string {
    switch e {
        case query:
            return "query"
        case urlParam:
            return "param"
        case header:
            return "header"
        case cookie:
            return "cookie"
        case json:
            return "json"
        case form:
            return "form"
        case multipartForm:
            return "multipart"
        default:
            return fmt.Sprintf("Unknown Tag (%d)", e.Value())
    }
}

func (e *Tag) Parse(str string) error {

    str = strings.TrimSpace(str)

    if val, ok := tagValues[str]; ok {
        *e = val
        return nil
    }
    return fmt.Errorf("invalid Tag: %s", str)
}

func (e Tag) Value() uint {
    return uint(e)
}

func (e Tag) MarshalText() ([]byte, error) {
    return []byte(e.String()), nil
}

func (e *Tag) UnmarshalText(text []byte) error {
    return e.Parse(string(text))
}

func (e Tag) IsValid() bool {
    if _, ok := tagKeys[e]; !ok {
        return false
    }
    return true
}

type tags struct {
    Query Tag
    UrlParam Tag
    Header Tag
    Cookie Tag
    Json Tag
    Form Tag
    MultipartForm Tag
}

var Tags = tags{
    Query: query,
    UrlParam: urlParam,
    Header: header,
    Cookie: cookie,
    Json: json,
    Form: form,
    MultipartForm: multipartForm,
}

func (e tags) Values() []Tag {
    return slices.Clone(tagKeysArray)
}

func (e tags) Args() []string {
    return slices.Clone(tagValuesArray)
}

func (e tags) Description() string {
    var strBuilder strings.Builder
    strBuilder.WriteString("\tAvailable Tags:\n")
    for _, enumVal := range e.Values() {
        strBuilder.WriteString(fmt.Sprintf("=> %d -> %s\n", enumVal.Value(), enumVal.String()))
    }
    return strBuilder.String()
}

func (e tags) Cast(value uint) (Tag, error) {
    if _, ok := tagKeys[Tag(value)]; !ok {
        return 0, fmt.Errorf("invalid cast Tag: %d", value)
    }
    return Tag(value), nil
}
