package room

type ID int

type Room struct {
	ID              ID
	UpdatedUnixTime int64
	PlayerIDs       map[PlayerID]struct{}
	Kyokus          []Kyoku
	RiichiPlayerIDs map[PlayerID]struct{}
}

func (r *Room) Joined(pid PlayerID) bool {
	if _, ok := r.PlayerIDs[pid]; ok {
		return true
	}
	return false
}

func (r *Room) CurrentState() (*CurrentState, error) {
	return &CurrentState{
		Field: CurrentFieldState{
			Fan:     FanTypeTon,
			Stack:   2,
			Deposit: 3,
		},
		Players: []CurrentPlayerState{
			{PlayerID: 1, Fan: FanTypeTon, Point: 25000, IsRiichi: false},
			{PlayerID: 2, Fan: FanTypeNan, Point: 25000, IsRiichi: false},
			{PlayerID: 3, Fan: FanTypeSha, Point: 25000, IsRiichi: false},
			{PlayerID: 4, Fan: FanTypePei, Point: 25000, IsRiichi: false},
		},
	}, nil
}
