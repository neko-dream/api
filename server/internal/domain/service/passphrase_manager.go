package service

import (
	"github.com/neko-dream/api/internal/domain/model/talksession"
	"github.com/neko-dream/api/internal/infrastructure/config"
	"github.com/neko-dream/api/pkg/hash"
	"golang.org/x/crypto/bcrypt"
)

type passphraseManager struct {
	cfg *config.Config
}

func NewPassphraseManager(cfg *config.Config) talksession.PassphraseManager {
	return &passphraseManager{cfg: cfg}
}

func (p *passphraseManager) Hash(plaintext string) (string, error) {
	return hash.BcryptHashHMAC(plaintext, p.cfg.HASH_PEPPER, bcrypt.DefaultCost)
}

func (p *passphraseManager) Verify(plaintext, hashedPassphrase string) bool {
	return hash.BcryptVerifyHMAC(plaintext, p.cfg.HASH_PEPPER, hashedPassphrase)
}
