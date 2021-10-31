package service

import (
	"context"
	"fmt"

	"github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/aadhaarapi"
	api "github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	ApiSuccessCode        = "2000"
	ValidationError       = "4000"
	ApiUnknownError       = "5000"
	AadhaarTechnicalIssue = "5001"

	InvalidCaptchaError = "4002"
	InvalidOTPError     = "4003"
	InvalidUIDError     = "4004"

	SessionExpiredError   = "4005"
	InvalidSessionIdError = "4006"

	RedownloadError = "4090"
)

var codeToMessageMapping = map[string]string{
	ApiSuccessCode:        "SUCCESS",
	ApiUnknownError:       "Some error occurred.",
	InvalidSessionIdError: "Invalid session",
}

func stringOR(args ...string) string {
	for _, arg := range args {
		if len(arg) > 0 {
			return arg
		}
	}
	return ""
}

func mapAadhaarErrToStatus(ctx context.Context, err error) *api.ResponseStatus {
	if aadhaarapi.IsInvalidCaptcha(err) {
		return mapToStatus(ctx, InvalidCaptchaError, err.Error())
	} else if aadhaarapi.IsOtpFailure(err) {
		return mapToStatus(ctx, InvalidOTPError, err.Error())
	} else if aadhaarapi.IsInvalidUIDOrVID(err) {
		return mapToStatus(ctx, InvalidUIDError, err.Error())
	} else if aadhaarapi.IsSessionExpired(err) {
		return mapToStatus(ctx, SessionExpiredError, err.Error())
	} else if aadhaarapi.IsInvalidSessionId(err) {
		return mapToStatus(ctx, InvalidSessionIdError, err.Error())
	} else if aadhaarapi.IsRedownloadError(err) {
		return mapToStatus(ctx, RedownloadError, err.Error())
	} else if aadhaarapi.IsTechnicalError(err) {
		return mapToStatus(ctx, AadhaarTechnicalIssue, err.Error())
	}
	return mapToStatus(ctx, ApiUnknownError, "")
}

func mapToStatus(ctx context.Context, code, msg string) *api.ResponseStatus {
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", code[0:3]))
	if definedMsg, ok := codeToMessageMapping[code]; ok {
		return &api.ResponseStatus{
			Code:    code,
			Message: stringOR(msg, definedMsg),
		}
	}
	return &api.ResponseStatus{
		Code:    code,
		Message: msg,
	}
}

type validationErr interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
}

func validationErrToStaus(ctx context.Context, e validationErr) *api.ResponseStatus {
	if e == nil {
		return nil
	}
	status := mapToStatus(ctx, ValidationError, fmt.Sprintf("Invalid %s", e.Field()))
	status.ValidationDetails = []*api.ResponseStatus_ValidationErrDetail{
		{
			Field:  e.Field(),
			Reason: e.Reason(),
		},
	}
	return status
}
