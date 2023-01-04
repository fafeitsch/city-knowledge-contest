package keygen

import "github.com/jaevor/go-nanoid"

var roomKeyGen, _ = nanoid.Standard(21)

func SetRoomKeyLength(length int) {
	var err error
	roomKeyGen, err = nanoid.Standard(length)
	if err != nil {
		panic(err)
	}
}

func RoomKey() string {
	return roomKeyGen()
}

var playerKeyGen, _ = nanoid.Standard(21)

func SetPlayerKeyLength(length int) {
	var err error
	playerKeyGen, err = nanoid.Standard(length)
	if err != nil {
		panic(err)
	}
}

func PlayerKey() string {
	return playerKeyGen()
}
