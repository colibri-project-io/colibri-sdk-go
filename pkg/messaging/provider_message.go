package messaging

import (
	"bytes"
	"encoding/json"

	"github.com/google/uuid"
)

type ProviderMessage struct {
	Id       uuid.UUID   `json:"id"`
	Origin   string      `json:"origin"`
	Action   string      `json:"action"`
	TenantId uuid.UUID   `json:"tenantId"`
	UserId   uuid.UUID   `json:"userId"`
	Message  interface{} `json:"message"`
	n        interface{}
}

// String convert struct into json string
func (msg *ProviderMessage) String() string {
	message, _ := json.Marshal(msg)

	return string(message)
}

// DecodeMessage transform interface into ProviderMessage
func (msg *ProviderMessage) DecodeMessage(model interface{}) error {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(msg.Message); err != nil {
		return err
	}

	if err := json.NewDecoder(buf).Decode(model); err != nil {
		return err
	}

	return nil
}

// addOriginBrokerNotification add reference of origin broker message to send dlq if an error occurs
func (msg *ProviderMessage) addOriginBrokerNotification(n interface{}) {
	msg.n = n
}
