package v2

import (
	"crypto/tls"
	"net"
	"os"
	"time"

	"github.com/ONSdigital/dp-document-db/certs"
	"github.com/globalsign/mgo"
)

const caFilePath = "rds-combined-ca-bundle.pem"

func NewSession(username, password string) (*mgo.Session, error) {
	host := os.Getenv("DB_ENDPOINT")

	tlsConfig, err := certs.GetCustomTLSConfig(caFilePath)
	if err != nil {
		return nil, err
	}

	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:        []string{host},
		Timeout:      time.Second * 5,
		Username:     username,
		Password:     password,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		ReadPreference: &mgo.ReadPreference{
			Mode: mgo.SecondaryPreferred,
		},
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", host, tlsConfig)
		},
	})
}
