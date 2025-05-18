package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type Encryption struct {
	Key []byte
}

// Шифруем текст с помощью симметричного алгоритма AES (в режиме CFB) и возвращаем результат в виде base64-строки
func (e *Encryption) Encrypt(text string) (string, error) {
	// Преобразовываем полученный текст в срез байт
	plaintext := []byte(text)

	// создаем новыц шифр-блок AES с ключом e.Key
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	// Создаем зашифрованный текст
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// Шифрование данных
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	// Преобразовываем зашифрованный текст в base64
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// Расшифровка зашифрованного текста
func (e *Encryption) Decrypt(cryptoText string) (string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText), nil
}
