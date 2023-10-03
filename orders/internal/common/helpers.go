package common

import (
	"math/rand"
	"time"
)

func GenerateRandomString() string {
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	length := 20
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func SetPageSize(input, def, max int) int {
	if input <= 0 {
		input = def
	} else if input > max {
		input = max
	}
	return input
}
func SetPageToken(input int) int {
	if input <= 0 {
		input = 0
	}
	return input
}

func IsFieldOutputOnly(field string) bool {
	list := [...]string{
		"order_id",
		"product_id",
		"tracking_number",
		"customer_id",
		"order_details_id",
		"create_time",
		"update_time",
		"delete_time",
		"version",
	}

	for _, curr := range list {
		if curr == field {
			return true
		}
	}

	return false
}
