package utils

import (
	"math/rand"
	"strconv"
)

// var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// func RandSeq(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

func RandSeq(n int) string {
	num := 1000000 + rand.Intn(9000000)
	return strconv.Itoa(num)
}
