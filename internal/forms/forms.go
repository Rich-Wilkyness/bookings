package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0 // if we have 0 errors this will return true (shorthand in Go)
}

// initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}), // need {} here becuase we are declaring it empty
	}
}

// server validation for required fields
func (f *Form) Required(fields ...string) { // the ...string is called a variatic function. This means it can take any number of arguments that are strings
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" { // TrimSpace removes extra space the user may have input into the field
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// checks for valid email
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
