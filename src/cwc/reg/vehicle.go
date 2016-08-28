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
