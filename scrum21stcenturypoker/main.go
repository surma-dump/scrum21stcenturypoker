package scrum21stcenturypoker

import (
	"http"
	"os"
	"appengine"
	"appengine/datastore"
	"fmt"
)

func init() {
	http.HandleFunc("/", roomchooser)
	http.Handle("/static/", http.FileServer("./", ""))
	http.HandleFunc("/enterRoom", enterRoom)
	http.HandleFunc("/room/", poker)
}

func error(w http.ResponseWriter, r *http.Request, prefix string, e os.Error) {
	ctx := appengine.NewContext(r)
	ctx.Errorf(prefix + ": " + e.String())
	http.Error(w, "There was an error. Sorry about that", 500)
}

func roomchooser(w http.ResponseWriter, r *http.Request) {
	_, e := getUser(w, r)
	if e != nil {
		error(w, r, "User management", e)
	}
	roomchooser_template.Execute(w, nil)
}

type Room struct {
	Name  string
	Admin User
	Scale []string
}

func enterRoom(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	e := r.ParseForm()
	if e != nil {
		error(w, r, "Form parsing", e)
		return
	}

	user, e := getUser(w, r)
	if e != nil {
		error(w, r, "User management", e)
	}
	rooms, ok := r.Form["room"]
	room_name := rooms[0] // ?
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	room_key := datastore.NewKey("Room", room_name, 0, nil)
	room := Room{}
	e = datastore.Get(ctx, room_key, &room)
	if e == nil {
		// Room exists
		http.Redirect(w, r, "/room/"+room_name, http.StatusTemporaryRedirect)
	} else if e == datastore.ErrNoSuchEntity || e == datastore.ErrInvalidEntityType {
		ctx.Infof("Creating room: \"%s\" (Reason: %s)", room_name, e.String())
		// Room has to be created
		_, e = datastore.Put(ctx, room_key, &Room{
			Name:  room_name,
			Admin: user,
			Scale: []string{"0", "0.5", "1", "2", "3", "5", "8", "13", "21", "40", "80", "120", "Infinite"},
		})
		if e != nil {
			error(w, r, "Room creation", e)
			return
		}
		http.Redirect(w, r, "/room/"+room_name, http.StatusTemporaryRedirect)
	} else {
		error(w, r, "Room existence check", e)
		return
	}
}

func poker(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	user, e := getUser(w, r)
	if e != nil {
		error(w, r, "User management", e)
	}
	room_name := r.URL.Path[len("/room/"):]
	room_key := datastore.NewKey("Room", room_name, 0, nil)
	room := Room{}
	e = datastore.Get(ctx, room_key, &room)
	if e == datastore.ErrNoSuchEntity {
		fmt.Fprintf(w, "Invalid room")
		return
	} else if e != nil {
		error(w, r, "Room entering", e)
		return
	}
	room_template.Execute(w, map[string]interface{}{
		"Name":    room.Name,
		"IsAdmin": room.Admin == user,
		"Scale":   room.Scale,
	})
}
