package main

import (
	"ekyu.moe/soda/i18n"

	"bytes"
	"compress/zlib"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"io/ioutil"

	"ekyu.moe/base91"
	"golang.org/x/crypto/nacl/box"
)

// 该函数会对输入尝试进行压缩，如果体积有所减小，则返回压缩后的内容和 true。
// 反之，返回一个空 []byte 和 false。
func tryCompess(content []byte) ([]byte, bool) {
	// 尝试压缩
	buf := new(bytes.Buffer)
	z, err := zlib.NewWriterLevel(buf, zlib.BestCompression)
	// 遇到错误直接选择不压缩
	if err != nil {
		return content, false
	}
	defer z.Close()

	if _, err := z.Write(content); err != nil {
		return content, false
	}
	if err := z.Close(); err != nil {
		return content, false
	}

	compressed := buf.Bytes()
	if len(compressed) < len(content) {
		return compressed, true
	}

	return content, false
}

// 打包明文，返回 base91 字符串。
// 顺序为：压缩或不压缩明文 -> 生成 nonce -> 加密 -> 写入 crc32 头 -> base91 encode
func pack(plain string) (string, error) {
	// 预处理（选择性压缩）
	content, _ := tryCompess([]byte(plain))
	// log.Println("compress:", ok)

	// 生成 nonce
	// nonce 和 shared 都是全局变量
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", err
	}

	// 加密
	encrypted := box.SealAfterPrecomputation(nonce[:], content, nonce, shared)

	// 计算并写入 crc32 头
	payload := make([]byte, 4+len(encrypted))

	c := crc32.Checksum(encrypted, crc32.IEEETable)
	binary.BigEndian.PutUint32(payload[:4], c)

	copy(payload[4:], encrypted)

	// base91 encode
	encoded := base91.EncodeToString(payload)

	return encoded, nil
}

// 该函数会对输入尝试进行解压缩，如果是合法的 zlib 压缩，则返回解压后的内容和 true
// 反之，返回一个空 []byte 和 false
func tryDecompess(content []byte) ([]byte, bool) {
	// 尝试解压
	r := bytes.NewReader(content)
	z, err := zlib.NewReader(r)
	// 遇到错误直接不解压
	// 但是这样有隐藏的风险，也许只是偶然的错误，但内容确实是压缩过的
	if err != nil {
		return content, false
	}
	defer z.Close()

	decompressed, err := ioutil.ReadAll(z)
	if err != nil {
		return content, false
	}
	if err := z.Close(); err != nil {
		return content, false
	}

	return decompressed, true
}

// 从密文中解包明文。
// 由于 prompt.go 和 validate.go 已经分别接管了 base91 decode 和 crc32 头校验，
// 这里仅仅处理解密和可能的解压。
// 顺序为：剥离 nonce -> 解密 -> 尝试解压
func unpack(encrypted []byte) (string, error) {
	// 获取 nonce
	copy(nonce[:], encrypted[:24])

	// 解密
	plain, ok := box.OpenAfterPrecomputation(nil, encrypted[24:], nonce, shared)
	if !ok {
		return "", errors.New(i18n.DECRYPT_FAIL)
	}

	// 尝试解压
	content, _ := tryDecompess(plain)
	// log.Println("decompress:", ok)

	return string(content), nil
}
