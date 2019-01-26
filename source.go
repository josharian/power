package power

type Source uint8

//go:generate stringer -type=Source
const (
	Unknown Source = iota
	Battery
	AC
	UPS
)
