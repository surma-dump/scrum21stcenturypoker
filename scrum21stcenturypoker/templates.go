package scrum21stcenturypoker

import (
	"template"
)

var (
	roomchooser_template = template.MustParseFile("roomchooser.html", nil)
	room_template        = template.MustParseFile("room.html", nil)
)
