package main

type Code interface {
	getLocalPath() string
	getRemotePath() string
	Lock()
	Unlock()
}
