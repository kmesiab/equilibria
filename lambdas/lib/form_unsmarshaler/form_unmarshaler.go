package form_unsmarshaler

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// UnMarshalBody unmarshal an url.Values struct into a struct annotated with `form` tags.
func UnMarshalBody(request events.APIGatewayProxyRequest, val interface{}) error {

	var err error
	var formValues url.Values

	if request.Body == "" {
		return fmt.Errorf("request body is empty")
	}

	formValues, err = url.ParseQuery(request.Body)
	if err != nil {
		return fmt.Errorf("error parsing request body: %s", err)
	}

	if err = UnmarshalForm(formValues, val); err != nil {
		return fmt.Errorf("error unmarshaling request body: %s", err)
	}

	return nil
}

func UnmarshalForm(src url.Values, dst interface{}) error {
	if src == nil {

		return fmt.Errorf("src cannot be nil")
	}

	dstVal := reflect.ValueOf(dst)

	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {

		return fmt.Errorf("dst must be a non-nil pointer")
	}

	dstVal = dstVal.Elem()
	if dstVal.Kind() != reflect.Struct {

		return fmt.Errorf("dst must be a pointer to a struct")
	}

	for i := 0; i < dstVal.NumField(); i++ {

		field := dstVal.Type().Field(i)
		formKey := field.Tag.Get("form")

		if formKey == "" {
			formKey = field.Name
		}

		if value, ok := src[formKey]; ok && len(value) > 0 {
			fieldVal := dstVal.Field(i)
			if fieldVal.CanSet() {

				switch fieldVal.Kind() {

				case reflect.String:
					fieldVal.SetString(value[0])

				case reflect.Int,
					reflect.Int8,
					reflect.Int16,
					reflect.Int32,
					reflect.Int64:

					intVal, err := strconv.ParseInt(value[0], 10, 64)
					if err != nil {

						return fmt.Errorf("invalid integer value for %s: %s",
							formKey, value[0])
					}
					fieldVal.SetInt(intVal)

				case reflect.Float32, reflect.Float64:
					floatVal, err := strconv.ParseFloat(value[0], 64)
					if err != nil {

						return fmt.Errorf("invalid float value for %s: %s",
							formKey, value[0])
					}

					fieldVal.SetFloat(floatVal)
				case reflect.Bool:
					boolVal, err := strconv.ParseBool(value[0])
					if err != nil {

						return fmt.Errorf("invalid boolean value for %s: %s",
							formKey, value[0])
					}
					fieldVal.SetBool(boolVal)
				default:

					return fmt.Errorf("unsupported field type: %s",
						fieldVal.Kind())
				}
			}
		}
	}

	return nil
}
