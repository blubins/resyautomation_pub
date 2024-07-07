package brute

type Root struct {
	Results Venues `json:"results"`
}

type Venues struct {
	Venues []Venue `json:"venues"`
}

type Venue struct {
	Slots        []Slot       `json:"slots"`
	VenueWrapper VenueWrapper `json:"venue"`
}
type VenueWrapper struct {
	Name string `json:"name"`
}

type Slot struct {
	Availability Availability `json:"availability"`
	Payment      Payment      `json:"payment"`
	Config       Config       `json:"config"`
	Date         Date         `json:"date"`
}

type Config struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	Token string `json:"token"`
}

type Date struct {
	End   string `json:"end"`
	Start string `json:"start"`
}

type Availability struct {
	Id int `json:"id"`
}

type Payment struct {
	IsPaid           bool    `json:"is_paid"`
	CancellationFee  float32 `json:"cancellation_fee"`
	TimeCancelCutOff string  `json:"time_cancel_cut_off"`
}
