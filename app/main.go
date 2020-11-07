package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

type TokenStruct struct {
	Token string `json:"token"`
}

var privateKey []byte

func init() {
	privateKey, _ = ioutil.ReadFile("./keys/id_rsa")
}

func main() {
	http.HandleFunc("/", initial)
	http.HandleFunc("/token", getToken)

	http.ListenAndServe(":80", nil)
}

func getToken(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	var secret = r.URL.Query().Get("secret")

	db, err := sql.Open("sqlite3", "./database/JWT_database.db")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(strings.Trim(id, " ")) == 0 || len(strings.Trim(secret, " ")) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No se recibieron los parametros obligatorios.")

		return
	}

	var count int
	err = db.QueryRow("SELECT COUNT(1) FROM Users WHERE id=? AND secret=?", id, secret).Scan(&count)

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Usuario o secret no válidos.")

		return
	}

	rows, _ := db.QueryRow("SELECT scopes FROM Users WHERE id=? AND secret=?", id, secret)
	var scopes string

	for rows.Next() {
		err = rows.Scan(&scopes)
	}

	//* Se genera el token *//
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	claims := make(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Unix() + 36000
	claims["iss"] = "sa_g1"
	claims["scope"] = strings.Split(scopes, ",")

	token.Claims = claims
	tokenString, _ := token.SignedString(key)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var structToken = TokenStruct{tokenString}
	json.NewEncoder(w).Encode(structToken)
}

func initial(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Servidor de tokens - GRUPO1 - Deploy Final")
}
