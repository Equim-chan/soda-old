package packager // import "ekyu.moe/soda/packager"

import (
	"bytes"
	"compress/zlib"
	"crypto/subtle"
	"encoding/binary"
	"hash/crc32"
	"io/ioutil"
	"log"

	"github.com/awnumar/memguard"
)

var (
	globalBuffer *memguard.LockedBuffer
)

func init() {
	buf, err := memguard.NewImmutable(16 * 1024)
	if err != nil {
		panic("buffer init: " + err.Error())
	}

	globalBuffer = buf
}

func AttachCrc32(content []byte) []byte {
	ret := make([]byte, len(content)+4)

	sum := crc32.Checksum(content, crc32.IEEETable)
	binary.BigEndian.PutUint32(ret[:4], sum)
	copy(ret[4:], content)

	return ret
}

func DetachCrc32(content []byte) ([]byte, bool) {
	actual := crc32.Checksum(content[4:], crc32.IEEETable)
	expected := binary.BigEndian.Uint32(content[:4])

	if subtle.ConstantTimeEq(int32(actual), int32(expected)) == 1 {
		return content[4:], true
	}

	return nil, false
}

// Pack tries to compress orig, the return value may be compress or not.
// If ever an error is returned, then it must be a fatal one.
func Pack(orig *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	// This is an assertion.
	if orig.IsMutable() {
		panic("packet must be immutable")
	}

	if err := globalBuffer.MakeMutable(); err != nil {
		return nil, err
	}
	defer globalBuffer.MakeImmutable()
	defer globalBuffer.Wipe()

	// 尝试压缩
	buf := bytes.NewBuffer(globalBuffer.Buffer()[:0])
	z, err := zlib.NewWriterLevel(buf, zlib.BestCompression)
	// 遇到错误直接选择不压缩
	if err != nil {
		return orig, nil
	}
	defer z.Close()

	if _, err := z.Write(orig.Buffer()); err != nil {
		return orig, nil
	}
	if err := z.Close(); err != nil {
		return orig, nil
	}

	// Compare the sizes
	l := buf.Len()
	if l >= orig.Size() {
		log.Println("no compress")
		return orig, nil
	}
	log.Println("compress")

	// these errors are all fatal
	ret, err := memguard.Trim(globalBuffer, 0, l)
	if err != nil {
		return nil, err
	}

	if err := ret.MakeImmutable(); err != nil {
		ret.Destroy()
		return nil, err
	}

	return ret, nil
}

// Unpack tries to decompress orig, the return value may be decompress or not.
// If ever an error is returned, then it must be a fatal one.
func Unpack(orig *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	// Assert
	if orig.IsMutable() {
		panic("packet must be immutable")
	}

	// Try decompress
	r := bytes.NewReader(orig.Buffer())
	z, err := zlib.NewReader(r)
	// Once an error is met, return orig.
	// However this is risky, as it is possible to be an occasional error
	// with its content indeed compressed.
	if err != nil {
		return orig, nil
	}
	defer z.Close()

	// I can't come up with a safer idea at the moment
	decompressed, err := ioutil.ReadAll(z)
	if err != nil {
		return orig, nil
	}
	if err := z.Close(); err != nil {
		return orig, nil
	}

	ret, err := memguard.NewImmutableFromBytes(decompressed)
	if err != nil {
		ret.Destroy()
		return nil, err
	}

	return ret, nil
}
