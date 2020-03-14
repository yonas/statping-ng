// Statping
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/statping/statping
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package notifiers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/statping/statping/types/failures"
	"github.com/statping/statping/types/notifications"
	"github.com/statping/statping/types/notifier"
	"github.com/statping/statping/types/services"
	"github.com/statping/statping/utils"
	"net/url"
	"strings"
	"time"
)

var _ notifier.Notifier = (*twilio)(nil)

type twilio struct {
	*notifications.Notification
}

func (t *twilio) Select() *notifications.Notification {
	return t.Notification
}

var Twilio = &twilio{&notifications.Notification{
	Method:      "twilio",
	Title:       "Twilio",
	Description: "Receive SMS text messages directly to your cellphone when a service is offline. You can use a Twilio test account with limits. This notifier uses the <a href=\"https://www.twilio.com/docs/usage/api\">Twilio API</a>.",
	Author:      "Hunter Long",
	AuthorUrl:   "https://github.com/hunterlong",
	Icon:        "far fa-comment-alt",
	Delay:       time.Duration(10 * time.Second),
	Form: []notifications.NotificationForm{{
		Type:        "text",
		Title:       "Account SID",
		Placeholder: "Insert your Twilio Account SID",
		DbField:     "api_key",
		Required:    true,
	}, {
		Type:        "text",
		Title:       "Account Token",
		Placeholder: "Insert your Twilio Account Token",
		DbField:     "api_secret",
		Required:    true,
	}, {
		Type:        "text",
		Title:       "SMS to Phone Number",
		Placeholder: "18555555555",
		DbField:     "Var1",
		Required:    true,
	}, {
		Type:        "text",
		Title:       "From Phone Number",
		Placeholder: "18555555555",
		DbField:     "Var2",
		Required:    true,
	}}},
}

// Send will send a HTTP Post to the Twilio SMS API. It accepts type: string
func (u *twilio) sendMessage(message string) error {
	twilioUrl := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%v/Messages.json", u.GetValue("api_key"))

	v := url.Values{}
	v.Set("To", "+"+u.Var1)
	v.Set("From", "+"+u.Var2)
	v.Set("Body", message)
	rb := *strings.NewReader(v.Encode())

	contents, _, err := utils.HttpRequest(twilioUrl, "POST", "application/x-www-form-urlencoded", nil, &rb, time.Duration(10*time.Second), true)
	success, _ := twilioSuccess(contents)
	if !success {
		errorOut := twilioError(contents)
		out := fmt.Sprintf("Error code %v - %v", errorOut.Code, errorOut.Message)
		return errors.New(out)
	}
	return err
}

// OnFailure will trigger failing service
func (u *twilio) OnFailure(s *services.Service, f *failures.Failure) error {
	msg := fmt.Sprintf("Your service '%v' is currently offline!", s.Name)
	return u.sendMessage(msg)
}

// OnSuccess will trigger successful service
func (u *twilio) OnSuccess(s *services.Service) error {
	msg := fmt.Sprintf("Your service '%v' is currently online!", s.Name)
	return u.sendMessage(msg)
}

// OnTest will test the Twilio SMS messaging
func (u *twilio) OnTest() error {
	msg := fmt.Sprintf("Testing the Twilio SMS Notifier")
	return u.sendMessage(msg)
}

func twilioSuccess(res []byte) (bool, twilioResponse) {
	var obj twilioResponse
	json.Unmarshal(res, &obj)
	if obj.Status == "queued" {
		return true, obj
	}
	return false, obj
}

func twilioError(res []byte) twilioErrorObj {
	var obj twilioErrorObj
	json.Unmarshal(res, &obj)
	return obj
}

type twilioErrorObj struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   int    `json:"status"`
}

type twilioResponse struct {
	Sid                 string      `json:"sid"`
	DateCreated         string      `json:"date_created"`
	DateUpdated         string      `json:"date_updated"`
	DateSent            interface{} `json:"date_sent"`
	AccountSid          string      `json:"account_sid"`
	To                  string      `json:"to"`
	From                string      `json:"from"`
	MessagingServiceSid interface{} `json:"messaging_service_sid"`
	Body                string      `json:"body"`
	Status              string      `json:"status"`
	NumSegments         string      `json:"num_segments"`
	NumMedia            string      `json:"num_media"`
	Direction           string      `json:"direction"`
	APIVersion          string      `json:"api_version"`
	Price               interface{} `json:"price"`
	PriceUnit           string      `json:"price_unit"`
	ErrorCode           interface{} `json:"error_code"`
	ErrorMessage        interface{} `json:"error_message"`
	URI                 string      `json:"uri"`
	SubresourceUris     struct {
		Media string `json:"media"`
	} `json:"subresource_uris"`
}
