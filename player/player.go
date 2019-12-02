package player

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/dgraph-io/badger"
	"gitlab.com/glatteis/earthwalker/database"
	"math/rand"
)

type PlayerSession struct {
	// UniqueIdentifier is the session identifier stored in the key.
	UniqueIdentifier string
	// CurrentGameID is game identifier the player might be currently in.
	CurrentGameID string
	// CurrentRound is the round the player is in.
	CurrentRound int
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func StorePlayerSession(session PlayerSession) error {
	err := database.GetDB().Update(func(txn *badger.Txn) error {
		var buffer bytes.Buffer
		gob.NewEncoder(&buffer).Encode(session)
		return txn.Set([]byte("session-"+session.UniqueIdentifier), buffer.Bytes())
	})

	if err != nil {
		return err
	}
	return nil
}

func LoadPlayerSession(id string) (PlayerSession, error) {
	var playerBytes []byte

	err := database.GetDB().Update(func(txn *badger.Txn) error {
		result, err := txn.Get([]byte("session-" + id))
		if err != nil {
			return err
		}

		var res []byte
		err = result.Value(func(val []byte) error {
			res = append([]byte{}, val...)
			return nil
		})

		if err != nil {
			return err
		}

		playerBytes = res
		return nil
	})

	if err == badger.ErrKeyNotFound {
		return PlayerSession{}, errors.New("this player does not exist")
	} else if err != nil {
		return PlayerSession{}, err
	}

	var foundSession PlayerSession
	gob.NewDecoder(bytes.NewBuffer(playerBytes)).Decode(&foundSession)

	return foundSession, nil
}
