package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
	"golang.org/x/crypto/argon2"
)

type cryptoService struct {
	encryptionKey    []byte
	newCryptoStorage interfaces.NewCryptoStorage
}

func NewCryptoService(newCryptoStorage interfaces.NewCryptoStorage) interfaces.CryptoService {
	return &cryptoService{newCryptoStorage: newCryptoStorage}
}

func (s *cryptoService) InitCrypto(userID, login, password string) (interfaces.Storage, error) {
	s.encryptionKey = s.generateKey(userID, login, password)
	return s.newCryptoStorage(userID, s.encryptionKey)
}

func (s *cryptoService) generateKey(userID, login, password string) []byte {
	salt := []byte(login + userID)
	return argon2.IDKey([]byte(password), salt, 2, 128*1024, 4, 32)
}

func (s *cryptoService) EncryptRecord(record *models.Record) error {
	aesGCM, err := s.newGCM()
	if err != nil {
		return err
	}
	record.Nonce = make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, record.Nonce); err != nil {
		return err
	}
	record.Data = aesGCM.Seal(nil, record.Nonce, record.Data, nil)
	return nil
}

func (s *cryptoService) DecryptRecord(record *models.Record) error {
	aesGCM, err := s.newGCM()
	if err != nil {
		return err
	}
	record.Data, err = aesGCM.Open(nil, record.Nonce, record.Data, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *cryptoService) newGCM() (cipher.AEAD, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
