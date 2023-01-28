package room

type ID int

type Room struct {
	ID        ID
	PlayerIDs map[PlayerID]struct{}
}
