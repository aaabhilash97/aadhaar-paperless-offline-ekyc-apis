package aadhaarapi

import (
	"crypto/x509"
	"encoding/pem"
	"regexp"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

var d = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 2, 3, 4, 0, 6, 7, 8, 9, 5},
	{2, 3, 4, 0, 1, 7, 8, 9, 5, 6},
	{3, 4, 0, 1, 2, 8, 9, 5, 6, 7},
	{4, 0, 1, 2, 3, 9, 5, 6, 7, 8},
	{5, 9, 8, 7, 6, 0, 4, 3, 2, 1},
	{6, 5, 9, 8, 7, 1, 0, 4, 3, 2},
	{7, 6, 5, 9, 8, 2, 1, 0, 4, 3},
	{8, 7, 6, 5, 9, 3, 2, 1, 0, 4},
	{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
}
var p = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 5, 7, 6, 2, 8, 3, 0, 9, 4},
	{5, 8, 0, 3, 7, 9, 6, 1, 4, 2},
	{8, 9, 1, 6, 0, 4, 3, 5, 2, 7},
	{9, 4, 5, 3, 1, 2, 6, 8, 7, 0},
	{4, 2, 8, 6, 5, 7, 3, 9, 0, 1},
	{2, 7, 9, 3, 8, 0, 6, 4, 1, 5},
	{7, 0, 4, 6, 9, 1, 3, 2, 5, 8},
}
var j = []int{0, 4, 3, 2, 1, 5, 6, 7, 8, 9}

func IsValidAadhaarNo(aadhaar string) bool {
	c := 0
	if re := regexp.MustCompile(`^\d+$`); !re.MatchString(aadhaar) {
		return false
	}

	valSplit := strings.Split(aadhaar, "")
	var checkSum, _ = strconv.Atoi(valSplit[len(valSplit)-1])
	valSplit = valSplit[:len(valSplit)-1]
	// Reverse array
	for i, j := 0, len(valSplit)-1; i < j; i, j = i+1, j-1 {
		valSplit[i], valSplit[j] = valSplit[j], valSplit[i]
	}
	for i := 0; i < len(valSplit); i++ {
		val, _ := strconv.Atoi(valSplit[i])
		c = d[c][p[((i + 1) % 8)][val]]
	}
	return checkSum == j[c]
}

const aadhaarLatestPublicKey = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`

func ValidateXMLSignature(xmlData []byte) (bool, error) {
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(xmlData)
	if err != nil {
		return false, err
	}

	block, _ := pem.Decode([]byte(aadhaarLatestPublicKey))
	if err != nil {
		return false, err
	}
	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, err
	}
	ctx := dsig.NewDefaultValidationContext(&dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{pub},
	})

	if el := doc.FindElement("OfflinePaperlessKyc"); el != nil {
		el, err = ctx.Validate(el)
		if err != nil {
			return false, err
		}
		return true, err
	}
	return false, err
}
