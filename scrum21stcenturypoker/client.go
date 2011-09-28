package scrum21stcenturypoker

import (
	"http"
	"rand"
	"fmt"
	"time"
	"appengine"
	"json"
)

func (this *PokerClient) HandleError(e Error) {
	this.ctx.Errorf(e.String())
	if e.HasUserMessage() {
		this.SendError(e.UserMessage())
	} else {
		this.SendError("An error occured!")
	}
}

func (this *PokerClient) SendError(msg string) {
	this.WriteHeader(500)
	this.SendData(NewErrorMessage(msg))
}

func (this *PokerClient) SendData(data interface{}) {
	serial, e := json.Marshal(data)
	if e != nil {
		panic(e)
	}

	this.Write(serial)
}

func (this *PokerClient) RenewIdCookie() {
	cookie := newIdCookie(this.id)
	http.SetCookie(this, cookie)
}

type PokerClient struct {
	ctx appengine.Context
	http.ResponseWriter
	req  *http.Request
	id   string
	meta PokerClientMeta
}

type PokerClientMeta struct {
	Name string
	Vote int
}

func NewPokerClient(w http.ResponseWriter, r *http.Request) *PokerClient {
	id, _ := getUserId(r)
	client := &PokerClient{
		ctx: appengine.NewContext(r),
		id:  id,
		req: r,
		meta: PokerClientMeta{
			Name: id,
		},
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
		if cookie.Name == ID_COOKIE && len(cookie.Value) == ID_LEN*2 {
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
