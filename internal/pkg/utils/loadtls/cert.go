package loadtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

//func LoadTLSCredentials(name string) (credentials.TransportCredentials, error) {
//	serverCert, err := tls.LoadX509KeyPair(fmt.Sprintf("cert/techno_%s.crt", name), "cert/techno.key")
//	if err != nil {
//		return nil, err
//	}
//
//	config := &tls.Config{ //nolint
//		Certificates: []tls.Certificate{serverCert},
//		ClientAuth:   tls.NoClientCert,
//	}
//
//	return credentials.NewTLS(config), nil
//}

func LoadTLSCredentials(name string) (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(fmt.Sprintf("cert/techno_%s.crt", name), "cert/techno.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certification: %w", err)
	}

	data, err := os.ReadFile("cert/CA.pem")
	if err != nil {
		return nil, fmt.Errorf("faild to read CA certificate: %w", err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("unable to append the CA certificate to CA pool")
	}

	tlsConfig := &tls.Config{ //nolint
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig), nil
}

//func LoadTLSClientCredentials() (credentials.TransportCredentials, error) {
//	pemServerCA, err := os.ReadFile("cert/CA.pem")
//	if err != nil {
//		return nil, err
//	}
//
//	certPool := x509.NewCertPool()
//	if !certPool.AppendCertsFromPEM(pemServerCA) {
//		return nil, fmt.Errorf("failed to add server CA's certificate")
//	}
//
//	config := &tls.Config{ //nolint
//		RootCAs: certPool,
//	}
//
//	return credentials.NewTLS(config), nil
//}

func LoadTLSClientCredentials() (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair("cert/techno_auth.crt", "cert/techno.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load client certification: %w", err)
	}

	ca, err := os.ReadFile("cert/CA.pem")
	if err != nil {
		return nil, fmt.Errorf("faild to read CA certificate: %w", err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("faild to append the CA certificate to CA pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      capool,
	}

	return credentials.NewTLS(tlsConfig), nil
}
