package scrum21stcenturypoker

import (
	"http"
	"os"
	"strings"
	"strconv"
)

type Action interface {
	GetName() string
}

type EnterRoomAction struct {
	Room string
}

func (EnterRoomAction) GetName() string {
	return "Enter"
}

type ExitRoomAction struct {
	Room string
}

func (ExitRoomAction) GetName() string {
	return "Exit"
}

type VoteAction struct {
	Room string
	Vote int
}

func (VoteAction) GetName() string {
	return "Vote"
}

var (
	ErrInvalidResource = os.NewError("Invalid resource path")
	ErrInvalidAction   = os.NewError("Invalid action path")
)

func (pc *PokerClient) GetAction() (Action, os.Error) {
	path := pc.req.URL.Path
	if !strings.HasPrefix(path, "/rooms") {
		return nil, ErrInvalidResource
	}
	elems := strings.Split(path[1:], "/", -1)
	if len(elems) != 3 {
		return nil, ErrInvalidAction
	}
	room, action_name := elems[1], elems[2]
	return parseAction(room, action_name, pc.req.Form)
}

func parseAction(room, action_name string, data http.Values) (Action, os.Error) {
	switch action_name {
	case "enter":
		return parseEnterRoomAction(room, action_name, data)
	case "exit":
		return parseExitRoomAction(room, action_name, data)
	case "vote":
		return parseVoteAction(room, action_name, data)
	}
	return nil, ErrInvalidAction
}

func parseEnterRoomAction(room, action_name string, data http.Values) (EnterRoomAction, os.Error) {
	return EnterRoomAction{
		Room: room,
	}, nil
}

func parseExitRoomAction(room, action_name string, data http.Values) (ExitRoomAction, os.Error) {
	return ExitRoomAction{
		Room: room,
	}, nil
}

func parseVoteAction(room, action_name string, data http.Values) (VoteAction, os.Error) {
	vote, e := strconv.Atoi(data.Get("vote"))
	return VoteAction{
		Room: room,
		Vote: vote,
	}, e
}
