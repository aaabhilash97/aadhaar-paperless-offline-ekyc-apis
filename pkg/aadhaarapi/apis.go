package aadhaarapi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/yeka/zip"
)

const (
	valOtp     = "valOtp"
	genOtp     = "genOtp"
	genCaptcha = "genCaptcha"
)

// GetCaptcha - To fetch captcha and session id from UIDAI.
func GetCaptcha() (captchaImg []byte, sessionCookie string, err error) {
	task := genCaptcha
	url := "https://resident.uidai.gov.in/CaptchaSecurityImages.php?width=100&height=40&characters=5"
	var response *http.Response
	response, err = makeHttpCall("GET", url, nil, nil, 1)
	if err != nil {
		err = newAadhaarError(task, aadhaarError{
			msg:       err.Error(),
			retryable: true,
		})
		return
	}
	defer response.Body.Close()

	captchaImg, err = ioutil.ReadAll(response.Body)
	if err != nil {
		err = newAadhaarError(task, aadhaarError{
			msg:       err.Error(),
			retryable: true,
		})
		return
	} else if len(captchaImg) < 500 {
		err = newAadhaarError(task, aadhaarError{
			msg:       "Empty captcha image",
			retryable: true,
		})
		return
	}
	sessionCookie = response.Header.Get("Set-Cookie")
	return
}

type VerifyCaptchaOpt struct {
	SessionId    string // SessionId received from GetCaptcha
	UidNo        string // Aadhaar no
	SecurityCode string // Captcha
}

type VerifyCaptchaResult struct {
	Msg string
}

// VerifyCaptcha
// Once captcha fetched, need to call this API with aadhaar number(uid_no) and captcha(security_code).
// On success user will receive an OTP on phone number associated with aadhaar.
func VerifyCaptcha(opt VerifyCaptchaOpt) (result VerifyCaptchaResult, err error) {
	if !IsValidAadhaarNo(opt.UidNo) {
		err = &aadhaarError{
			msg:        "Invalid aadhaar no",
			invalidUid: true,
		}
		return
	}

	task := genOtp
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("uidno", opt.UidNo)
	_ = writer.WriteField("security_code", opt.SecurityCode)
	_ = writer.WriteField("task", task)
	_ = writer.WriteField("boxchecked", "0")
	_ = writer.WriteField("task", "genOtp")
	_ = writer.WriteField("boxchecked", "0")
	err = writer.Close()
	if err != nil {
		err = &aadhaarError{
			msg: err.Error(),
		}
		return
	}

	url := "https://resident.uidai.gov.in/offline-kyc"
	method := "POST"
	aadhaarRes, err := makeHttpCall(method, url, map[string]string{
		"Cookie":       opt.SessionId,
		"Content-Type": writer.FormDataContentType(),
	}, payload, 2)
	if err != nil {
		err = newAadhaarError(task, aadhaarError{
			msg: err.Error(),
		})
		return
	}
	defer aadhaarRes.Body.Close()
	pageRes, err := mapAadhaarPageResult(task, aadhaarRes.Body)
	if err != nil {
		return
	}

	result = VerifyCaptchaResult(pageRes)
	return
}

type OfflineAAdhaarXML struct {
	XMLName xml.Name `xml:"OfflinePaperlessKyc"`
	UidData struct {
		Poi struct {
			Dob        string `xml:"dob,attr"`
			EmailHash  string `xml:"e,attr"`
			Gender     string `xml:"gender,attr"`
			MobileHash string `xml:"m,attr"`
			Name       string `xml:"name,attr"`
		} `xml:"Poi"`
		Poa struct {
			CareOf     string `xml:"careof,attr"`
			Country    string `xml:"country,attr"`
			District   string `xml:"dist,attr"`
			House      string `xml:"house,attr"`
			Landmark   string `xml:"landmark,attr"`
			Locality   string `xml:"loc,attr"`
			Pincode    string `xml:"pc,attr"`
			Postoffice string `xml:"po,attr"`
			State      string `xml:"state,attr"`
			Street     string `xml:"street,attr"`
			Subdist    string `xml:"subdist,attr"`
			Vtc        string `xml:"vtc,attr"`
		} `xml:"Poa"`
		Photo string `xml:"Pht"`
	} `xml:"UidData"`
}

type VerifyOTPAndGetAadhaarResult struct {
	Details OfflineAAdhaarXML

	ZipFile []byte
	XmlFile []byte
}

type VerifyOTPAndGetAadhaarOpt struct {
	SessionId string // SessionId received from GetCaptcha
	OTP       string // OPT received in mobile or mAadhaar TOTP
	ZipCode   string
}

// VerifyOTPAndGetAadhaar
// download zip file from UIDAI and fetch details.
func VerifyOTPAndGetAadhaar(opt VerifyOTPAndGetAadhaarOpt) (result VerifyOTPAndGetAadhaarResult, err error) {
	task := valOtp

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("task", task)
	_ = writer.WriteField("boxchecked", "0")
	_ = writer.WriteField("zipcode", opt.ZipCode)
	_ = writer.WriteField("totp", opt.OTP)
	_ = writer.WriteField("task", task)
	_ = writer.WriteField("boxchecked", "0")
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	url := "https://resident.uidai.gov.in/offline-kyc"
	method := "POST"
	aadhaarRes, err := makeHttpCall(method, url, map[string]string{
		"Cookie":       opt.SessionId,
		"Content-Type": writer.FormDataContentType(),
	}, payload, 2)
	if err != nil {
		err = newAadhaarError(task, aadhaarError{
			msg: err.Error(),
		})
		return
	}
	defer aadhaarRes.Body.Close()

	if strings.Contains(aadhaarRes.Header.Get("Content-Type"), "/html") {
		// If response type is html, assuming request is error
		var pageRes aadhaarPageResult
		if pageRes, err = mapAadhaarPageResult(task, aadhaarRes.Body); err != nil {
			return
		} else {
			err = newAadhaarError(task, aadhaarError{
				aadhaarPageError: true,
				msg:              pageRes.Msg,
			})
			return
		}
	} else if strings.Contains(aadhaarRes.Header.Get("Content-Type"), "/zip") {
		// If response type is zip, assuming request is success

		result.ZipFile, err = ioutil.ReadAll(aadhaarRes.Body)
		if err != nil {
			err = newAadhaarError(task, aadhaarError{
				msg: "Unknown error",
			})
			return
		}
		ll := bytes.NewReader(result.ZipFile)
		var f *zip.Reader
		f, err = zip.NewReader(ll, int64(len(result.ZipFile)))

		if err != nil {
			err = newAadhaarError(task, aadhaarError{
				msg: "Unknown error",
			})
			return
		}
		for _, f := range f.File {
			if re := regexp.MustCompile(`offlineaadhaar\d+\.xml`); re.MatchString(f.Name) {
				if f.IsEncrypted() {
					f.SetPassword(opt.ZipCode)
				}
				var r io.ReadCloser
				r, err = f.Open()
				if err != nil {
					err = newAadhaarError(task, aadhaarError{
						msg: "Unknown error",
					})
					return
				}

				result.XmlFile, err = ioutil.ReadAll(r)
				if err != nil {
					err = newAadhaarError(task, aadhaarError{
						msg: "Unknown error",
					})
					return
				}
				xmlDetails := OfflineAAdhaarXML{}
				err = xml.Unmarshal(result.XmlFile, &xmlDetails)
				if err != nil {
					err = newAadhaarError(task, aadhaarError{
						msg: "Unknown error",
					})
					return
				}
				result.Details = xmlDetails
				return
			}
		}

	}
	// Assuming Some error occurred
	err = newAadhaarError(task, aadhaarError{
		msg: "Unknown error",
	})
	return
}
