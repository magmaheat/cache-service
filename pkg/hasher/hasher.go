package hasher

import "golang.org/x/crypto/bcrypt"

type HashManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) bool
}

type Crypto struct{}

func NewCryptoHasher() *Crypto {
	return &Crypto{}
}

func (c *Crypto) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (c *Crypto) CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
