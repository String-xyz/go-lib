package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Id             string `json:"id" db:"id"`
	SomeField      string `json:"somefield" db:"somefield"`
	SomeOtherField int    `json:"someotherfield" db:"someotherfield"`
}

type ContactToPlatform struct {
	ContactId  string `json:"contactId" db:"contact_id"`
	PlatformId string `json:"platformId" db:"platform_id"`
}

func TestSanitizeModelInput(t *testing.T) {
	m := User{Id: "usr_0923840923840923840923"}
	err := SanitizeIdInput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "0923840923840923840923", m.Id)
}

func TestSanitizeModelOutput(t *testing.T) {
	m := User{Id: "0923840923840923840923"}
	err := SanitizeIdOutput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "usr_0923840923840923840923", m.Id)
}

func TestSanitizeRelationalModelInput(t *testing.T) {
	m := ContactToPlatform{ContactId: "cont_123", PlatformId: "plat_456"}
	err := SanitizeIdInput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "123", m.ContactId)
	assert.Equal(t, "456", m.PlatformId)
}

func TestSanitizeRelationalModelOutput(t *testing.T) {
	m := ContactToPlatform{ContactId: "123", PlatformId: "456"}
	err := SanitizeIdOutput(&m)
	assert.NoError(t, err)
	assert.Equal(t, "cont_123", m.ContactId)
	assert.Equal(t, "plat_456", m.PlatformId)
}

func TestSanitizeInputInline(t *testing.T) {
	m := "usr_123"
	m2 := "plat_456"

	err := SanitizeIdInput(&struct{ UserId, PlatformId string }{m, m2}, &m, &m2)
	assert.NoError(t, err)
	assert.Equal(t, "123", m)
	assert.Equal(t, "456", m2)
}

func TestSanitizeOutputInline(t *testing.T) {
	m := "123"
	m2 := "456"

	err := SanitizeIdOutput(&struct{ UserId, PlatformId string }{m, m2}, &m, &m2)
	assert.NoError(t, err)
	assert.Equal(t, "usr_123", m)
	assert.Equal(t, "plat_456", m2)
}
