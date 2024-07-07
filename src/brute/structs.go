package brute

import (
	"time"

	tls_client "github.com/bogdanfinn/tls-client"
)

type ResyTask struct {
	Client         tls_client.HttpClient
	Proxy          Proxy
	Venue          BookingData
	TaskIdentifier string
	EmailAddress   string
	Password       string
	XResyAuthToken string
	Day            string
	Time           string
	RoomName       string
	PartySize      int
	PaymentId      string
	RestaurantId   string
	Delay          int
	TaskSchedule   string
	RunTime        int
	WebhookUrl     string
	TaskStartTime  time.Time
	StopChannel    chan struct{}
	Running        bool
}

type BookingData struct {
	Name              string
	BookingToken      string
	RoomType          string
	Date              string
	BookingTokenFinal string
}

type Proxy struct {
	ProxyHttpFormat string // username:password@host:port
	Host            string
	Port            string
	Username        string
	Password        string
}
