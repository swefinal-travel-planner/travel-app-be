package redis

import "fmt"

func Concat(baseKey string, number int64) string {
	return fmt.Sprintf("%s:%d", baseKey, number)
}
