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
		client.SendError("Parsing data", e)
		return
	}

	action, e := client.GetAction()
	if e != nil {
		client.SendError("Parsing action", e)
		return
	}

	client.Send("You did this: " + action.GetName())
}

/*func roomchooser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	_, e := getUser(ctx, w, r)
	if e != nil {
		error(w, r, "User management", e)
	}
	roomchooser_template.Execute(w, nil)
}

type Room struct {
	Name  string
	AdminID string
	Scale []string
}

func vote(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	e := r.ParseForm()
	if e != nil {
		error(w, r, "Form parsing", e)
		return
	}

	user, e := getUser(ctx, w, r)
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

	} else {
		error(w, r, "Room existence check", e)
		return
	}

	user_query := datastore.NewQuery("User").Ancestor(room_key)
	var users []User
	_, e = user_query.GetAll(ctx, &users)
	if e != nil {
		error(w, r, "Getting users in room", e)
		return
	}
	for _, ouser := range users {
		msg := fmt.Sprintf("%s voted %s", user.ID[0:4], r.Form["vote"][0])
		channel.Send(ctx, ouser.ID, msg)
	}
}

func enterRoom(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	e := r.ParseForm()
	if e != nil {
		error(w, r, "Form parsing", e)
		return
	}

	user, e := getUser(ctx, w, r)
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
			AdminID: user.ID,
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

	user_key := datastore.NewKey("User", user.ID, 0, room_key)
	_, e = datastore.Put(ctx, user_key, &user)
	if e != nil {
		error(w, r, "Adding to room", e)
		return
	}
}

func poker(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	user, e := getUser(ctx, w, r)
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
		"IsAdmin": room.AdminID == user.ID,
		"Scale":   room.Scale,
		"ChannelToken": user.Channel,
	})
}
*/
