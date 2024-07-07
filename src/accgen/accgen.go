package accgen

import (
	"fmt"
	"time"
)

func Start() {
	var task AccGenTask
	// sets faker values here at InitAccountGenConfig()
	err := task.InitAccountGenConfig()
	if err != nil {
		task.LogError(fmt.Sprintf("Error loading account generation configuration: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}
	// validate config if invalid exit task
	task.Log("Validating configuration")
	err = task.ValidateAccGen()
	if err != nil {
		task.LogError(fmt.Sprintf("Configuration invalid, %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}
	// horrible code to change blank text to localhost if user is
	// not using a proxy
	var proxyText string
	if len(task.Proxy.ProxyHttpFormat) == 0 {
		proxyText = "localhost"
	} else {
		proxyText = task.Proxy.ProxyHttpFormat
	}
	task.Log(fmt.Sprintf("Starting task with: email:%s password:%s proxy:%s",
		task.EmailAddress, task.Password, proxyText))

	captchaKey := task.SolveCaptcha()
	task.Log("Captcha successfully solved")
	payload := fmt.Sprintf("first_name=%s&last_name=%s&mobile_number=%s&em_address=%s&policies_accept=1&marketing_opt_in=0&complete=1&device_type_id=3&device_token=%s&isNonUS=0&password=%s&captcha_token=%s",
		task.FirstName, task.LastName, task.Config.PhoneNumber, task.EmailAddress, task.TaskIdentifier, task.Password, captchaKey)
	// initialize tls client that will be used for all subsequent requests
	task.InitHttpClient()
	task.Log("Posting registration")
	err = task.PostReigstration(payload)
	if err != nil {
		task.LogError(fmt.Sprintf("Error posting registration: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

	task.LogSuccess("Successfully registered account")
	// if no cc is set then exit task and save account
	if len(task.Config.CCNumber) < 5 {
		task.SaveAccountData()
		err = task.SendCreatedWebhook(false)
		task.LogError(fmt.Sprintf("Error sending success webhook: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

	task.Log("Getting stripe pk token")
	err = task.GetStripePkToken()
	if err != nil {
		task.LogError(fmt.Sprintf("Error getting pk token: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

	task.Log("Creating stripe intent")
	err = task.CreateStripeIntent()
	if err != nil {
		task.LogError(fmt.Sprintf("Error creating stripe intent: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

	task.Log("Getting stripe UUIDs")
	err = task.GetStripeUUIDs()
	if err != nil {
		task.LogError(fmt.Sprintf("Error getting stripe UUIDs: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

	task.Log("Setting stripe intents")
	err = task.SetStripeIntents()
	if err != nil {
		task.LogError(fmt.Sprintf("Error setting stripe intents: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

	task.Log("Setting payment method")
	err = task.SetPaymentMethod()
	if err != nil {
		task.LogError(fmt.Sprintf("Error setting payment method: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

	task.LogSuccess("Successfully registered and set payment method")
	err = task.SaveAccountData()
	if err != nil {
		task.LogError(fmt.Sprintf("Error saving account data: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}
	err = task.SendCreatedWebhook(true)
	if err != nil {
		task.LogError(fmt.Sprintf("Error sending success webhook: %s", err.Error()))
		time.Sleep(time.Second * 3)
		return
	}

}
