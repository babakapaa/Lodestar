package chat

import "github.com/google/uuid"

func validID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
