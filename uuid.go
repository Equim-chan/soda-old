package main

import (
	"crypto/rand"
	"encoding/hex"
)

func uuidv4() (string, error) {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return "", err
	}

	// Per 4.4, set bits for version and `clock_seq_hi_and_reserved`
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	str := hex.EncodeToString(uuid)
	str = str[:8] + "-" + str[8:12] + "-" + str[12:16] + "-" + str[16:20] + "-" + str[20:]

	return str, nil
}
