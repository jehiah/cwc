package reg

//go:generate stringer -type=Vehicle

type Vehicle int

const (
	Unknown Vehicle = 1 << iota
	Taxi
	FHV
	Other
)

func PossibleTaxi(license string) bool {
	return len(license) == 4
}

func (v Vehicle) IncludesTaxi() bool {
	return v&Taxi != 0
}
func (v Vehicle) IncludesFHV() bool {
	return v&FHV != 0
}
