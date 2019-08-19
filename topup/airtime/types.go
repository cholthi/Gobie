package airtime

import (
	"encoding/xml"
	"fmt"
	"time"
)

/*Context this	is	a	structure	used	in	all	requests	to
  identify	and	authorize	the	client	performing
  the	transaction
*/
type Context struct {
	Channel  string        `xml:"channel"`
	ID       string        `xml:"clientId"`
	TimeOut  time.Duration `xml:"clientRequestTimeout"`
	User     User          `xml:"initiatorPrincipalId"`
	Password string        `xml:"password"`
	Ref      string        `xml:"clientReference"`
	Tag      string        `xml:"clientTag"`
	Comment  string        `xml:"clientComment"`
}

//User Authenticated and registered API user at the server
type User struct {
	XMLName     xml.Name
	ID          string `xml:"id"`
	Type        string `xml:"type"`
	RequestType string `xml:"userId"`
}

//Account the account info of the party involved in the transaction
type Account struct {
	XMLName xml.Name
	ID      string `xml:"accountId"`
	Type    string `xml:"accountTypeId"`
}

//Amount represents the airtime value in a certain currency
type Amount struct {
	XMLName  xml.Name `xml:"balance"`
	Currency string   `xml:"currency"`
	Value    float64  `xml:"value"`
}

//TopUpRequest represents request sent to the API endpoint[used to marshal to xml]
type TopUpRequest struct {
	XMLName           xml.Name `xml:"ns2:requestTopupResponse"`
	Context           Context  `xml:"context"`
	Sender            User
	SenderAccount     Account
	Subscriber        User
	SubscriberAccount Account
	ProductID         string `xml:"productId"`
	Amount            Amount
}

//TopupResponse represents response from the API server[used by xml.Unmarshal]
type TopupResponse struct {
	XMLName      xml.Name `xml:"requestTopupResponse"`
	Reference    string   `xml:"return>ersReference"`
	ResultCode   int      `xml:"return>ResultCode"`
	Status       string   `xml:"return>resultDescription"`
	SenderMsisdn string   `xml:"return>topupPrincipal>principalId>id"`
	Balance      Amount   `xml:"return>senderPrincipal>accounts>account>balance"`
}
type ErrorResponse interface {
	Message() string
	Code() int
	Error() string
}

// Implements ErrorResponse interface
func (err *TopupResponse) Error() string {

	return fmt.Sprintf("%d:%s", err.ResultCode, err.Status)

}

// Implements ErrorResponse interface
func (err *TopupResponse) Message() string {
	return err.Status
}

// Implements ErrorResponse interface
func (err *TopupResponse) Code() int {
	return err.ResultCode
}
