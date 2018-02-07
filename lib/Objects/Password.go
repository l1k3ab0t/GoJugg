package Objects

import (
	"github.com/coreos/bbolt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type PWHandler struct {
	db string
}

type Password struct {
	Key string
	Pw  string
}

func CreateDB(file string, pwList []Password) PWHandler {
	handler := PWHandler{file}
	db, err := bolt.Open(handler.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println(db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("PW"))
		if err != nil {
			return err
		}
		return nil
	}))

	for _, v := range pwList {
		log.Println("Key: ", v.Key, " PW: ", v.Pw)
		log.Println(db.Update(func(tx *bolt.Tx) error {
			hash, err := bcrypt.GenerateFromPassword([]byte(v.Pw), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			b := tx.Bucket([]byte("PW"))
			err = b.Put([]byte(v.Key), hash)
			return err
		}))
	}
	return handler
}

func (h PWHandler) ChangePW(oldPW Password, newPW Password) {
	db, err := bolt.Open(h.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if h.Check(oldPW) == true && oldPW.Key == newPW.Key {
		log.Println(db.Update(func(tx *bolt.Tx) error {
			hash, err := bcrypt.GenerateFromPassword([]byte(newPW.Pw), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			b := tx.Bucket([]byte("PW"))
			err = b.Put([]byte(oldPW.Key), hash)
			return err
		}))
	}

}

func (h PWHandler) Check(PW Password) bool {
	pwCorrect := false
	db, err := bolt.Open(h.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println(db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PW"))
		v := b.Get([]byte(PW.Key))
		if err := bcrypt.CompareHashAndPassword(v, []byte(PW.Pw)); err != nil {
			return err
		} else {
			pwCorrect = true
			return nil
		}
	}))
	return pwCorrect
}
