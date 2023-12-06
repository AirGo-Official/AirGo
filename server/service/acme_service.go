package service

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/providers/dns/godaddy"
	"github.com/go-acme/lego/v4/providers/dns/hetzner"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
	"github.com/go-acme/lego/v4/registration"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"os"
	"strings"
	"time"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func AcmeApply(a *model.Acme) error {
	var provider challenge.Provider
	var err error
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if a.AcmeEmail == "" {
		a.AcmeEmail = fmt.Sprintf("%s@airgo.com", encrypt_plugin.RandomString(8))
	}
	myUser := MyUser{
		Email: a.AcmeEmail,
		key:   privateKey,
	}
	config := lego.NewConfig(&myUser)
	client, err := lego.NewClient(config)
	if err != nil {
		global.Logrus.Error(err.Error())
		return err
	}
	switch a.DNSProvider {
	case "AliCloud":
		provider, err = AliCloudProvider(a)
	case "Cloudflare":
		provider, err = CloudflareProvider(a)
	case "GoDaddy":
		provider, err = GoDaddyProvider(a)
	case "Hetzner":
		provider, err = HetznerProvider(a)
	case "TencentCloud":
		provider, err = TencentCloudProvider(a)
	default:
		global.Logrus.Error(err.Error())
		return err
	}
	client.Challenge.SetDNS01Provider(provider)
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		global.Logrus.Error(err.Error())
		return err
	}
	addressArr := strings.Fields(a.Address)
	var domainArr []string
	for k, _ := range addressArr {
		domainArr = append(domainArr, addressArr[k][:strings.Index(addressArr[k], ":")])
	}
	myUser.Registration = reg
	request := certificate.ObtainRequest{
		Domains: domainArr,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		global.Logrus.Error(err.Error())
		return err
	}
	a.PrivateKey = string(certificates.PrivateKey)
	a.Pem = string(certificates.Certificate)
	os.WriteFile("air.key", certificates.PrivateKey, 0644)
	os.WriteFile("air.cer", certificates.Certificate, 0644)

	return global.DB.Save(&a).Error

	//fmt.Printf("Certificate:%s\nIssuerCertificate:%s\nPrivateKey:%s\n", string(certificates.Certificate), string(certificates.IssuerCertificate), string(certificates.PrivateKey))
	//os.WriteFile("PrivateKey.key", certificates.PrivateKey, 0777)
	//os.WriteFile("Certificate.cer", certificates.Certificate, 0777)
	//os.WriteFile("IssuerCertificate.cer", certificates.IssuerCertificate, 0777)
}

func IsExpired(a *model.Acme) bool {
	if a.IsExpired {
		return true
	}
	conn, _ := tls.Dial("tcp", a.Address, nil)
	cert := conn.ConnectionState().PeerCertificates[0]

	fmt.Printf("NotAfter: %v\n", cert.NotAfter)
	now := time.Now()
	NotAfter := cert.NotAfter
	a.ExpiredAt = NotAfter

	if now.After(NotAfter) {
		fmt.Println("过期")
		a.IsExpired = true
	}
	fmt.Printf("没过期")
	return false
}

func AliCloudProvider(acme *model.Acme) (challenge.Provider, error) {
	config := alidns.NewDefaultConfig()
	config.APIKey = acme.AliCloudAccessKey
	config.SecretKey = acme.AliCloudSecretKey
	config.TTL = 3600
	return alidns.NewDNSProviderConfig(config)
}
func CloudflareProvider(acme *model.Acme) (challenge.Provider, error) {
	config := cloudflare.NewDefaultConfig()
	config.AuthToken = acme.CloudflareDnsApiToken
	config.TTL = 3600
	return cloudflare.NewDNSProviderConfig(config)
}

func GoDaddyProvider(acme *model.Acme) (challenge.Provider, error) {
	config := godaddy.NewDefaultConfig()
	config.APIKey = acme.GodaddyApiKey
	config.APISecret = acme.GodaddyApiSecret
	config.TTL = 3600
	return godaddy.NewDNSProviderConfig(config)
}
func HetznerProvider(acme *model.Acme) (challenge.Provider, error) {
	config := hetzner.NewDefaultConfig()
	config.APIKey = acme.HetznerApiKey
	config.TTL = 3600
	return hetzner.NewDNSProviderConfig(config)
}

func TencentCloudProvider(acme *model.Acme) (challenge.Provider, error) {
	config := tencentcloud.NewDefaultConfig()
	config.SecretID = acme.TencentCloudSecretId
	config.SecretKey = acme.TencentCloudSecretKey
	config.Region = ""
	config.SessionToken = ""
	config.TTL = 3600
	return tencentcloud.NewDNSProviderConfig(config)
}
