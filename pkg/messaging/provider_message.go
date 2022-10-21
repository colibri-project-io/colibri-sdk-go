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
}

func (msg ProviderMessage) String() string {
	message, _ := json.Marshal(msg)

	return string(message)
}

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
