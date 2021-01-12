package main

import "github.com/lusantisuper/Mundipagg-SDK-GO/internal/req"

// NewMundipagg Create a Mundipagg object
func NewMundipagg(basicAuthUserName string, basicAuthPassword string, basicSecretAuthKey string) (*req.Login, error) {
	return req.NewMundipagg(basicAuthUserName, basicAuthPassword, basicSecretAuthKey)
}

func main() {}
