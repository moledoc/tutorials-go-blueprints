package main

import (
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"path"
	"strings"
)

// ErrNoAvatarURL is the  error that is returned when the Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing user profile picutres.
type Avatar interface {
	// GetAvatarURL get the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get a URL for the spcified client.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	return urlStr, nil
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	userid, ok := c.userData["userid"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	useridStr, ok := userid.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	m := md5.New()
	io.WriteString(m, strings.ToLower(useridStr))
	return "//www.gravatar.com/avatar/" + useridStr, nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	userid, ok := c.userData["userid"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	useridStr, ok := userid.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := path.Match(useridStr+"*", file.Name()); match {
			return "/avatars/" + file.Name(), nil
		}
	}
	return "", ErrNoAvatarURL
}
