package accgen

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type UserAccGenConfig struct {
	CatchAll    string
	CCNumber    string
	ExpYear     string
	ExpMonth    string
	CVV         string
	ZipCode     string
	PhoneNumber string
	CaptchaKey  string
	Quantity    string
	Webhook     string
}

type Proxy struct {
	ProxyHttpFormat string // username:password@host:port
	Host            string
	Port            string
	Username        string
	Password        string
}

type AccGenTask struct {
	Client         tls_client.HttpClient
	Config         UserAccGenConfig
	Proxy          Proxy
	TaskIdentifier string
	FirstName      string
	LastName       string
	EmailAddress   string
	Password       string
	AccountData    *User
	StripeData     Stripe
}

type User struct {
	AccountData RegistrationResponse `json:"user"`
}

type RegistrationResponse struct {
	ID                     int    `json:"id"`
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	MobileNumber           string `json:"mobile_number"`
	EmailAddress           string `json:"em_address"`
	EmailIsVerified        int    `json:"em_is_verified"`
	MobileNumberIsVerified int    `json:"mobile_number_is_verified"`
	IsActive               int    `json:"is_active"`
	ReferralCode           string `json:"referral_code"`
	IsMarketable           int    `json:"is_marketable"`
	IsConcierge            int    `json:"is_concierge"`
	DateUpdated            int    `json:"date_updated"`
	DateCreated            int    `json:"date_created"`
	HasSetPassword         int    `json:"has_set_password"`
	ViewedGDAWelcome       bool   `json:"viewed_gda_welcome"`
	NumBookings            int    `json:"num_bookings"`
	ResySelect             int    `json:"resy_select"`
	ProfileImageURL        string `json:"profile_image_url"`
	IsGlobalDiningAccess   bool   `json:"is_global_dining_access"`
	IsRGA                  bool   `json:"is_rga"`
	GuestID                int    `json:"guest_id"`
	Token                  string `json:"token"`
	LegacyToken            string `json:"legacy_token"`
	RefreshToken           string `json:"refresh_token"`
}

type Stripe struct {
	StripePkToken   string
	ClientSecret    string
	CustomerId      string
	EphemeralKey    string
	Guid            string
	Muid            string
	Sid             string
	PaymentMethodID string
}

type SetupIntentResponse struct {
	ClientSecret string `json:"client_secret"`
	CustomerId   string `json:"customer_id"`
	EphemeralKey string `json:"ephemeral_key"`
}

type SiteConfigResponse struct {
	StripePkToken string `json:"stripe_publishable_key"`
}

type StripeUUIDsResponse struct {
	Guid string `json:"guid"`
	Muid string `json:"muid"`
	Sid  string `json:"sid"`
}

type SetStripeIntents struct {
	PaymentMethodID string `json:"payment_method"`
}

type PaymentMethod struct {
	ID int `json:"id"`
}

type PaymentMethodList struct {
	PaymentMethod []PaymentMethod `json:"payment_methods"`
}
