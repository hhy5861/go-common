package common

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

type (
	Password struct {
		Password string
		Salt     string
		Mode     string
	}
)

const (
	Md5Mode    = "md5"
	Sha256Mode = "sha256"
	Sha512Mode = "sha512"
)

func NewPassword(pwd *Password) *Password {
	return pwd
}

func (p *Password) GenerateSalt(n int) *Password {
	p.Salt = tools.GetRandomString(n)

	return p
}

func (p *Password) formatPassword() string {
	return fmt.Sprintf("%s:%s", p.Password, p.Salt)
}

func (p *Password) HashMD5() *Password {
	p.Mode = Md5Mode
	h := md5.New()

	h.Write([]byte(p.formatPassword()))
	p.Password = hex.EncodeToString(h.Sum(nil))

	return p
}

func (p *Password) Hash256() *Password {
	p.Mode = Sha256Mode
	h := sha256.New()

	h.Write([]byte(p.formatPassword()))
	p.Password = hex.EncodeToString(h.Sum(nil))

	return p
}

func (p *Password) Hash512() *Password {
	p.Mode = Sha512Mode
	h := sha512.New()

	h.Write([]byte(p.formatPassword()))
	p.Password = hex.EncodeToString(h.Sum(nil))

	return p
}

func (p *Password) EqualsPassword(hashPwd string) bool {
	return strings.EqualFold(p.Password, hashPwd)
}
