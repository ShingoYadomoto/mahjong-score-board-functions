package room

type ID int

type Room struct {
	ID        ID
	PlayerIDs map[PlayerID]struct{}
}

func (r *Room) Joined(pid PlayerID) bool {
	if _, ok := r.PlayerIDs[pid]; ok {
		return true
	}
	return false
}
