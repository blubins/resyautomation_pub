package brute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"resysniper/src/accgen"
	"resysniper/src/colors"
	"resysniper/src/utils"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

// Logs error to ./config/logs.txt and stdout in cyan
func (rt *ResyTask) Log(s string) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("%s[%s] [%s] [%s] - %s%s",
		colors.CYAN, rt.TaskIdentifier, currentDate, "BRUTE", s, colors.RESET)

	logMessageForFile := fmt.Sprintf("[%s] [%s] [%s] - %s",
		rt.TaskIdentifier, currentDate, "BRUTE", s)
	err := utils.AppendFileSync("./config/logs.txt", logMessageForFile+"\n")
	if err != nil {
		return
	}
	fmt.Printf("%s\n", logMessage)
}

// Logs error to ./config/logs.txt and stdout in red
func (rt *ResyTask) LogError(s string) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("%s[%s] [%s] [%s] - %s%s",
		colors.RED, rt.TaskIdentifier, currentDate, "BRUTE", s, colors.RESET)

	logMessageForFile := fmt.Sprintf("[%s] [%s] [%s] - %s",
		rt.TaskIdentifier, currentDate, "BRUTE", s)

	err := utils.AppendFileSync("./config/logs.txt", logMessageForFile+"\n")
	if err != nil {
		return
	}
	fmt.Printf("%s\n", logMessage)
}

// Logs success to ./config/logs.txt and stdout in green
func (rt *ResyTask) LogSuccess(s string) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("%s[%s] [%s] [%s] - %s%s",
		colors.GREEN, rt.TaskIdentifier, currentDate, "BRUTE", s, colors.RESET)

	logMessageForFile := fmt.Sprintf("[%s] [%s] [%s] - %s",
		rt.TaskIdentifier, currentDate, "BRUTE", s)

	err := utils.AppendFileSync("./config/logs.txt", logMessageForFile+"\n")
	if err != nil {
		return
	}
	fmt.Printf("%s\n", logMessage)
}

// Loads a proxy into the tasks proxy struct
func (rt *ResyTask) LdrProxy() error {
	proxyFileText, err := utils.OpenConfigFile("./config/proxies.txt")
	if err != nil {
		return err
	}
	proxies := strings.Split(proxyFileText, "\n")
	if len(proxies) == 1 {
		return nil
	}
	proxyOut := strings.Split(proxies[rand.Intn(len(proxies))], ":")
	rt.Proxy.Host = proxyOut[0]
	rt.Proxy.Port = proxyOut[1]
	rt.Proxy.Username = proxyOut[2]
	rt.Proxy.Password = proxyOut[3]
	rt.Proxy.ProxyHttpFormat = fmt.Sprintf("http://%s:%s@%s:%s",
		rt.Proxy.Username, rt.Proxy.Password, rt.Proxy.Host, rt.Proxy.Port)
	return nil
}

// Initilizes the http client for the task
func (rt *ResyTask) InitHttpClient() error {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithProxyUrl(rt.Proxy.ProxyHttpFormat),
	}
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return err
	}
	rt.Client = client
	rt.Client.SetProxy(rt.Proxy.ProxyHttpFormat)
	return nil
}

// validates a given tasks data, call after initialization
func (rt *ResyTask) ValidateTaskData() error {
	// validate very basically if email is valid
	if len(rt.EmailAddress) == 0 {
		return fmt.Errorf("email address invalid, too short length 0")
	}
	if !strings.Contains(rt.EmailAddress, "@") {
		return fmt.Errorf("invalid email address no @ sign")
	}
	// at bare minimum at least make sure field is filled
	if len(rt.XResyAuthToken) == 0 {
		return fmt.Errorf("x resy auth token invalid, too short length 0")
	}
	// minimum delay of 2s
	if rt.Delay < 2000 {
		rt.Delay = 2000
	}
	// make sure day format is correct
	if strings.Count(rt.Day, "-") != 2 {
		return fmt.Errorf("invalid day field, must follow 2024-06-12 format")
	}

	return nil
}

// Gets the venues name by id and sets it into rt.Venue.Name
func (rt *ResyTask) GetVenueData() error {
	url := fmt.Sprintf("https://api.resy.com/2/config?venue_id=%s", rt.RestaurantId)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":                {"application/json, text/plain, */*"},
		"accept-encoding":       {"gzip, deflate, br, zstd"},
		"accept-language":       {"en-US,en;q=0.9"},
		"authorization":         {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""},
		"cache-control":         {"no-cache"},
		"origin":                {"https://resy.com"},
		"pragma":                {"no-cache"},
		"priority":              {"u=1, i"},
		"referer":               {"https://resy.com"},
		"sec-ch-ua":             {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":      {"?0"},
		"sec-ch-ua-platform":    {"\"Windows\""},
		"sec-fetch-dest":        {"empty"},
		"sec-fetch-mode":        {"cors"},
		"sec-fetch-site":        {"same-site"},
		"user-agent":            {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-origin":              {"https://resy.com"},
		"x-resy-auth-token":     {rt.XResyAuthToken},
		"x-resy-universal-auth": {rt.XResyAuthToken},
	}

	resp, err := rt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// success 200
	resByteArr, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		var res BookingData
		err = json.Unmarshal(resByteArr, &res)
		if err != nil {
			return err
		}
		rt.Venue.Name = res.Name
		return nil
	}

	return fmt.Errorf(fmt.Sprintf("non 200 response received: %d", resp.StatusCode))
}

// Will grab random reservation slot from users day, party size and venue id
func (rt *ResyTask) GetRandomReservation() error {
	// day format = 2024-06-12
	url := fmt.Sprintf("https://api.resy.com/4/find?lat=0&long=0&day=%s&party_size=%d&venue_id=%s&sort_by=available",
		rt.Day, rt.PartySize, rt.RestaurantId)

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":             {"application/json, text/plain, */*"},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"en-US,en;q=0.9"},
		"authorization":      {`ResyAPI api_key="VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5"`},
		"cache-control":      {"no-cache"},
		"origin":             {"https://resy.com"},
		"pragma":             {"no-cache"},
		"priority":           {"u=1,i"},
		"referer":            {"https://resy.com"},
		"sec-ch-ua":          {"\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"Windows\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-origin":           {"https://resy.com"},
	}

	resp, err := rt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resByteArr, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error getting reservations non 200 code received: %s", resp.Status)
	}

	var res Root
	err = json.Unmarshal(resByteArr, &res)
	if err != nil {
		return err
	}

	// results.venues[0].venue.name
	if len(res.Results.Venues) > 0 {
		if len(res.Results.Venues[0].Slots) == 0 {
			return fmt.Errorf("no reservation slots found")
		}
	} else {
		return fmt.Errorf("no reservation slots found")
	}

	rt.Venue.Name = res.Results.Venues[0].VenueWrapper.Name
	// picks random slot from slot 0 and slice length
	randomSlot := rand.Intn(len(res.Results.Venues[0].Slots))
	rt.Venue.BookingToken = res.Results.Venues[0].Slots[randomSlot].Config.Token
	rt.Venue.RoomType = res.Results.Venues[0].Slots[randomSlot].Config.Type
	rt.Venue.Date = res.Results.Venues[0].Slots[randomSlot].Date.Start

	return nil
}

// Grabs the booking token for a given reservation config and users day and party size parameters,
// will then set the field in the task field under bookingTokenFinal
func (rt *ResyTask) GetBookingToken() error {
	url := "https://api.resy.com/3/details"
	payload := fmt.Sprintf(`{"commit":1,"config_id":"%s","day":"%s","party_size":%d}`,
		rt.Venue.BookingToken, rt.Venue.Date, rt.PartySize)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBufferString(payload),
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"content-type":          {"application/json"},
		"accept":                {"application/json"},
		"accept-encoding":       {"gzip, deflate, br"},
		"accept-language":       {"en-US,en;q=0.9"},
		"authorization":         {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""},
		"referer":               {"https://widgets.resy.com"},
		"sec-ch-ua":             {"\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\""},
		"user-agent":            {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-resy-auth-token":     {rt.XResyAuthToken},
		"x-resy-universal-auth": {rt.XResyAuthToken},
	}

	resp, err := rt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// success = 201 Created
	if resp.StatusCode != 201 {
		return fmt.Errorf(fmt.Sprintf("non 200 response received: %s", resp.Status))
	}

	var respJson Obj
	resByteArr, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resByteArr, &respJson)
	if err != nil {
		return err
	}

	if len(respJson.BookToken.Value) == 0 {
		return fmt.Errorf("could not get valid booking token")
	}
	rt.Venue.BookingTokenFinal = respJson.BookToken.Value

	return nil
}

// Posts a reservation given the previous methods have been followed
// for getting the booking token
func (rt *ResyTask) PostReservation() error {
	url := "https://api.resy.com/3/book"

	paymentId := fmt.Sprintf(`{"id":%s}`, rt.PaymentId)

	encodedPaymentId := utils.UrlEncode(paymentId)
	encodedBookingToken := utils.UrlEncode(rt.Venue.BookingTokenFinal)
	encodedBookingToken = strings.ReplaceAll(encodedBookingToken, "|", "%7C")
	payloadString := fmt.Sprintf(`book_token=%s&struct_payment_method=%s&source_id=resy.com-venue-details&venue_marketing_opt_in=0`,
		encodedBookingToken, encodedPaymentId)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBufferString(payloadString),
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":                {"application/json, text/plain, */*"},
		"authorization":         {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""},
		"cache-control":         {"no-cache"},
		"content-type":          {"application/x-www-form-urlencoded"},
		"origin":                {"https://widgets.resy.com"},
		"referer":               {"https://widgets.resy.com"},
		"user-agent":            {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-resy-auth-token":     {rt.XResyAuthToken},
		"x-resy-universal-auth": {rt.XResyAuthToken},
	}

	resp, err := rt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// success 200
	if resp.StatusCode != 201 {
		return fmt.Errorf("error posting reservation non 200 status received: %s", resp.Status)
	}

	return nil
}

// Will grab default payment stripe id from users account via their x resy auth token and set it to the task paymentid field
func (rt *ResyTask) GetPaymentId() error {
	url := "https://api.resy.com/2/user"

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":                {"application/json, text/plain, */*"},
		"accept-encoding":       {"gzip, deflate, br, std"},
		"accept-language":       {"en-US,en;q=0.9"},
		"authorization":         {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""},
		"cache-control":         {"no-cache"},
		"origin":                {"https://resy.com"},
		"pragma":                {"no-cache"},
		"priority":              {"u=1, i"},
		"referer":               {"https://resy.com/"},
		"sec-ch-ua":             {"\"Google Chrome\";v=\"120\", \"Chromium\";v=\"120\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":      {"?0"},
		"sec-ch-ua-platform":    {"\"Windows\""},
		"sec-fetch-dest":        {"empty"},
		"sec-fetch-mode":        {"cors"},
		"sec-fetch-site":        {"same-site"},
		"user-agent":            {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-origin":              {"https://resy.com"},
		"x-resy-auth-token":     {rt.XResyAuthToken},
		"x-resy-universal-auth": {rt.XResyAuthToken},
	}

	resp, err := rt.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	//success is 200
	if resp.StatusCode != 200 {
		return fmt.Errorf("non 200 status received: %d", resp.StatusCode)
	}

	resByteArr, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resOut accgen.PaymentMethodList
	json.Unmarshal(resByteArr, &resOut)

	if len(resOut.PaymentMethod) == 0 {
		return fmt.Errorf("no payment methods found on account")
	}
	rt.PaymentId = strconv.Itoa(resOut.PaymentMethod[0].ID)

	return nil
}

// TODO finish login function so users do not have to mess around with tokens and shit within config besides user:pass
func (rt *ResyTask) Login() error {
	payload := fmt.Sprintf("email=%s&password=%s", utils.UrlEncode(rt.EmailAddress), rt.Password)
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.resy.com/3/auth/password",
		bytes.NewBufferString(payload),
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"authorization": {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""},
		"accept":        {"application/json, text/plain, */*"},
		"content-type":  {"application/x-www-form-urlencoded"},
		"origin":        {"https://resy.com"},
		// "accept-encoding": {"gzip, deflate, br, zstd"},
		// "accept-language": {"en-US,en;q=0.9"},
		// "cache-control":   {"no-cache"},
		// "pragma":          {"no-cache"},
		// "priority":        {"u=1,i"},
		// "referer":            {"https://resy.com/"},
		// "sec-ch-ua":          {"\"Google Chrome\";v=\"120\", \"Chromium\";v=\"120\", \"Not-A.Brand\";v=\"24\""},
		// "sec-ch-ua-mobile":   {"?0"},
		// "sec-ch-ua-platform": {"\"Windows\""},
		// "sec-fetch-dest":     {"empty"},
		// "sec-fetch-mode":     {"cors"},
		// "sec-fetch-site":     {"same-site"},
		// "user-agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		// "x-origin":           {"https://resy.com"},
	}

	resp, err := rt.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// success is 200
	resBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	rt.Log(string(resBytes))
	if resp.StatusCode == 200 {
	}

	return nil
}

func (rt *ResyTask) Kill() {
	rt.Running = false
	close(rt.StopChannel)
}
