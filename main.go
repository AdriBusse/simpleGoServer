package main

import (
	"bytes"
	"encoding/gob"

    "errors"
	"fmt"
    "log"
    "net/http"
	"strings"

	"example.com/simpleGoServer/internal/cookies"
)

var secretKey string
var secretKeyHex []byte

type User struct {
    Username string
    Role  int
}

func main() {
	gob.Register(&User{})

	var err error
	secretKey = "secretKeySecret1"
	secretKeyHex= []byte(secretKey)
	if err != nil {
        log.Fatal(err)
    }

    // Start a web server with 3 endpoints on port 3333.
    mux := http.NewServeMux()
    mux.HandleFunc("/", pongHandler)
    mux.HandleFunc("/set", setCookieHandler)
    mux.HandleFunc("/get", getCookieHandler)

    port := ":3333"
    log.Print("Listening on port 3333...\n")
    err = http.ListenAndServe(port, mux)
    if err != nil {
        log.Fatal(err)
    }
}

func pongHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("ping pong..."))
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
    user := User{Username: "Secret User", Role: 1}

	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
        Name:     "EthicalHackingCookie",
        Value:    buf.String(),
        Path:     "/",
        MaxAge:   360000000,
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteLaxMode,
    }

	err = cookies.WriteEncrypted(w, cookie, secretKeyHex)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

    w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	gobEncodedValue, err := cookies.ReadEncrypted(r, "exampleCookie", secretKeyHex)

	if err != nil {
		switch {
        case errors.Is(err, http.ErrNoCookie):
            http.Error(w, "cookie not found", http.StatusBadRequest)
        case errors.Is(err, cookies.ErrInvalidValue):
            http.Error(w, "invalid cookie", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        return	
	}

	var user User

	reader := strings.NewReader(gobEncodedValue)

	if err := gob.NewDecoder(reader).Decode(&user); err != nil{
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
  

	fmt.Fprintf(w, "Username: %s", user.Username)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Role: %s", fmt.Sprint(user.Role))
}

