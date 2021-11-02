# Aadhaar Paperless Offline e-KYC APIs

Aadhaar Paperless Offline e-KYC is a secure and shareable document which can be used by any Aadhaar holder for offline verification of identification.(https://resident.uidai.gov.in/offline-kyc)

## Introduction

This server exposing 3 APIs to download a ZIP file containing the Resident's Paperless Offline eKYC from UIDAI website.

- API to fetch captcha from UIDAI - GetCaptcha.
- API to enter aadhaar number and captcha received in GetCaptcha API.
  - By calling this API, user will get an OTP to the phone number associated to aadhaar.
- API to enter OTP received and four digit zipCode(share code).

These APIs depend on `https://resident.uidai.gov.in/offline-kyc` website.
By using these APIs client can expose their own interface to user.

### How To Run Server

```sh
go run cmd/server/main.go
```

OR 

```sh
go run cmd/server/main.go --conf configs/local.yml 
```

This server exposing 2 type of API interfaces(gRPC and REST).
Default port for gRPC server is 4444 and REST is 3333. (configs/default.yml).

This can be changed by adding environment specific configs.

Example:

Add `configs/local.yml` and run `go run cmd/server/main.go --conf configs/local.yml`

```yml
# sample config
server:
  grpcPort: 4444
  httpPort: 3333
  host: "0.0.0.0"

logger:
  outputPath: "stdout"
  level: "info"
  disableStackTrace: true

accessLogger:
  outputPath: "stdout"
  level: "error"
  disableStackTrace: true

```

### Use as a package

```sh
import "github.com/aaabhilash97/aadhaar_scrapper_apis/pkg/aadhaarapi"
```

# API Details

If client wants to use gRPC interface, client can generate client code from `api/proto` using protobuf compiler.

# REST API Details

## Get Captcha

To fetch captcha and session id from UIDAI.

* #### **gRPC FUNCTION NAME** - GetCaptcha
* #### **METHOD** - `GET`
* #### **URL** - `{{baseUrl}}/api/v1/GetCaptcha`
* #### **RESPONSE**
```json
      {
          "status": {
              "code": <string | response status code>,
              "message": <string | message>
          },
          "data": {
              "session_id": <string | session id>,
              "captcha_image": <string | base64 encoded image>
          }
      }
```

## Verify Captcha

Once captcha fetched, need to call this API with aadhaar number(uid_no) and captcha(security_code).

On success user will receive an OTP on phone number associated with aadhaar.

* #### **gRPC FUNCTION NAME** - VerifyCaptcha
* #### **METHOD** - `POST`
* #### **URL** - `{{baseUrl}}/api/v1/VerifyCaptcha`
* #### **REQUEST**
```json
      {
          "session_id": <string | session id>,
          "uid_no": <string | Enter your 12 digit Aadhaar number or 16 digit Virtual ID to begin.>,
          "security_code": <string | captcha received by user>
      }
```

* #### **RESPONSE**
```json
      {
          "status": {
              "code": <string | response status code>,
              "message": <string | message>
          }
      }
```

## Verify OTP and Fetch aadhaar details

Call this API with OTP received to download zip file from UIDAI and fetch details.

* #### **gRPC FUNCTION NAME** - VerifyOtpAndGetAadhaar
* #### **METHOD** - `POST`
* #### **URL** - `{{baseUrl}}/api/v1/VerifyOtpAndGetAadhaar`
* #### **REQUEST**
```json
      {
          "session_id": <string | session id>,
          "otp": <string | 6 digit OTP or 8 digit TOTP.>,
          "zip_code": <string | 4 characters>,
          "include_zip_file": <bool | to include zip file in response>,
          "include_xml_file": <bool | to include xml file in response>
      }
```

* #### **RESPONSE**
```json
{
    "status": {
        "code": <string | response status code>,
        "message": <string | message>
    },
    "data": {
        "details": {
            "poi": { // Proof of identity
                "dob": <string | Example: 10-09-1991>,
                "gender": <string | F or M>,
                "mobile_hash": <string>,
                "name": <string | Name as in aadhaar>
            },
            "poa": { // Proof of address
                "careof": <string>,
                "country":<string>,
                "dist": <string>,
                "house": <string>,
                "landmark":<string>,
                "pincode": <string>,
                "postoffice":<string>,
                "state": <string>,
                "street":<string>,
                "subdist": <string>,
                "vtc": <string>,
            },
            "photo": <string | Base64 encoded image of face photo>
        },
        "zip_file":  <string | Base64 encoded zip file downloaded from uidai>,
        "xml_file":  <string | Base64 encoded xml file adhaar Paperless Offline e-KYC>
    }
}
```

# Response Status Codes
* #### **2000** - Success
* #### **4000** - Request validation errors
* #### **5000** - Unknown error
* #### **5001** - UIDAI technical issue error
* #### **4002** - Invalid captcha
* #### **4003** - Invalid OTP
* #### **4004** - Invalid aadhaar number
* #### **4005** - Session expired
* #### **4006** - Invalid session id
* #### **4090** - Aadhaar already downloaded error.

