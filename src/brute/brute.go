package brute

import (
	"fmt"
	"time"
)

// main brute task loop for booking reservations
func (task *ResyTask) Start() {
	// main task loop that needs to be repeated
	task.TaskStartTime = time.Now()
	for {
		select {
		case <-task.StopChannel:
			// task.LogError(fmt.Sprintf("Stopping task for %s", task.EmailAddress))
			return
		default:
			//
		}
		task.Log("Starting task")
		err := task.InitHttpClient()
		if err != nil {
			task.LogError("error getting http client")
		}

		// err = task.Login()
		// if err != nil {
		// 	task.LogError(err.Error())
		// }

		// return
		// if err != nil {
		// 	task.LogError(fmt.Sprintf("Error initalizing http client: %s", err.Error()))
		// 	time.Sleep(time.Millisecond * time.Duration(task.Delay))
		// 	continue
		// }

		if len(task.PaymentId) == 0 {
			task.Log("Getting payment ID")
			err = task.GetPaymentId()
			if err != nil {
				task.LogError(err.Error())
				time.Sleep(time.Millisecond * time.Duration(task.Delay))
				continue
			}
		}

		task.Log("Getting reservations")
		err = task.GetRandomReservation()
		if err != nil {
			task.LogError(fmt.Sprintf("Error getting reservation slots: %s", err.Error()))
			time.Sleep(time.Millisecond * time.Duration(task.Delay))
			continue
		}

		task.Log("Getting booking token")
		err = task.GetBookingToken()
		if err != nil {
			task.LogError(fmt.Sprintf("Error getting booking token: %s", err.Error()))
			time.Sleep(time.Millisecond * time.Duration(task.Delay))
			continue
		}

		task.Log("Posting reservation")
		// TODO check for status 412 Precondition Failed
		// means reservation was already booked and to tell
		// user to check email/acc to see if it already has been booked by the bot
		err = task.PostReservation()
		if err != nil {
			task.LogError(fmt.Sprintf("Error posting reservation: %s", err.Error()))
			time.Sleep(time.Millisecond * time.Duration(task.Delay))
			continue
		}

		task.LogSuccess(fmt.Sprintf("Successfully booked %s in %s", task.Venue.Name, time.Since(task.TaskStartTime)))
		if len(task.WebhookUrl) > 10 {
			// ignore error if any dont care if it goes through
			task.SendSuccessWebhook()
		}
		// kill and break make sure task is exitted completely
		task.Kill()
		break
	}
}
