package ssl_plugin

import (
	"fmt"
	"testing"
)
import "github.com/ppoonk/AirGo/model"

func TestDnsCert(t *testing.T) {
	var a = model.Acme{
		AcmeEmail:             "ponk@oicq.com",
		AcmeMode:              "dns",
		AccountType:           "letsencrypt",
		KeyType:               "P256",
		Address:               "test.airgo.link",
		DNSProvider:           "TencentCloud",
		TencentCloudSecretId:  "",
		TencentCloudSecretKey: "",
		IsExpired:             false,
		IsRenewal:             false,
		PrivateKey:            "",
		Pem:                   "",
		CertURL:               "",
		CommonName:            "",
		Organization:          "",
	}
	client, _ := NewRegisterClient(&a)
	err := client.DnsCert(&a)
	fmt.Println("DnsCertTest a:", a)
	fmt.Println("DnsCertTest err:", err)

}
