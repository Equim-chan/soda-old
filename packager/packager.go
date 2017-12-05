package packager // import "ekyu.moe/soda/packager"

import (
	"bytes"
	"compress/zlib"
	"crypto/subtle"
	"encoding/binary"
	"hash/crc32"
	"io/ioutil"

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

// 该函数会对输入尝试进行压缩，如果体积有所减小，则返回压缩后的内容和 true。
// 反之，返回原 Buffer 和 false。
// This function is never thread safe.
// TODO: assert orig immute, destroy
func Pack(orig *memguard.LockedBuffer) (*memguard.LockedBuffer, bool) {
	if err := globalBuffer.MakeMutable(); err != nil {
		return orig, false
	}
	defer globalBuffer.MakeImmutable()
	defer globalBuffer.Wipe()

	// 尝试压缩
	buf := bytes.NewBuffer(globalBuffer.Buffer()[:0])
	z, err := zlib.NewWriterLevel(buf, zlib.BestCompression)
	// 遇到错误直接选择不压缩
	if err != nil {
		return orig, false
	}
	defer z.Close()

	if _, err := z.Write(orig.Buffer()); err != nil {
		return orig, false
	}
	if err := z.Close(); err != nil {
		return orig, false
	}

	if l := buf.Len(); l < orig.Size() {
		ret, err := memguard.Trim(globalBuffer, 0, l)
		if err != nil {
			return orig, false
		}

		if err := ret.MakeImmutable(); err != nil {
			return orig, false
		}

		return ret, true
	}

	return orig, false
}

// 该函数会对输入尝试进行解压缩，如果是合法的 zlib 压缩，则返回解压后的内容和 nil
// 反之，返回 nil 和 error。
func Unpack(packet *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	// 尝试解压
	r := bytes.NewReader(packet.Buffer())
	z, err := zlib.NewReader(r)
	// 遇到错误直接不解压
	// 但是这样有隐藏的风险，也许只是偶然的错误，但内容确实是压缩过的
	if err != nil {
		return nil, err
	}
	defer z.Close()

	// 这没办法了
	decompressed, err := ioutil.ReadAll(z)
	if err != nil {
		return nil, err
	}
	if err := z.Close(); err != nil {
		return nil, err
	}

	return memguard.NewImmutableFromBytes(decompressed)
}
