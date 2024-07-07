package setup

type fileList struct {
	Path    string
	Content string
}

type userConfigJson struct {
	License    string
	Webhook    string
	CaptchaKey string
}
