package brute

type DetailsPayload struct {
	Commit    int    `json:"commit"`     // 0
	ConfigId  string `json:"config_id"`  // config id = slot token rgs://resy/50806/1667589...
	Day       string `json:"day"`        // "2024-06-29"
	PartySize int    `json:"party_size"` // 2
}
