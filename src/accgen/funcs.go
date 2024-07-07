package accgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"resysniper/src/colors"
	"resysniper/src/utils"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func (agt *AccGenTask) Log(s string) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("%s[%s] [%s] [%s] - %s%s",
		colors.CYAN, agt.TaskIdentifier, currentDate, "ACC-GEN", s, colors.RESET)

	logMessageForFile := fmt.Sprintf("[%s] [%s] [%s] - %s",
		agt.TaskIdentifier, currentDate, "ACC-GEN", s)
	err := utils.AppendFileSync("./config/logs.txt", logMessageForFile+"\n")
	if err != nil {
		return
	}
	fmt.Printf("%s\n", logMessage)
}

func (agt *AccGenTask) LogError(s string) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("%s[%s] [%s] [%s] - %s%s",
		colors.RED, agt.TaskIdentifier, currentDate, "ACC-GEN", s, colors.RESET)

	logMessageForFile := fmt.Sprintf("[%s] [%s] [%s] - %s",
		agt.TaskIdentifier, currentDate, "ACC-GEN", s)

	err := utils.AppendFileSync("./config/logs.txt", logMessageForFile+"\n")
	if err != nil {
		return
	}
	fmt.Printf("%s\n", logMessage)
}

func (agt *AccGenTask) LogSuccess(s string) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("%s[%s] [%s] [%s] - %s%s",
		colors.GREEN, agt.TaskIdentifier, currentDate, "ACC-GEN", s, colors.RESET)

	logMessageForFile := fmt.Sprintf("[%s] [%s] [%s] - %s",
		agt.TaskIdentifier, currentDate, "ACC-GEN", s)

	err := utils.AppendFileSync("./config/logs.txt", logMessageForFile+"\n")
	if err != nil {
		return
	}
	fmt.Printf("%s\n", logMessage)
}

func (agt *AccGenTask) LdrProxy() error {
	proxyFileText, err := utils.OpenConfigFile("./config/proxies.txt")
	if err != nil {
		return err
	}
	proxies := strings.Split(proxyFileText, "\n")
	if len(proxies) == 1 {
		return nil
	}
	proxyOut := strings.Split(proxies[rand.Intn(len(proxies))], ":")
	agt.Proxy.Host = proxyOut[0]
	agt.Proxy.Port = proxyOut[1]
	agt.Proxy.Username = proxyOut[2]
	agt.Proxy.Password = proxyOut[3]
	agt.Proxy.ProxyHttpFormat = fmt.Sprintf("http://%s:%s@%s:%s",
		agt.Proxy.Username, agt.Proxy.Password, agt.Proxy.Host, agt.Proxy.Port)
	return nil
}

func (c *AccGenTask) ValidateAccGen() error {
	// validate the catchall
	if len(c.Config.CatchAll) < 3 {
		return fmt.Errorf("catchall field invalid")
	}
	// make sure @ sign is formatted correctly
	if !strings.Contains(c.Config.CatchAll, "@") {
		c.Config.CatchAll = "@" + c.Config.CatchAll
	}
	// validate cc length is not empty field
	if len(c.Config.CCNumber) != 0 {
		// validate exp year is of length 2
		if len(c.Config.ExpYear) != 2 {
			return fmt.Errorf("expiration year invalid, must be 2 digits in length i.e 27")
		}
		// validate exp month is of length 2
		if len(c.Config.ExpMonth) != 2 {
			return fmt.Errorf("expiration month invalid, must be of 2 digits in length i.e 07")
		}
		// validate month is within range 1-12
		intMonth, _ := strconv.Atoi(c.Config.ExpMonth)
		if intMonth > 12 || intMonth < 1 {
			return fmt.Errorf("expiration month invalid, out of range 01-12")
		}
		// validate cvv is 3 or 4 characters long
		if len(c.Config.CVV) < 2 || len(c.Config.CVV) > 4 {
			return fmt.Errorf("cvv invalid, field either empty or too long")
		}
		// put the cc number into the correct format
		// from whatever user has inputted
		invalidChars := []string{" ", "+", "-"}
		for _, char := range invalidChars {
			c.Config.CCNumber = strings.ReplaceAll(c.Config.CCNumber, char, "")
		}
		var outString []rune
		for i, letter := range c.Config.CCNumber {
			outString = append(outString, letter)
			if (i+1)%4 == 0 && i != len(c.Config.CCNumber)-1 {
				outString = append(outString, '+')
			}
		}
		c.Config.CCNumber = string(outString)
	}
	// validate zipcode is not empty
	if len(c.Config.ZipCode) < 3 {
		return fmt.Errorf("zipcode invalid, too short")
	}
	// validate phone number is not empty
	if len(c.Config.PhoneNumber) < 3 {
		return fmt.Errorf("phone number invalid, too short")
	}
	// validate quantityfield is greater than 0
	intQuantity, _ := strconv.Atoi(c.Config.Quantity)
	if intQuantity < 1 {
		return fmt.Errorf("invalid quantity field must be >= 1")
	}

	//now set faker values for rest of config
	fake := faker.New()
	c.FirstName = fake.Person().FirstName()
	c.LastName = fake.Person().LastName()
	if len(c.Config.PhoneNumber) == 3 {
		c.Config.PhoneNumber += fake.Numerify("#######")
	}
	// set url encode manually to + sign and prefix 1
	c.Config.PhoneNumber = "%2B1" + c.Config.PhoneNumber
	c.EmailAddress = c.FirstName + c.LastName + fake.Numerify("#####") + c.Config.CatchAll
	c.Password = fake.Bothify("RS?????###!")

	return nil
}

// initalizes a client that will be used in subsequent member funcs
func (agt *AccGenTask) InitHttpClient() error {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithProxyUrl(agt.Proxy.ProxyHttpFormat),
	}
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return err
	}
	agt.Client = client
	// double make sure proxy is set
	agt.Client.SetProxy(agt.Proxy.ProxyHttpFormat)
	return nil
}

// will post to the resy registration endpoint with a given string
func (agt *AccGenTask) PostReigstration(payload string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.resy.com/2/user/registration",
		bytes.NewBufferString(payload),
	)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"accept":          {"application/json, text/plain, */*"},
		"accept-encoding": {"gzip, deflate, br, zstd"},
		"accept-language": {"en-US,en;q=0.9"},
		"authorization":   {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""}, // hardcoded api key intentional
		"cache-control":   {"no-cache"},
		"content-type":    {"application/x-www-form-urlencoded"},
		"origin":          {"https://resy.com"},
		"pragma":          {"no-cache"},
		"priority":        {"u=1, i"},
		"referer":         {"https://resy.com"},
		"user-agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
	}
	resp, err := agt.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		resByteArr, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var resOut User
		err = json.Unmarshal(resByteArr, &resOut)
		if err != nil {
			return err
		}
		agt.AccountData = &resOut
		return nil
	}

	return fmt.Errorf("received non 201 status %d", resp.StatusCode)
}

// Initalizes the account gen config
func (agt *AccGenTask) InitAccountGenConfig() error {
	agt.TaskIdentifier = uuid.NewString()
	agt.Log("Initializing task")
	usersConfigFile, err := utils.OpenConfigFile("./config/accountgenconfig.csv")
	if err != nil {
		return err
	}
	splitConfig := strings.Split(usersConfigFile, "\n")
	if len(splitConfig) < 2 {
		return fmt.Errorf("config not properly populated")
	}
	splitConfig = strings.Split(splitConfig[1], ",")
	if len(splitConfig) != 9 {
		return fmt.Errorf("config not properly populated")
	}
	agt.Config.CatchAll = splitConfig[0]
	agt.Config.CCNumber = splitConfig[1]
	agt.Config.ExpYear = splitConfig[2]
	agt.Config.ExpMonth = splitConfig[3]
	agt.Config.CVV = splitConfig[4]
	agt.Config.ZipCode = splitConfig[5]
	agt.Config.PhoneNumber = splitConfig[6]
	agt.Config.CaptchaKey = splitConfig[7]
	agt.Config.Quantity = splitConfig[8]

	agt.LdrProxy()

	jsonConfigFile, err := utils.OpenConfigJson("./config/config.json")
	if err != nil {
		return err
	}
	agt.Config.Webhook = jsonConfigFile.WebhookUrl
	return nil
}

// will solve recaptcha and return the token
func (agt *AccGenTask) SolveCaptcha() string {
	for {
		agt.Log("Solving captcha")
		captchaKey, err := utils.SolveCaptcha(agt.Config.CaptchaKey)
		if err != nil {
			agt.Log(fmt.Sprintf("Error solving captcha: %s", err.Error()))
			time.Sleep(time.Second * 1)
			continue
		}
		return captchaKey
	}
}

// saves the account data to resy-sniper/config/accounts.csv
// first name,last name,email,password,phone number,token,payment id
func (agt *AccGenTask) SaveAccountData() error {
	//first name,last name,email,password,phone number,token,payment id
	phoneNumberFix := strings.ReplaceAll(agt.Config.PhoneNumber, "%2B", "+")

	stringOut := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s\n",
		agt.FirstName, agt.LastName, agt.EmailAddress,
		agt.Password, phoneNumberFix, agt.AccountData.AccountData.RefreshToken,
		agt.StripeData.PaymentMethodID,
	)
	err := utils.AppendFileSync("./config/accounts.csv", stringOut)
	if err != nil {
		return err
	}
	return nil
}

func (agt *AccGenTask) GetStripePkToken() error {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.resy.com/2/config",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":                {"application/json, text/plain, */*"},
		"accept-language":       {"en-US,en;q=0.9"},
		"authorization":         {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""},
		"cache-control":         {"no-cache"},
		"sec-ch-ua":             {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":      {"?0"},
		"sec-ch-ua-platform":    {"\"Windows\""},
		"sec-fetch-dest":        {"empty"},
		"sec-fetch-mode":        {"cors"},
		"sec-fetch-site":        {"same-site"},
		"x-origin":              {"https://resy.com"},
		"x-resy-auth-token":     {agt.AccountData.AccountData.RefreshToken},
		"x-resy-universal-auth": {agt.AccountData.AccountData.RefreshToken},
		"Referer":               {"https://resy.com/"},
		"Referrer-Policy":       {"strict-origin-when-cross-origin"},
		"user-agent":            {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
	}

	resp, err := agt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		resByteArr, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var resOut SiteConfigResponse

		err = json.Unmarshal(resByteArr, &resOut)
		if err != nil {
			err = fmt.Errorf("error unmarshaling site config data: %s", err.Error())
			return err
		}
		agt.StripeData.StripePkToken = resOut.StripePkToken
		return nil
	}

	return fmt.Errorf("non 200 response received: %d", resp.StatusCode)
}

func (agt *AccGenTask) CreateStripeIntent() error {
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.resy.com/3/stripe/setup_intent",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"authority":       {"api.resy.com"},
		"accept":          {"application/json, text/plain, */*"},
		"accept-encoding": {"gzip, deflate, br, zstd"},
		"accept-language": {"en-US,en;q=0.9"},
		"authorization":   {"ResyAPI api_key=\"VbWk7s3L4KiK5fzlO7JD3Q5EYolJI7n5\""},
		"cache-control":   {"no-cache"},
		"origin":          {"https://resy.com"},
		"pragma":          {"no-cache"},
		"priority":        {"u=1, i"},
		"referer":         {"https://resy.com/"},
		// "sec-ch-ua":             {"\"Google Chrome\";v=\"120\", \"Chromium\";v=\"120\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"Windows\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		// "user-agent":            {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-origin":              {"https://resy.com"},
		"x-resy-auth-token":     {agt.AccountData.AccountData.Token},
		"x-resy-universal-auth": {agt.AccountData.AccountData.Token},
	}

	resp, err := agt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// success is 201 (created) status
	if resp.StatusCode == 201 {
		resByteArr, err := io.ReadAll(resp.Body)
		if err != nil {
			agt.Log(err.Error())
			return nil
		}

		var resOut SetupIntentResponse
		err = json.Unmarshal(resByteArr, &resOut)
		if err != nil {
			agt.Log(fmt.Sprintf("Error unmarshaling stripe setup intent response: %s", err.Error()))
			return nil
		}
		agt.StripeData.ClientSecret = resOut.ClientSecret
		agt.StripeData.CustomerId = resOut.CustomerId
		agt.StripeData.EphemeralKey = resOut.EphemeralKey
		return nil
	}

	return fmt.Errorf("non 201 status received: %d", resp.StatusCode)
}

func (agt *AccGenTask) GetStripeUUIDs() error {
	req, err := http.NewRequest(
		http.MethodPost,
		"https://m.stripe.com/6",
		// this buffer string is constant but response is not
		bytes.NewBufferString("JTdCJTIydjIlMjIlM0ExJTJDJTIyaWQlMjIlM0ElMjIzNjQ4YmE4YzVlNGRhNzk5ZmNkNWY5NjMxMWUxOGZlZSUyMiUyQyUyMnQlMjIlM0E0LjUlMkMlMjJ0YWclMjIlM0ElMjI0LjUuNDMlMjIlMkMlMjJzcmMlMjIlM0ElMjJqcyUyMiUyQyUyMmElMjIlM0FudWxsJTJDJTIyYiUyMiUzQSU3QiUyMmElMjIlM0ElMjIlMjIlMkMlMjJiJTIyJTNBJTIyaHR0cHMlM0ElMkYlMkY1SnpobElKVVFzMEVMYzl2RUd3QnhFNldSR1M0RDFXdVdPNW1OaFoyUTNrLmcydTktaHFadkdJcVlKY1BsUGZ3SkFmLXYzUmd5S194MU5wcHpBbEExMk0lMkZyTDVsX0h2Z2lzbm9MaXdOZ0JNUElkQ0pDa2ZLV0NZT1MwLTNxSVJjcU9RJTJGNE40WFJBRWM0akQtMUVDX2toM01vRElGandMekN2WFFOcDJmeG1XT1dpdyUyMiUyQyUyMmMlMjIlM0ElMjJZM2oxN3NCRFkycjdRSXZDLUZKZzdSNDlna2k2cUEyZU1qdXU3Uy1oWjBBJTIyJTJDJTIyZCUyMiUzQSUyMjU0NzM0MWVmLTMzZWQtNGJhZS1hMTIwLTRlOWU1YmQzYjU1MThkOTI3MCUyMiUyQyUyMmUlMjIlM0ElMjI3ODQwOTY5YS1jOTg1LTQ5ZmUtOGNjYi0xZWIyNTY4MTExNDliNDUxMDElMjIlMkMlMjJmJTIyJTNBZmFsc2UlMkMlMjJnJTIyJTNBdHJ1ZSUyQyUyMmglMjIlM0F0cnVlJTJDJTIyaSUyMiUzQSU1QiUyMmxvY2F0aW9uJTIyJTVEJTJDJTIyaiUyMiUzQSU1QiU1RCUyQyUyMm4lMjIlM0EzMzIuMTk5OTk5OTg4MDc5MDclMkMlMjJ1JTIyJTNBJTIycmVzeS5jb20lMjIlMkMlMjJ3JTIyJTNBJTIyMTcxNzYxNzU4NzE3MCUzQTQzMWQ0YzE5YzVlYTM1YTllZTAwMzA5NmZkYjY4Mjk0MThlYjkzNjZiYzFkN2FiNWQ2ZDhlZTQ1YTg1M2NhM2QlMjIlN0QlMkMlMjJoJTIyJTNBJTIyZTkyYWRiNWFiMDk0MjkzMGRlOWMlMjIlN0Q="),
	)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"accept":             {"*/*"},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"cache-control":      {"no-cache"},
		"content-type":       {"text/plain;charset=UTF-8"},
		"origin":             {"https://m.stripe.network"},
		"pragma":             {"no-cache"},
		"priority":           {"u=1, i"},
		"referer":            {"https://m.stripe.network/"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"120\", \"Chromium\";v=\"120\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"Windows\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
	}

	resp, err := agt.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		resByteArr, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var resOut StripeUUIDsResponse
		err = json.Unmarshal(resByteArr, &resOut)
		if err != nil {
			return err
		}
		agt.StripeData.Guid = resOut.Guid
		agt.StripeData.Muid = resOut.Muid
		agt.StripeData.Sid = resOut.Sid
		return nil
	}

	return fmt.Errorf("non 200 status received: %d", resp.StatusCode)
}

func (agt *AccGenTask) SetStripeIntents() error {

	clientSecretSplit := strings.Split(agt.StripeData.ClientSecret, "_")
	url := fmt.Sprintf("https://api.stripe.com/v1/setup_intents/seti_%s/confirm", clientSecretSplit[1])

	payload := fmt.Sprintf("return_url=https://resy.com/account/payment-methods&payment_method_data[type]=card&payment_method_data[card][number]=%s&payment_method_data[card][cvc]=%s&payment_method_data[card][exp_year]=%s&payment_method_data[card][exp_month]=%s&payment_method_data[allow_redisplay]=unspecified&payment_method_data[billing_details][address][postal_code]=%s&payment_method_data[billing_details][address][country]=US&payment_method_data[pasted_fields]=number&payment_method_data[payment_user_agent]=stripe.js/417cd13f1a;+stripe-js-v3/417cd13f1a;+payment-element&payment_method_data[referrer]=https://resy.com&payment_method_data[time_on_page]=8503935&payment_method_data[client_attribution_metadata][client_session_id]=%s&payment_method_data[client_attribution_metadata][merchant_integration_source]=elements&payment_method_data[client_attribution_metadata][merchant_integration_subtype]=payment-element&payment_method_data[client_attribution_metadata][merchant_integration_version]=2021&payment_method_data[client_attribution_metadata][payment_intent_creation_flow]=standard&payment_method_data[client_attribution_metadata][payment_method_selection_flow]=merchant_specified&payment_method_data[guid]=%s&payment_method_data[muid]=%s&payment_method_data[sid]=%s&expected_payment_method_type=card&radar_options[hcaptcha_token]=%s&use_stripe_sdk=true&key=%s&client_secret=%s",
		agt.Config.CCNumber, agt.Config.CVV,
		agt.Config.ExpYear, agt.Config.ExpMonth,
		agt.Config.ZipCode, uuid.NewString(),
		agt.StripeData.Guid, agt.StripeData.Muid,
		agt.StripeData.Sid, "s",
		agt.StripeData.StripePkToken, agt.StripeData.ClientSecret,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBufferString(payload),
	)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":             {"application/json"},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"en-US,en;q=0.9"},
		"cache-control":      {"no-cache"},
		"content-type":       {"application/x-www-form-urlencoded"},
		"origin":             {"https://js.stripe.com"},
		"pragma":             {"no-cache"},
		"priority":           {"u=1,i"},
		"referer":            {"https://js.stripe.com"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"120\", \"Chromium\";v=\"120\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"Windows\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-site"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
	}

	resp, err := agt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		resByteArr, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var resOut SetStripeIntents
		err = json.Unmarshal(resByteArr, &resOut)
		if err != nil {
			return err
		}
		agt.StripeData.PaymentMethodID = resOut.PaymentMethodID
		return nil
	}

	return fmt.Errorf("non 200 status received: %d", resp.StatusCode)
}

func (agt *AccGenTask) SetPaymentMethod() error {
	payload := fmt.Sprintf("stripe_payment_method_id=%s", agt.StripeData.PaymentMethodID)
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.resy.com/3/stripe/payment_method",
		bytes.NewBufferString(payload),
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
		"content-type":          {"application/x-www-form-urlencoded"},
		"origin":                {"https://resy.com"},
		"pragma":                {"no-cache"},
		"priority":              {"u=1,i"},
		"referer":               {"https://resy.com/"},
		"sec-ch-ua":             {"\"Google Chrome\";v=\"120\", \"Chromium\";v=\"120\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":      {"?0"},
		"sec-ch-ua-platform":    {"\"Windows\""},
		"sec-fetch-dest":        {"empty"},
		"sec-fetch-mode":        {"cors"},
		"sec-fetch-site":        {"same-site"},
		"user-agent":            {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"x-origin":              {"https://resy.com"},
		"x-resy-auth-token":     {agt.AccountData.AccountData.Token},
		"x-resy-universal-auth": {agt.AccountData.AccountData.Token},
	}

	resp, err := agt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// success is 201 created
	if resp.StatusCode == 201 {
		resBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var resOut PaymentMethodList
		json.Unmarshal(resBytes, &resOut)
		agt.StripeData.PaymentMethodID = strconv.Itoa(resOut.PaymentMethod[0].ID)
		return nil
	}

	return fmt.Errorf("non 200 status received: %d", resp.StatusCode)
}
