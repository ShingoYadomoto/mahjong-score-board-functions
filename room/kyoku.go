package room

type PointType int

const (
	PointTypeNone PointType = iota
	PointType1
	PointType2
	PointType3
	PointType4
	PointTypeManGan
	PointTypeHaneMan
	PointTypeBaiMan
	PointTypeSanbaiMan
	PointTypeYakuMan
	PointTypeDoubleYakuMan
	PointTypeTripleYakuMan
)

type FuType int

const (
	FuTypeNone FuType = 0
	FuType20   FuType = 20
	FuType25   FuType = 25
	FuType30   FuType = 30
	FuType40   FuType = 40
	FuType50   FuType = 50
	FuType60   FuType = 60
	FuType70   FuType = 70
	FuType80   FuType = 80
	FuType90   FuType = 90
	FuType100  FuType = 100
	FuType110  FuType = 110
)

type FanType int

const (
	FanTypeTon FanType = iota + 1
	FanTypeNan
	FanTypeSha
	FanTypePei
)

type Kyoku struct {
	IsWin        bool
	FromPlayerID PlayerID
	ByPlayerID   PlayerID
	PointType    PointType
	FuType       FuType
}
