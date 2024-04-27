package secexport

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
)

func IsJSON(s *string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(*s), &js) == nil
}

func IsJSONArray(s *string) bool {
	var array []interface{}
	return json.Unmarshal([]byte(*s), &array) == nil
}

func GetSHA1(s *string) string {
	hasher := sha1.New()
	io.WriteString(hasher, *s)
	sha := hasher.Sum(nil)

	return fmt.Sprintf("%x", sha)
}
