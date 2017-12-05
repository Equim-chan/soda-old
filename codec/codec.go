// package codec provides some binary-to-text encode/decode interfaces
package codec // import "ekyu.moe/soda/codec"

import (
	"regexp"

	"ekyu.moe/base256"
	"ekyu.moe/base91"
	"github.com/kyokomi/emoji"
)

type EncodeFunc func([]byte) string

var (
	Base91Encode   = base91.EncodeToString
	EmojiEncode    = base256.EncodeToString
	EmojiTagEncode = encodeEmojiTag
)

var (
	nonAscii        = regexp.MustCompile(`[[:^ascii:]]`)
	spaces          = regexp.MustCompile(`[[:space:]]`)
	emojiToTagSheet = make(map[string]string)
)

func init() {
	// rip time, rip space
	for i, v := range emoji.CodeMap() {
		emojiToTagSheet[v] = i
	}
}

func encodeEmojiTag(p []byte) string {
	raw := base256.EncodeToString(p)

	ret := ""
	for _, v := range raw {
		ret += emojiToTagSheet[string(v)]
	}

	return ret
}

func DetectCodecAndDecode(s string) []byte {
	filtered := spaces.ReplaceAllString(s, "")
	i := emoji.Sprint(filtered)

	if nonAscii.MatchString(i) {
		return base256.DecodeString(i)
	}

	return base91.DecodeString(i)
}
