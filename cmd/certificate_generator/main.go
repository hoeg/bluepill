package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gopkg.in/yaml.v2"
)

// take a commandline argument that is the common name and generate a self signed certificate.
func main() {
	// Check if the common name argument is provided
	if len(os.Args) < 3 {
		fmt.Println("Please provide the service name and namespace as a commandline argument. e.g. bluepill default")
		return
	}

	// Get the common name from commandline argument
	commonName := fmt.Sprintf("%s.%s.svc", os.Args[1], os.Args[2])

	// Generate a private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		return
	}

	// Create a self-signed certificate template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{commonName},
	}

	// Create a self-signed certificate using the private key and template
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Println("Failed to create certificate:", err)
		return
	}

	//encode certificate and private key to PEM format in memory, do not write them to a file
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	// Create a Kubernetes Secret object
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-tls", os.Args[1]),
			Namespace: os.Args[2],
		},
		StringData: map[string]string{
			"cert.pem": base64.StdEncoding.EncodeToString(certPEM),
			"key.pem":  base64.StdEncoding.EncodeToString(keyPEM),
		},
		Type: corev1.SecretTypeOpaque,
	}

	// Convert the Secret object to YAML encoding, omit empty fields
	secretYAML, err := yaml.Marshal(secret)
	if err != nil {
		fmt.Println("Failed to marshal Secret to YAML:", err)
		return
	}

	//write secretYAML to a file named secret.yaml
	err = os.WriteFile("secret.yaml", secretYAML, 0644)
	if err != nil {
		fmt.Println("Failed to write Secret YAML to file:", err)
		return
	}

	//also write the files to testdata/cert.pem and testdata/key.pem
	err = os.WriteFile("testdata/cert.pem", certPEM, 0644)
	if err != nil {
		fmt.Println("Failed to write cert.pem to file:", err)
		return
	}
	err = os.WriteFile("testdata/key.pem", keyPEM, 0644)
	if err != nil {
		fmt.Println("Failed to write key.pem to file:", err)
		return
	}
}
