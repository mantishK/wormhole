package model

import (
	"time"

	"github.com/coopernurse/gorp"
)

type UserToken struct {
	UserId   int       `db:"user_id"`
	Token    string    `db:"token"`
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}

func (ut *UserToken) GetUserIdFromToken(dbMap *gorp.DbMap) error {
	err := dbMap.SelectOne(ut, "SELECT * FROM user_token WHERE token = ?", ut.Token)
	if err != nil {
		return err
	}
	return nil
}

func (ut *UserToken) Add(dbMap *gorp.DbMap) error {
	ut.Created = time.Now()
	ut.Modified = time.Now()
	ut.Token = generateRandomToken(dbMap)
	err := dbMap.Insert(ut)
	if err != nil {
		return err
	}
	return nil
}

func (ut *UserToken) Delete(dbMap *gorp.DbMap) error {
	_, err := dbMap.Exec("DELETE from user_token where token=?", ut.Token)
	if err != nil {
		return err
	}
	return nil
}

func (ut *UserToken) Update(dbMap *gorp.DbMap) error {
	ut.Modified = time.Now()
	token := generateRandomToken(dbMap)
	_, err := dbMap.Exec("UPDATE user_token SET token = ?, modified = ? WHERE user_id= ? AND token = ?", token, ut.Modified, ut.UserId, ut.Token)
	if err != nil {
		return err
	}
	ut.Token = token
	return nil
}
