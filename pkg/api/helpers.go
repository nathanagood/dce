package api

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"reflect"

	"github.com/Optum/dce/pkg/errors"
	"github.com/gorilla/schema"
)

// BuildNextURL merges the next parameters of pagination into the request parameters and returns an API URL.
func BuildNextURL(u url.URL, i interface{}) (url.URL, error) {
	req := url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}

	values := url.Values{}
	err := schema.NewEncoder().Encode(i, values)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return url.URL{}, errors.NewInternalServer("unable to encode query", err)
	}

	req.RawQuery = values.Encode()
	return req, nil
}

// WriteContext prints the context out
func WriteContext(ctx context.Context, w io.StringWriter) error {
	if _, err := w.WriteString("{"); err != nil {
		return err
	}
	rv := reflect.ValueOf(ctx)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Type().Field(i)

			if f.Name == "key" {
				if _, err := w.WriteString(fmt.Sprintf("\"key\": %q,", rv.Field(i))); err != nil {
					return err
				}
			}
			if f.Name == "Context" {

				// this is just a repetition of the above, so you can make a recursive
				// function from it, or for loop, that stops when there are no more
				// contexts to be inspected.

				rv := rv.Field(i)
				for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
					rv = rv.Elem()
				}

				if rv.Kind() == reflect.Struct {
					for i := 0; i < rv.NumField(); i++ {
						f := rv.Type().Field(i)

						if f.Name == "key" {
							if _, err := w.WriteString(fmt.Sprintf("\"key\": %q,", rv.Field(i))); err != nil {
								return err
							}
						}
						// ...
					}
				}
			}
		}
	}
	_, err := w.WriteString("}")
	return err
}
