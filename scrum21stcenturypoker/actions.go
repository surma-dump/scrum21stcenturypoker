package scrum21stcenturypoker

import (
	"http"
	"strings"
	"strconv"
)

type Action interface {
	GetName() string
	Execute(c *PokerClient) (interface{}, Error)
}

type CreateRoomAction struct {
	Room string
}

func (*CreateRoomAction) GetName() string {
	return "Create"
}

func (this *CreateRoomAction) Execute(pc *PokerClient) (interface{}, Error) {
	roommgr := NewRoomManager(pc.ctx)
	e := roommgr.NewRoom(this.Room, pc.id)
	if e != nil {
		return nil, e
	}
	return NewSuccessMessage("Room created", nil), nil
}

type EnterRoomAction struct {
	Room string
}

func (*EnterRoomAction) GetName() string {
	return "Enter"
}

var (
	ErrAlreadyInRoom = &ErrorData{
		InternalMessage: "User %s already in room \"%s\" - tried to Enter",
		ExternalMessage: "You are already in room \"%s\"",
	}
)

func (this *EnterRoomAction) Execute(pc *PokerClient) (interface{}, Error) {
	roommgr := NewRoomManager(pc.ctx)
	if roommgr.ClientIsInRoom(pc, this.Room) {
		return nil, ErrAlreadyInRoom.FormatInternalMessage(pc.id, this.Room).FormatExternalMessage(this.Room)
	}
	c, e := roommgr.EnterRoom(pc, this.Room)
	if e != nil {
		return nil, e
	}
	return NewSuccessMessage("Entered Room", c), nil
}

type ExitRoomAction struct {
	Room string
}

var (
	ErrNotInRoom = &ErrorData{
		InternalMessage: "User %s not in room \"%s\" - tried to %s",
		ExternalMessage: "You are not in room \"%s\"",
	}
)

func (*ExitRoomAction) GetName() string {
	return "Exit"
}

func (this *ExitRoomAction) Execute(pc *PokerClient) (interface{}, Error) {
	roommgr := NewRoomManager(pc.ctx)
	if !roommgr.ClientIsInRoom(pc, this.Room) {
		return nil, ErrNotInRoom.FormatInternalMessage(pc.id, this.Room, this.GetName()).FormatExternalMessage(this.Room)
	}
	roommgr.ExitRoom(pc, this.Room)
	return NewSuccessMessage("Exited Room", this.Room), nil
}

type VoteAction struct {
	Room string
	Vote int
}

var (
	ErrMultipleVotes = &ErrorData{
		InternalMessage: "User %s tried to vote twice in \"%s\"",
		ExternalMessage: "You voted already",
	}
)

func (*VoteAction) GetName() string {
	return "Vote"
}

func (this *VoteAction) Execute(pc *PokerClient) (interface{}, Error) {
	roommgr := NewRoomManager(pc.ctx)
	if !roommgr.ClientIsInRoom(pc, this.Room) {
		return nil, ErrNotInRoom.FormatInternalMessage(pc.id, this.Room, this.GetName()).FormatExternalMessage(this.Room)
	}
	if roommgr.ClientHasVoted(pc, this.Room) {
		return nil, ErrMultipleVotes.FormatInternalMessage(pc.id, this.Room)
	}
	roommgr.Vote(pc, this.Room, this.Vote)
	return NewSuccessMessage("You voted", this.Room), nil
}

var (
	ErrInvalidResource = &ErrorData{
		InternalMessage: "Invalid resource path \"%s\"",
	}
	ErrInvalidAction = &ErrorData{
		InternalMessage: "Invalid action \"%s\"",
	}
)

func (pc *PokerClient) GetAction() (Action, Error) {
	path := pc.req.URL.Path
	if !strings.HasPrefix(path, "/rooms") {
		return nil, ErrInvalidResource.Format(path)
	}
	elems := strings.Split(path[1:], "/", -1)
	if len(elems) != 3 {
		return nil, ErrInvalidAction.Format(path)
	}
	room, action_name := elems[1], elems[2]
	return parseAction(room, action_name, pc.req.Form)
}

func parseAction(room, action_name string, data http.Values) (Action, Error) {
	switch action_name {
	case "create":
		return parseCreateRoomAction(room, action_name, data)
	case "enter":
		return parseEnterRoomAction(room, action_name, data)
	case "exit":
		return parseExitRoomAction(room, action_name, data)
	case "vote":
		return parseVoteAction(room, action_name, data)
	}
	return nil, ErrInvalidAction
}

func parseCreateRoomAction(room, action_name string, data http.Values) (*CreateRoomAction, Error) {
	return &CreateRoomAction{
		Room: room,
	}, nil
}

func parseEnterRoomAction(room, action_name string, data http.Values) (*EnterRoomAction, Error) {
	return &EnterRoomAction{
		Room: room,
	}, nil
}

func parseExitRoomAction(room, action_name string, data http.Values) (*ExitRoomAction, Error) {
	return &ExitRoomAction{
		Room: room,
	}, nil
}

func parseVoteAction(room, action_name string, data http.Values) (*VoteAction, Error) {
	vote, e := strconv.Atoi(data.Get("vote"))
	return &VoteAction{
		Room: room,
		Vote: vote,
	}, FromError("Parsing vote", e)
}
