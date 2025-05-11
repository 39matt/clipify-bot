package models

type PendingVerification struct {
	Code     string
	Platform string
	Username string
}

var Verifications = map[string]PendingVerification{}
