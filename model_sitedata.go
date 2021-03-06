package main

import (
	"errors"
	"strconv"
	"strings"
)

/**
 * SiteData
 * Contains configuration for the website
 */
type siteData struct {
	Title       string
	Ip          string
	Port        int
	SessionName string
	ServerDir   string
	authMode    int
	publicMode  int

	DevMode bool
	Mode    int

	m       *model
	mPath   []string // The path in the db to this site data
	changed bool

	sessionSecret string
}

// NewSiteData returns a siteData object with the default values
func NewSiteData(m *model) *siteData {
	ret := new(siteData)
	ret.Title = "ICT GameJam"
	ret.Ip = "127.0.0.1"
	ret.Port = 8080
	ret.SessionName = "ict-gamejam"
	ret.ServerDir = "./"
	ret.mPath = []string{"site"}
	ret.m = m
	return ret
}

// Authentication Modes: Flags for which clients are able to vote
const (
	AuthModeAuthentication = iota
	AuthModeAll
	AuthModeError
)

// Mode flags for how the site is currently running
const (
	SiteModeWaiting = iota
	SiteModeVoting
	SiteModeError
)

// load the site data out of the database
// If fields don't exist in the DB, don't clobber what is already in s
func (s *siteData) LoadFromDB() error {
	if err := s.m.openDB(); err != nil {
		return err
	}
	defer s.m.closeDB()

	if title, _ := s.m.bolt.GetValue(s.mPath, "title"); strings.TrimSpace(title) != "" {
		s.Title = title
	}
	if ip, err := s.m.bolt.GetValue(s.mPath, "ip"); err == nil {
		s.Ip = ip
	}
	if port, err := s.m.bolt.GetInt(s.mPath, "port"); err == nil {
		s.Port = port
	}
	if sessionName, _ := s.m.bolt.GetValue(s.mPath, "session-name"); strings.TrimSpace(sessionName) != "" {
		s.SessionName = sessionName
	}
	if serverDir, _ := s.m.bolt.GetValue(s.mPath, "server-dir"); strings.TrimSpace(serverDir) != "" {
		s.ServerDir = serverDir
	}
	s.changed = false
	if secret, _ := s.m.bolt.GetValue(s.mPath, "session-secret"); strings.TrimSpace(secret) != "" {
		s.sessionSecret = secret
	}
	return nil
}

// Return if the site data in memory has changed
func (s *siteData) NeedsSave() bool {
	return s.changed
}

// Save the site data into the DB
func (s *siteData) SaveToDB() error {
	var err error
	if err = s.m.openDB(); err != nil {
		return err
	}
	defer s.m.closeDB()

	if err = s.m.bolt.SetValue(s.mPath, "title", s.Title); err != nil {
		return err
	}
	if err = s.m.bolt.SetValue(s.mPath, "ip", s.Ip); err != nil {
		return err
	}
	if err = s.m.bolt.SetInt(s.mPath, "port", s.Port); err != nil {
		return err
	}
	if err = s.m.bolt.SetValue(s.mPath, "session-name", s.SessionName); err != nil {
		return err
	}
	if err = s.m.bolt.SetValue(s.mPath, "server-dir", s.ServerDir); err != nil {
		return err
	}
	s.changed = false
	if err = s.m.bolt.SetValue(s.mPath, "session-secret", s.sessionSecret); err != nil {
		return err
	}
	return nil
}

// Return the Auth Mode
func (s *siteData) GetAuthMode() int {
	return s.authMode
}

// Set the auth mode
func (s *siteData) SetAuthMode(mode int) error {
	if mode < AuthModeAuthentication || mode >= AuthModeError {
		return errors.New("Invalid Authentication Mode: " + strconv.Itoa(mode))
	}
	if mode != s.authMode {
		s.authMode = mode
		s.changed = true
	}
	return nil
}

// Return the public site mode
func (s *siteData) GetPublicMode() int {
	return s.publicMode
}

// Set the public site mode
func (s *siteData) SetPublicMode(mode int) error {
	if mode < SiteModeWaiting || mode >= SiteModeError {
		return errors.New("Invalid Public Mode: " + strconv.Itoa(mode))
	}
	if mode != s.publicMode {
		s.publicMode = mode
		s.changed = true
	}
	return nil
}
