package model

import (
	"crypto/sha512"
	"encoding/base64"
	"strings"
	"time"

	"github.com/coopernurse/gorp"
)

type User struct {
	UserID   int       `db:"user_id" json:"user_id"`
	UserName string    `db:"user_name" json:"user_name"`
	Password string    `db:"password" json:"-"`
	Created  time.Time `db:"created" json:"created"`
	Modified time.Time `db:"modified" json:"modified"`
}

func (u *User) Get(dbMap *gorp.DbMap) error {
	err := dbMap.SelectOne(u, "SELECT * FROM users WHERE user_id = ?", u.UserID)
	if err != nil {
		return err
	}
	u.Password = ""
	return nil
}

func (u *User) GetUserFromUserName(dbMap *gorp.DbMap) error {
	err := dbMap.SelectOne(u, "SELECT * FROM users WHERE user_name = ?", u.UserName)
	if err != nil {
		return err
	}
	u.Password = ""
	return nil
}

func (u *User) Save(dbMap *gorp.DbMap) error {
	u.Created = time.Now()
	u.Modified = time.Now()
	err := dbMap.Insert(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UserNameExists(dbMap *gorp.DbMap) (bool, error) {
	user := []User{}
	_, err := dbMap.Select(&user, "SELECT * FROM users WHERE user_name = ?", u.UserName)
	if err != nil {
		return false, err
	}
	if len(user) == 0 {
		return false, nil
	}
	u.Password = ""
	return true, nil
}

func (u *User) IsValidUser(dbMap *gorp.DbMap) (bool, error) {
	err := dbMap.SelectOne(u, "SELECT * FROM users WHERE user_name = ? AND password=?", u.UserName, u.Password)
	if err != nil {
		return false, err
	}
	if u.UserID == 0 {
		return false, nil
	}
	u.Password = ""
	return true, nil
}

func (u *User) IsValidPassword(dbMap *gorp.DbMap) (bool, error) {
	err := dbMap.SelectOne(u, "SELECT * FROM users WHERE user_id = ? AND password=?", u.UserID, u.Password)
	if err != nil {
		return false, err
	}
	if u.UserID == 0 {
		return false, nil
	}
	u.Password = ""
	return true, nil
}

func (u *User) UpdatePassword(dbMap *gorp.DbMap) error {
	u.Modified = time.Now()
	_, err := dbMap.Update(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetPasswordSalt(dbMap *gorp.DbMap) (string, error) {
	password, err := dbMap.SelectStr("SELECT password FROM users WHERE user_id = ?", u.UserID)
	if err != nil {
		return "", nil
	}
	passwordSlice := strings.Split(password, ".")
	salt := passwordSlice[len(passwordSlice)-1]
	return salt, nil
}

func (u *User) HashPassword(salt string) {
	password := u.Password
	sha := sha512.New()
	if len(salt) == 0 {
		salt = generateRandomString(8)
	}
	for i := 0; i < 1; i++ {
		result := base64.StdEncoding.EncodeToString(sha.Sum([]byte(password + salt)))
		password = string(result)
	}
	u.Password = password + "." + salt
}
