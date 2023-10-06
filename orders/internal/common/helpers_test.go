package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	randomString := GenerateRandomString()
	assert.NotEmpty(t, randomString)
	assert.Len(t, randomString, 20)
}

func TestSetPageSize(t *testing.T) {
	// Test with input greater than max
	result := SetPageSize(30, 10, 20)
	assert.Equal(t, 20, result)

	// Test with input less than max
	result = SetPageSize(15, 10, 20)
	assert.Equal(t, 15, result)

	// Test with input less than or equal to 0
	result = SetPageSize(0, 10, 20)
	assert.Equal(t, 10, result)

	// Test with input equal to max
	result = SetPageSize(20, 10, 20)
	assert.Equal(t, 20, result)
}

func TestSetPageToken(t *testing.T) {
	// Test with input greater than 0
	result := SetPageToken(5)
	assert.Equal(t, 5, result)

	// Test with input less than or equal to 0
	result = SetPageToken(0)
	assert.Equal(t, 0, result)

	// Test with negative input
	result = SetPageToken(-5)
	assert.Equal(t, 0, result)
}

func TestIsFieldOutputOnly(t *testing.T) {
	// Test with an output-only field
	result := IsFieldOutputOnly("order_id")
	assert.True(t, result)

	// Test with a non-output-only field
	result = IsFieldOutputOnly("product_name")
	assert.False(t, result)
}
