package client

import (
	"oras.land/oras-go/v2/registry/remote/auth"
)

func GetDefaultAuthClient() *auth.Client {
	return &auth.Client{
		Credential: auth.StaticCredential("", auth.EmptyCredential),
	}
}