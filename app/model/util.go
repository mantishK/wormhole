package model

import (
	"math/rand"
	"time"

	"github.com/coopernurse/gorp"
)

func generateRandomToken(dbMap *gorp.DbMap) string {
	r := ""
	for {
		r = generateRandomString(80)
		user_token_test := UserToken{}
		user_token_test.Token = r
		err := user_token_test.GetUserIdFromToken(dbMap)
		if err != nil && user_token_test.UserId == 0 {
			break
		}
	}
	return r
}

func generateRandomString(length int) string {
	pool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rs := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	r := ""
	for i, _ := range b {
		b[i] = pool[rs.Intn(len(pool))]
	}
	r = string(b)
	return r
}
