package lvl

import (
	"encoding/json"

	"api.cloud.io/internal/models"
)
import "github.com/syndtr/goleveldb/leveldb"

func AddUser(user models.User) error {
	db, err := leveldb.OpenFile("./lvldb/users", nil)
	if err != nil {
		return err
	}
	defer db.Close()
	bt, _ := json.Marshal(user)
	err = db.Put([]byte(user.Username), bt, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(username string) (models.User, error) {
	db, err := leveldb.OpenFile("./lvldb/users", nil)
	if err != nil {
		return models.User{}, err
	}
	defer db.Close()
	data, err := db.Get([]byte(username), nil)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	err = json.Unmarshal(data, &user)
	return user, err
}
