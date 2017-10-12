package main

import (
	"ekyu.moe/soda/i18n"

	"crypto/subtle"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"strings"

	"ekyu.moe/base91"
)

func pubValidator(val interface{}) error {
	if str, ok := val.(string); ok {
		decoded := base91.DecodeString(str)
		// 校验长度，应包含 CRC32（长度 4）和密钥本体（长度 32）
		if len(decoded) == 36 {
			actual := crc32.Checksum(decoded[4:], crc32.IEEETable)
			expected := binary.BigEndian.Uint32(decoded[:4])
			// 校验 CRC32
			if subtle.ConstantTimeEq(int32(actual), int32(expected)) == 1 {
				return nil
			}
		}
	}

	return errors.New(i18n.INVALID_PUB)
}

func plainValidator(val interface{}) error {
	if str, ok := val.(string); ok && len(strings.TrimSpace(str)) > 0 {
		return nil
	}

	return errors.New(i18n.INVALID_PLAIN)
}

func encryptedValidator(val interface{}) error {
	if str, ok := val.(string); ok {
		decoded := base91.DecodeString(str)
		// 校验长度，应包含 CRC32（长度 4）和 nonce（长度 24）
		if len(decoded) > 28 {
			actual := crc32.Checksum(decoded[4:], crc32.IEEETable)
			expected := binary.BigEndian.Uint32(decoded[:4])
			// 校验 crc32
			if subtle.ConstantTimeEq(int32(actual), int32(expected)) == 1 {
				return nil
			}
		}
	}

	return errors.New(i18n.INVALID_ENCRYPTED)
}
