package utils

import (
	"github.com/google/uuid"
	"log"
)

func ParseStringToUUID(id string) uuid.UUID {
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
	}
	return uid
}
