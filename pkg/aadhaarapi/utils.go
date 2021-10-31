package aadhaarapi

import (
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
