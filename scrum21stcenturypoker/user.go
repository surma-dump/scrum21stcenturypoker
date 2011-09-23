package scrum21stcenturypoker

import (
	"http"
	"time"
	"rand"
	"fmt"
	"os"
)

var (
	r = rand.New(rand.NewSource(time.Nanoseconds()))
)

const (
	ID_LEN      = 16
	COOKIE_NAME = "s21cp_user"
)

type User string

func getUser(w http.ResponseWriter, r *http.Request) (User, os.Error) {
	cookie := getUserCookie(r)
	if cookie == nil {
		cookie = generateNewUserCookie()
		http.SetCookie(w, cookie)
	}
	return User(cookie.Value), nil
}

func getUserCookie(r *http.Request) *http.Cookie {
	for _, cookie := range r.Cookie {
		if cookie.Name == COOKIE_NAME {
			return cookie
		}
	}
	return nil
}

func generateUserID() (id string) {
	for i := 0; i < ID_LEN; i++ {
		id += fmt.Sprintf("%02x", r.Uint32()%256)
	}
	return
}

func generateNewUserCookie() *http.Cookie {
	return &http.Cookie{
		Name:  COOKIE_NAME,
		Value: generateUserID(),
		Path:  "/",
		// 10 Days
		Expires: *(time.SecondsToUTC(time.Nanoseconds()/1e9 + 10*24*60*60)),
	}
}
