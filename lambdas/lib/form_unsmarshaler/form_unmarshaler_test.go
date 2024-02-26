package form_unsmarshaler

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	StringField string `form:"stringField"`
	IntField    int    `form:"intField"`
	BoolField   bool   `form:"boolField"`
}

func TestUnMarshalBody(t *testing.T) {
	tests := []struct {
		name    string
		request events.APIGatewayProxyRequest
		target  interface{}
		wantErr bool
	}{
		{
			name: "Successful Unmarshaling",
			request: events.APIGatewayProxyRequest{
				Body: "stringField=Hello&intField=123&boolField=true",
			},
			target:  &TestStruct{},
			wantErr: false,
		},
		{
			name: "Empty Request Body",
			request: events.APIGatewayProxyRequest{
				Body: "",
			},
			target:  &TestStruct{},
			wantErr: true,
		},
		{
			name: "Invalid Form Data",
			request: events.APIGatewayProxyRequest{
				Body: "invalid",
			},
			target:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnMarshalBody(tt.request, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &TestStruct{}, tt.target)
			}
		})
	}
}

func TestUnmarshalForm(t *testing.T) {
	tests := []struct {
		name    string
		src     url.Values
		dst     interface{}
		wantErr bool
	}{
		{
			name: "Successful Unmarshaling",
			src: url.Values{
				"stringField": []string{"Hello"},
				"intField":    []string{"123"},
				"boolField":   []string{"true"},
			},
			dst:     &TestStruct{},
			wantErr: false,
		},
		{
			name:    "Nil Source",
			src:     nil,
			dst:     &TestStruct{},
			wantErr: true,
		},
		{
			name: "Destination Not a Pointer",
			src: url.Values{
				"stringField": []string{"Hello"},
			},
			dst:     TestStruct{},
			wantErr: true,
		},
		{
			name:    "Destination is Nil Pointer",
			src:     url.Values{},
			dst:     (*TestStruct)(nil),
			wantErr: true,
		},
		{
			name: "Destination Not a Struct",
			src: url.Values{
				"stringField": []string{"Hello"},
			},
			dst:     new(int),
			wantErr: true,
		},
		{
			name: "Unsupported Field Type",
			src: url.Values{
				"unsupportedField": []string{"unsupported"},
			},
			dst: &struct {
				UnsupportedField []string `form:"unsupportedField"`
			}{},
			wantErr: true,
		},
		{
			name: "Invalid Integer Value",
			src: url.Values{
				"intField": []string{"invalid"},
			},
			dst:     &TestStruct{},
			wantErr: true,
		},
		{
			name: "Invalid Boolean Value",
			src: url.Values{
				"boolField": []string{"invalid"},
			},
			dst:     &TestStruct{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalForm(tt.src, tt.dst)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if reflect.ValueOf(tt.dst).Kind() == reflect.Ptr {
					assert.Equal(t, reflect.Struct, reflect.ValueOf(tt.dst).Elem().Kind())
				}
			}
		})
	}
}
