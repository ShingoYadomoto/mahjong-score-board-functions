package room

type (
	CurrentFieldState struct {
		Fan     FanType
		Stack   int
		Deposit int
	}

	CurrentPlayerState struct {
		PlayerID PlayerID
		Fan      FanType
		Point    int
		IsRiichi bool
	}

	CurrentState struct {
		Field   CurrentFieldState
		Players []CurrentPlayerState
	}
)
