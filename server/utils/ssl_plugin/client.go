package ssl_plugin

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"os"
	"time"
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}

func (u *AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

type AcmeClient struct {
	Config *lego.Config
	Client *lego.Client
	User   *AcmeUser
}

func (c *AcmeClient) DnsCert(a *model.Acme) error {
	c.UseDns(a)
	privateKey, _ := certcrypto.GeneratePrivateKey(KeyType(a.KeyType))
	fmt.Println("privateKey:", privateKey)
	resource, err := c.ObtainSSL([]string{a.Address}, privateKey)
	if err != nil {
		return err
	}
	a.PrivateKey = string(resource.PrivateKey)
	a.Pem = string(resource.Certificate)
	a.CertURL = resource.CertURL
	certBlock, _ := pem.Decode(resource.Certificate)
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return err
	}
	a.ExpiredAt = cert.NotAfter
	a.StartAt = cert.NotBefore
	a.CommonName = cert.Issuer.CommonName
	a.Organization = cert.Issuer.Organization[0]
	//保存到数据库，上级保存

	//保存到文件
	c.SaveCertificateFile(resource)
	return nil
}

func (c *AcmeClient) UseDns(a *model.Acme) error {
	var (
		provider challenge.Provider
		err      error
	)
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

	return c.Client.Challenge.SetDNS01Provider(provider, dns01.AddDNSTimeout(10*time.Minute))

}

func (c *AcmeClient) ObtainSSL(domains []string, privateKey crypto.PrivateKey) (certificate.Resource, error) {
	request := certificate.ObtainRequest{
		Domains:    domains,
		Bundle:     true,
		PrivateKey: privateKey,
	}

	certificates, err := c.Client.Certificate.Obtain(request)
	if err != nil {
		return certificate.Resource{}, err
	}

	return *certificates, nil
}

func (c *AcmeClient) RenewSSL(certUrl string) (certificate.Resource, error) {
	certificates, err := c.Client.Certificate.Get(certUrl, true)
	if err != nil {
		return certificate.Resource{}, err
	}
	certificates, err = c.Client.Certificate.RenewWithOptions(*certificates, &certificate.RenewOptions{
		Bundle:         true,
		PreferredChain: "",
		MustStaple:     true,
	})
	if err != nil {
		return certificate.Resource{}, err
	}

	return *certificates, nil
}

func (c *AcmeClient) SaveCertificateFile(resource certificate.Resource) {
	os.WriteFile("air.key", resource.PrivateKey, 0644)
	os.WriteFile("air.cer", resource.Certificate, 0644)
}
