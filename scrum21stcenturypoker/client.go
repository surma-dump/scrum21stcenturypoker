package scrum21stcenturypoker

import (
	"http"
	"os"
	"rand"
	"fmt"
	"time"
	"appengine"
)

type PokerClient struct {
	ctx appengine.Context
	http.ResponseWriter
	req *http.Request
	id  string
}

func (pc *PokerClient) SendError(prefix string, e os.Error) {
	pc.ctx.Errorf("%s: %s", prefix, e.String())
	pc.WriteHeader(500)
	pc.Send("NAIN!")
}

func (pc *PokerClient) Send(msg string) {
	fmt.Fprint(pc, msg)
}

func (pc *PokerClient) RenewIdCookie() {
	cookie := newIdCookie(pc.id)
	http.SetCookie(pc, cookie)
}

func NewPokerClient(w http.ResponseWriter, r *http.Request) *PokerClient {
	id, _ := getUserId(r)
	client := &PokerClient{
		ctx: appengine.NewContext(r),
		id:  id,
		req: r,
	}
	client.ResponseWriter = w
	client.RenewIdCookie()
	return client
}

const (
	ID_LEN    = 16
	ID_COOKIE = "s21cp_id"
)

func getUserId(r *http.Request) (string, bool) {
	for _, cookie := range r.Cookie {
		if cookie.Name == ID_COOKIE && len(cookie.Value) == ID_LEN {
			return cookie.Value, true
		}
	}
	// No (valid) cookie found
	return generateNewId(), false
}

func newIdCookie(id string) *http.Cookie {
	return &http.Cookie{
		Name:  ID_COOKIE,
		Value: id,
		Path:  "/",
		// 10 Days
		Expires: *(time.SecondsToUTC(time.Nanoseconds()/1e9 + 10*24*60*60)),
	}
}

var (
	r = rand.New(rand.NewSource(time.Nanoseconds()))
)

func generateNewId() (id string) {
	for i := 0; i < ID_LEN; i++ {
		id += fmt.Sprintf("%02x", r.Uint32()%256)
	}
	return
}
