package main

import (
	"errors"
)

var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

type Avatar struct {
	GetAvatarURL(c *client)	(string, error)
}