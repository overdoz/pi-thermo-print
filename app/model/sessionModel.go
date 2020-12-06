package model

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
)

type Session struct {
	Fs *firestore.Client
	Ctx context.Context
}