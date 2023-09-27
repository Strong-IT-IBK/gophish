package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/csv"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
	"github.com/jordan-wright/email"
)

var (
	firstNameRegex 			= regexp.MustCompile(`(?i)first[\s_-]*name$`)
	lastNameRegex			= regexp.MustCompile(`(?i)last[\s_-]*name$`)
	emailRegex				= regexp.MustCompile(`(?i)email$`)
	positionRegex			= regexp.MustCompile(`(?i)position$`)
	departmentRegex			= regexp.MustCompile(`(?i)department$`)
	departmentNumberRegex	= regexp.MustCompile(`(?i)department\s?number$`)
	ageRegex				= regexp.MustCompile(`(?i)age$`)
	genderRegex				= regexp.MustCompile(`(?i)gender$`)
	siteRegex				= regexp.MustCompile(`(?i)site$`)
	phoneRegex				= regexp.MustCompile(`(?i)phone$`)
	degreeRegex				= regexp.MustCompile(`(?i)degree$`)
	descriptionRegex		= regexp.MustCompile(`(?i)description$`)
)

// ParseMail takes in an HTTP Request and returns an Email object
// TODO: This function will likely be changed to take in a []byte
func ParseMail(r *http.Request) (email.Email, error) {
	e := email.Email{}
	m, err := mail.ReadMessage(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(m.Body)
	e.HTML = body
	return e, err
}

// ParseCSV contains the logic to parse the user provided csv file containing Target entries
func ParseCSV(r *http.Request) ([]models.Target, error) {
	mr, err := r.MultipartReader()
	ts := []models.Target{}
	if err != nil {
		return ts, err
	}
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		// Skip the "submit" part
		if part.FileName() == "" {
			continue
		}
		defer part.Close()
		reader := csv.NewReader(part)
		reader.TrimLeadingSpace = true
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		fi		:= -1
		li		:= -1
		ei		:= -1
		pi		:= -1
		depi	:= -1
		depni	:= -1
		agi		:= -1
		gi		:= -1
		si		:= -1
		phi		:= -1
		degi	:= -1
		desi	:= -1
		fn		:= ""
		ln		:= ""
		ea		:= ""
		ps		:= ""
		depn	:= ""
		depnn	:= ""
		agn		:= ""
		gn		:= ""
		sn		:= ""
		phn		:= ""
		degn	:= ""
		desn	:= ""
		for i, v := range record {
			switch {
			case firstNameRegex.MatchString(v):
				fi = i
			case lastNameRegex.MatchString(v):
				li = i
			case emailRegex.MatchString(v):
				ei = i
			case positionRegex.MatchString(v):
				pi = i
			case departmentRegex.MatchString(v):
				depi = i
			case departmentNumberRegex.MatchString(v):
				depni = i
			case ageRegex.MatchString(v):
				agi = i
			case genderRegex.MatchString(v):
				gi = i
			case siteRegex.MatchString(v):
				si = i
			case phoneRegex.MatchString(v):
				phi = i
			case degreeRegex.MatchString(v):
				degi = i
			case descriptionRegex.MatchString(v):
				desi = i
			}
		}
		if fi == -1 && li == -1 && ei == -1 && pi == -1 && depi == -1 && depni == -1 && agi == -1 && gi == -1 && si == -1 && phi == -1 && degi == -1 && desi == -1 {
			continue
		}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if fi != -1 && len(record) > fi {
				fn = record[fi]
			}
			if li != -1 && len(record) > li {
				ln = record[li]
			}
			if ei != -1 && len(record) > ei {
				csvEmail, err := mail.ParseAddress(record[ei])
				if err != nil {
					continue
				}
				ea = csvEmail.Address
			}
			if pi != -1 && len(record) > pi {
				ps = record[pi]
			}
			if depi != -1 && len(record) > depi {
				depn = record[depi]
			}
			if depni != -1 && len(record) > depni {
				depnn = record[depni]
			}
			if agi != -1 && len(record) > agi {
				agn = record[agi]
			}
			if gi != -1 && len(record) > gi {
				gn = record[gi]
			}
			if si != -1 && len(record) > si {
				sn = record[si]
			}
			if phi != -1 && len(record) > phi {
				phn = record[phi]
			}
			if degi != -1 && len(record) > degi {
				degn = record[degi]
			}
			if desi != -1 && len(record) > desi {
				desn = record[desi]
			}
			t := models.Target{
				BaseRecipient: models.BaseRecipient{
					FirstName:			fn,
					LastName:			ln,
					Email:				ea,
					Position:			ps,
					Department:			depn,
					DepNumber:			depnn,
					Age:				agn,
					Gender:				gn,
					Site:				sn,
					Phone:				phn,
					Degree:				degn,
					Desc:				desn,
				},
			}
			ts = append(ts, t)
		}
	}
	return ts, nil
}

// CheckAndCreateSSL is a helper to setup self-signed certificates for the administrative interface.
func CheckAndCreateSSL(cp string, kp string) error {
	// Check whether there is an existing SSL certificate and/or key, and if so, abort execution of this function
	if _, err := os.Stat(cp); !os.IsNotExist(err) {
		return nil
	}
	if _, err := os.Stat(kp); !os.IsNotExist(err) {
		return nil
	}

	log.Infof("Creating new self-signed certificates for administration interface")

	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return fmt.Errorf("error generating tls private key: %v", err)
	}

	notBefore := time.Now()
	// Generate a certificate that lasts for 10 years
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return fmt.Errorf("tls certificate generation: failed to generate a random serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Gophish"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return fmt.Errorf("tls certificate generation: failed to create certificate: %s", err)
	}

	certOut, err := os.Create(cp)
	if err != nil {
		return fmt.Errorf("tls certificate generation: failed to open %s for writing: %s", cp, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(kp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("tls certificate generation: failed to open %s for writing", kp)
	}

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return fmt.Errorf("tls certificate generation: unable to marshal ECDSA private key: %v", err)
	}

	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	keyOut.Close()

	log.Info("TLS Certificate Generation complete")
	return nil
}
