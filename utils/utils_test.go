package utils

import "testing"

func Test_JWT(t *testing.T) {
	print(JWTGenerate("2022", "pk"))
}
