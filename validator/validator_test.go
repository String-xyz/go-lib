package validator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// tester has all its field as required fields
type tester struct {
	Email string `json:"email" validate:"required,email"`        // email must be a valid email
	Name  string `json:"name" validate:"required,gt=2"`          // the name field should have at least 2 letters(this can be adjusted)
	Age   int    `json:"age" validate:"required,gte=18,numeric"` // user must be 18yr or older
}

func TestValid(t *testing.T) {
	v := New()
	js := `{"email":"marlon@string.xyz", "name":"marlon", "age": 18}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.NoError(t, err)
	err = v.Validate(ts)
	assert.NoError(t, err)
}

func TestInvalidEmail(t *testing.T) {
	v := New()
	js := `{"email":"marlon@string", "name":"marlon", "age": 18}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.NoError(t, err)
	err = v.Validate(ts)
	assert.Error(t, err)
	bt, _ := json.Marshal(ExtractErrorParams(err))
	t.Log(string(bt))
}

func TestInvalidAge(t *testing.T) {
	v := New()
	js := `{"email":"marlon@string.xyz", "name":"marlon", "age": 10}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.NoError(t, err)
	err = v.Validate(ts)
	assert.Error(t, err)
	bt, _ := json.Marshal(ExtractErrorParams(err))
	t.Log(string(bt))

}

func TestInvalidName(t *testing.T) {
	v := New()
	js := `{"email":"marlon@string.xyz", "name":"m", "age": 18}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.NoError(t, err)
	err = v.Validate(ts)
	assert.Error(t, err)
	bt, _ := json.Marshal(ExtractErrorParams(err))
	t.Log(string(bt))

}

func TestMissingEmail(t *testing.T) {
	v := New()
	js := `{"name":"marlon", "age": 18}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.NoError(t, err)
	err = v.Validate(ts)
	assert.Error(t, err)
	bt, _ := json.Marshal(ExtractErrorParams(err))
	t.Log(string(bt))
}

func TestMissingName(t *testing.T) {
	v := New()
	js := `{"email":"marlon@string.xyz", "age": 18}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.NoError(t, err)
	err = v.Validate(ts)
	assert.Error(t, err)
	bt, _ := json.Marshal(ExtractErrorParams(err))
	t.Log(string(bt))
}

func TestMissingAge(t *testing.T) {
	v := New()
	js := `{"email":"marlon@string.xyz", "name":"m"}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.NoError(t, err)
	err = v.Validate(ts)
	assert.Error(t, err)
	bt, _ := json.Marshal(ExtractErrorParams(err))
	t.Log(string(bt))
}

func TestInvalidNumeric(t *testing.T) {
	v := New()
	js := `{"email":"marlon@string.xyz", "name":"marlon", "age":"string"}`
	ts := tester{}
	err := json.Unmarshal([]byte(js), &ts)
	assert.Error(t, err)
	err = v.Validate(ts)
	assert.Error(t, err)
	bt, _ := json.Marshal(ExtractErrorParams(err))
	t.Log(string(bt))
}
