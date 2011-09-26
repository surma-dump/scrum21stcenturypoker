package scrum21stcenturypoker

import (
	"http"
	// "os"
	// "appengine"
)

func init() {
	http.Handle("/", http.FileServer("static", ""))
	http.HandleFunc("/rooms/", Poker)
}

func Poker(w http.ResponseWriter, r *http.Request) {
	client := NewPokerClient(w, r)
	e := r.ParseForm()
	if e != nil {
		client.HandleError(FromError("Parsing data", e))
		return
	}

	action, e := client.GetAction()
	if e != nil {
		client.HandleError(FromError("Parsing action", e))
		return
	}

	data, err := action.Execute(client)
	if err != nil {
		client.HandleError(err)
	} else {
		client.SendData(data)
	}
}
