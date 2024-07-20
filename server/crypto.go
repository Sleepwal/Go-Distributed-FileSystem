package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func generateID() string {
	buf := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, buf) // 随机数生成器
	return hex.EncodeToString(buf)       // 将字节转换为16进制字符串
}

func hashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func newEncryptionKey() []byte {
	buf := make([]byte, 32)
	_, _ = io.ReadFull(rand.Reader, buf)
	return buf
}

// copyStream
// @Description: 用于将一个加密流（cipher.Stream）从一个输入源（io.Reader）复制到输出目标
// @param stream
// @param blockSize
// @param src 一个输入源对象，用于读取需要加密/解密的数据。
// @param dst 一个输出目标对象，用于写入加密后的数据。
// @return int
// @return error
func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {
	var (
		buf = make([]byte, 32*1024)
		nw  = blockSize
	)

	for {
		n, err := src.Read(buf)
		if n > 0 {
			// 将一个字节切片（buf）与另一个字节切片（buf[:n]）进行异或操作，从而实现加密或解密的效果。
			stream.XORKeyStream(buf, buf[:n])
			nn, err := dst.Write(buf[:n]) // 将加密数据写入目标文件
			if err != nil {
				return 0, err
			}
			nw += nn
		}

		if err == io.EOF { // 读取到文件末尾
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return nw, nil
}

// copyDecrypt
// @Description: 解密数据并复制
// @param key
// @param src
// @param dst
// @return int
// @return error
func copyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	//Create a new AES cipher using the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())
	//Read the iv from the src
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	//Create a new CTR stream with the block and iv
	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

// copyEncrypt
// @Description: 加密数据并复制
// @param key
// @param src
// @param dst
// @return int
// @return error
func copyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}
