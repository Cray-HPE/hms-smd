/*
 * MIT License
 *
 * (C) Copyright [2018-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

func (s *SmD) setupCerts(certPath string, keyPath string) error {
	err := s.checkCertPath(certPath, keyPath)
	// Generate certificate if one does not already exist.
	if err != nil {
		// This will obviously fail.
		if certPath == "" || keyPath == "" {
			s.LogAlways("Cert or key path was the empty string")
			return err
		}
		s.LogAlways("generate Certs")
		hostname, hostnameErr := os.Hostname()
		if hostnameErr != nil {
			s.LogAlways("Getting hostname: %s", err)
		}
		err = s.generateCerts(certPath, keyPath, hostname)
		if err != nil {
			s.LogAlways("Error: Couldn't create https certs.")
		}
		return err
	}
	s.LogAlways("Found existing certs")
	return nil
}

//CheckCertPath - See if certs already exist
func (s *SmD) checkCertPath(certPath string, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return err
	} else if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

//
// GenerateCerts - create rsa certs / keys
//
func (s *SmD) generateCerts(certPath string, keyPath string, host string) error {
	var priv interface{}
	var err error
	validFor := 10 * 365 * 24 * time.Hour
	rsaBits := 2048

	priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		s.LogAlways("Failed to generate private key: %s!", err)
		return err
	}

	var notBefore = time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		s.LogAlways("Failed to generate serial number: %s!", err)
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Cray Inc."},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader,
		&template,
		&template,
		&priv.(*rsa.PrivateKey).PublicKey,
		priv)
	if err != nil {
		s.LogAlways("Failed to create certificate: %s!", err)
		return err
	}
	certOut, err := os.Create(certPath)
	if err != nil {
		s.LogAlways("Failed to open %s for writing: %s!", certPath, err)
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	s.LogAlways("Wrote %s", certPath)

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		s.LogAlways("Failed to open %s for writing: %s", keyPath, err)
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv.(*rsa.PrivateKey))})
	keyOut.Close()
	s.LogAlways("Wrote %s", keyPath)
	return nil
}
