package aadhaarapi

import (
	"fmt"
	"strings"
)

type aadhaarError struct {
	task             string
	msg              string
	retryable        bool
	invalidSessionId bool
	aadhaarPageError bool
	invalidUid       bool
}

func newAadhaarError(task string, err aadhaarError) *aadhaarError {
	err.task = task
	return &err
}

func (e *aadhaarError) Error() string {
	return e.msg
}

func IsAadhaarError(err error) bool {
	if _, ok := err.(*aadhaarError); ok {
		return true
	}
	return false
}

func IsRetryableError(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.retryable
	}
	return false
}
func IsInvalidSessionId(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.invalidSessionId
	}
	return false
}

func IsAadhaarPageError(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.aadhaarPageError
	}
	return false
}

func IsInvalidCaptcha(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.aadhaarPageError && err.msg == "Please Enter Valid Captcha"
	}
	return false
}

func IsOtpFailure(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.aadhaarPageError && strings.HasPrefix(err.msg, "OTP/TOTP Fail")
	}
	return false
}

func IsInvalidUIDOrVID(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.invalidUid ||
			(err.aadhaarPageError &&
				(strings.HasPrefix(err.msg, "Invalid UID") && err.task == genOtp))
	}
	return false
}

func IsSessionExpired(err error) bool {

	if err, ok := err.(*aadhaarError); ok {
		fmt.Println(err)
		return (IsInvalidUIDOrVID(err) && err.task == valOtp)
	}
	return false
}

func IsRedownloadError(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.aadhaarPageError &&
			err.msg == "You have already downloaded offline ekyc XML. Please check your downloads."
	}
	return false
}

func IsTechnicalError(err error) bool {
	if err, ok := err.(*aadhaarError); ok {
		return err.aadhaarPageError &&
			err.msg == "Technical issue please try after some time."
	}
	return false
}
