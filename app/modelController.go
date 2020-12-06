package app

import (
	"log"
)

type Message struct {
	Text string
}




func (s Session)SendToFirebase(m, d string) {

	_, _, err := s.Fs.Collection("tweets").Add(s.Ctx, map[string]interface{}{
		"date": d,
		"message":  m,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)

	}
}