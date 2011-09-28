package scrum21stcenturypoker

import (
	"appengine"
	"appengine/datastore"
	"appengine/channel"
)

type Room struct {
	Name  string
	Admin string
	Scale []string
}

var (
	DEFAULT_SCALE = []string{"0", "0.5", "1", "2", "3", "5", "8", "13", "21", "40", "80", "120", "Infinity"}
)

type RoomManager struct {
	ctx appengine.Context
}

func NewRoomManager(ctx appengine.Context) *RoomManager {
	return &RoomManager{
		ctx: ctx,
	}
}

func generateRoomKey(room_name string) *datastore.Key {
	return datastore.NewKey("Room", room_name, 0, nil)
}

func (this *RoomManager) RoomExists(room_name string) bool {
	room := Room{}
	key := generateRoomKey(room_name)
	e := datastore.Get(this.ctx, key, &room)
	if e != nil && e != datastore.ErrNoSuchEntity {
		panic(e)
	}
	return e == nil
}

var (
	ErrRoomExists = &ErrorData{
		InternalMessage: "Room \"%s\" exists",
		ExternalMessage: "Room \"%s\" exists",
	}
)

func (this *RoomManager) NewRoom(room_name, adminid string) Error {
	if this.RoomExists(room_name) {
		return ErrRoomExists.Format(room_name)
	}

	key := generateRoomKey(room_name)
	room := Room{
		Name:  room_name,
		Admin: adminid,
		Scale: DEFAULT_SCALE,
	}
	_, e := datastore.Put(this.ctx, key, &room)
	if e != nil {
		panic(e)
	}
	return nil
}

func (this *RoomManager) ClientIsInRoom(pc *PokerClient, room_name string) bool {
	users := this.ClientsInRoom(room_name)
	for _, user := range users {
		if user == pc.id {
			return true
		}
	}
	return false
}

func (this *RoomManager) ClientsInRoom(room_name string) []string {
	room_key := generateRoomKey(room_name)
	keys, e := datastore.NewQuery("Client").Ancestor(room_key).KeysOnly().GetAll(this.ctx, nil)
	if e != nil {
		panic(e)
	}

	users := make([]string, 0)
	for _, key := range keys {
		users = append(users, key.StringID())
	}
	return users
}

var (
	ErrNoSuchRoom = &ErrorData{
		InternalMessage: "Room \"%s\" does not exists",
		ExternalMessage: "Room \"%s\" does not exists",
	}
)

func generateClientRoomKey(id, room_name string) *datastore.Key {
	room_key := generateRoomKey(room_name)
	client_key := datastore.NewKey("Client", id, 0, room_key)
	return client_key
}

func (this *RoomManager) EnterRoom(pc *PokerClient, room_name string) (string, Error) {
	if !this.RoomExists(room_name) {
		return "", ErrNoSuchRoom.Format(room_name)
	}
	// This effectively marks the membership in the room
	this.SetClientMetaForRoom(pc, room_name)
	c, e := channel.Create(this.ctx, room_name+"/"+pc.id)
	if e != nil {
		panic(e)
	}
	return c, nil
}

func (this *RoomManager) ExitRoom(pc *PokerClient, room_name string) {
	client_key := generateClientRoomKey(pc.id, room_name)
	datastore.Delete(this.ctx, client_key)
}

func (this *RoomManager) GetClientMetaForRoom(pc *PokerClient, room_name string) {
	client_key := generateClientRoomKey(pc.id, room_name)
	e := datastore.Get(this.ctx, client_key, &pc.meta)
	if e != nil {
		panic(e)
	}
}

func (this *RoomManager) SetClientMetaForRoom(pc *PokerClient, room_name string) {
	client_key := generateClientRoomKey(pc.id, room_name)
	_, e := datastore.Put(this.ctx, client_key, &(pc.meta))
	if e != nil {
		panic(e)
	}
}

func (this *RoomManager) ClientHasVoted(pc *PokerClient, room_name string) bool {
	this.GetClientMetaForRoom(pc, room_name)
	return pc.meta.Vote != 0
}

func (this *RoomManager) Vote(pc *PokerClient, room_name string, vote int) {
	this.GetClientMetaForRoom(pc, room_name)
	pc.meta.Vote = vote
	this.SetClientMetaForRoom(pc, room_name)
}
