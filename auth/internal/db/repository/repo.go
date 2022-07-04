package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"mts/auth-service/internal/db"
)

func init() {

	data, err := ioutil.ReadFile("./user_config.json")
	if err != nil {
		log.Fatalln(err)
		return
	}

	var obj []db.AuthData
	err = json.Unmarshal(data, &obj)
	if err != nil {
		log.Fatalln("error:", err)
	}

	for _, v := range obj {
		Repository[v.Login] = v
	}
}

var Repository = map[string]db.AuthData{}
