package brute

type Obj struct {
	BookToken BookToken `json:"book_token"`
}

type BookToken struct {
	//DateExpires string `json:"date_expires"`
	Value string `json:"value"`
}
