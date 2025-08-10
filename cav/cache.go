package cav

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type cachedSessions struct {
	Sessions map[string]map[string]string `json:"sessions"`
}

func (c *client) storeSessionsToCache(passphrase, path string) error {
	cs := cachedSessions{
		Sessions: map[string]map[string]string{},
	}

	for _, value := range c.clientsInitialized {
		cs.Sessions[value.getID()] = value.getCredential().getSession()
	}

	csJSON, err := json.Marshal(cs)
	if err != nil {
		return err
	}

	csEncrypted, err := encryptSessions([]byte(passphrase), string(csJSON))
	if err != nil {
		return err
	}

	return writeGobFile(path, csEncrypted)
}

func (c *client) restoreSessionsFromCache(passphrase, path string) error {
	c.cachePassphrase = passphrase
	c.cachePath = path

	// check if cache file exist. If not exist ignore
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.logger.Warn("Cache file does not exist, skipping restoration", "path", path)
		return nil
	}

	var csEncrypted string

	err := readGobFile(path, &csEncrypted)
	if err != nil {
		return err
	}

	csJSON, err := decryptSessions([]byte(passphrase), csEncrypted)
	if err != nil {
		return err
	}

	var cs cachedSessions
	if err := json.Unmarshal([]byte(csJSON), &cs); err != nil {
		return err
	}

	for id, session := range cs.Sessions {
		for _, value := range c.clientsInitialized {
			if value.getID() == id {
				if err := value.getCredential().restoreSession(session); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func encryptSessions(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decryptSessions(key []byte, ciphertext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext is too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

func writeGobFile(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return encoder.Encode(data)
}

func readGobFile(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	return decoder.Decode(data)
}
