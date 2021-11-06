package aadhaarapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func makeHttpCall(
	method, url string,
	headers map[string]string,
	body io.Reader,
	retry uint,
) (response *http.Response, err error) {
	for i := 0; i <= int(retry); i++ {
		client := &http.Client{}
		var req *http.Request
		if req, err = http.NewRequest(method, url, body); err != nil {
			return
		}

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		if response, err = client.Do(req); err != nil {
			continue
		} else {
			return
		}
	}
	return
}

type aadhaarPageResult struct {
	Msg string
}

func mapAadhaarPageResult(task string, body io.ReadCloser) (result aadhaarPageResult, err error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return
	}
	var aadhaarMessage string
	var isAadhaarReqError bool
	// Find error aadhaar page
	alertBox := doc.Find("#system-message-container>#system-message>div.alert")
	if alertBox.Length() == 0 {
		err = newAadhaarError(task, aadhaarError{
			msg: "Unknown error",
		})
		return
	}
	{

		alertBox.Each(func(i int, s *goquery.Selection) {
			if s.HasClass("alert-error") {
				isAadhaarReqError = true
			}
			aadhaarMessage = s.Find("div>div.alert-message").Text()
		})
		if isAadhaarReqError {
			err = newAadhaarError(task, aadhaarError{
				msg:              aadhaarMessage,
				aadhaarPageError: true,
			})
			return
		}
	}
	result = aadhaarPageResult{
		Msg: aadhaarMessage,
	}
	return
}

type VerifyAadhaarNumberPageResult struct {
	Msg          string
	IsVerified   bool
	AgeBand      string
	State        string
	Gender       string
	MobileNumber string
	Details      string
}

func mapVerifyAadhaarNumberPageResult(uidNo, task string, body io.ReadCloser) (result VerifyAadhaarNumberPageResult, err error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return
	}

	// Find error aadhaar page
	headerMsg := doc.Find("div.mb-15 > div.col-md-10.col-sm-10.col-xs-9.pl-0 > h2")
	if headerMsg == nil {
		err = newAadhaarError(task, aadhaarError{
			msg:              "unknown error",
			aadhaarPageError: true,
		})
	}
	if details := doc.Find("#maincontent"); details != nil {
		result.Details, err = details.Html()
		if err != nil {
			return
		}
	}

	if details := doc.Find("#maincontent > div > div.mb-15 > div.col-xs-12.my-20 > span:nth-child(1) > b"); details != nil {
		result.AgeBand, err = details.Html()
		if err != nil {
			return
		}
	}

	if details := doc.Find("#maincontent > div > div.mb-15 > div.col-xs-12.my-20 > span:nth-child(2) > b"); details != nil {
		result.Gender, err = details.Html()
		if err != nil {
			return
		}
	}

	if details := doc.Find("#maincontent > div > div.mb-15 > div.col-xs-12.my-20 > span:nth-child(3) > b"); details != nil {
		result.State, err = details.Html()
		if err != nil {
			return
		}
	}

	if details := doc.Find("#maincontent > div > div.mb-15 > div.col-xs-12.my-20 > span:nth-child(4) > b"); details != nil {
		result.MobileNumber, err = details.Html()
		if err != nil {
			return
		}
	}

	result.IsVerified = headerMsg.Text() == fmt.Sprintf("Aadhaar Number %s Exists!", uidNo)
	result.Msg = headerMsg.Text()
	return
}
