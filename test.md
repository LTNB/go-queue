#   Creation - Technical Design

## Introduction
This document describes the system architecture design of `Account Creation` service, including database design and API design

**Design Goal**

- Verify code via SMS
- Create/update user login

## Database design
*account*
Core:
- `phone_number`: (varchar) unique id
- `first_name`: (varchar) user's first name
- `last_name`: (varchar) user's last name
- `email`: (varchar) user's email
- `country_id`: (varchar) id country
- `province`: (varchar) user's province
- `city`: (varchar) user's city
- `address_line`: (varchar) user's address line
- `postal_code`: (varchar)

**data_country**
- `id`: (varchar) unique id
- `name`: (varchar) name of country

**dialing_code**
- `code`: (varchar) unique code
- `name`: (varchar) abbreviation of the country

## API design
`Account Creation` service is provided as RESTful, input/output data is packaged as JSON
The API response with the following general format:
```json
{
  "status": "(int) result/error-code of API",
  "message": "(string) describes details about the result/error-code",
  "data": "(mixed) data output content (if has)"
}
```
The statuses are defined:
- 200: call API success.
- 400: input error (ex: parameter is wrong format).
- 403: client is not allow to call this api.
- 404: data not found.
- 500: server exception.
- 501: API is not existed.

URI of API is prefixed with value: `/api`.
That means, when call to API: `/example/:id` then the full URL API will be: `http(s)://host:port/api/example/:id`.

** API `GET /phone-number/dialing-code`**
Get dialing code list.
- Input: none
- Output:
    - Success:
        ```json
        {
          "status": 200,
          "message": "ok",
          "data": [
              {
                "code": "+84",
                "name": "vn"
              },
              {
                "code": "+1",
                "name": "us"
              }
             ]
        }  
    ```
    
** API `GET /phone-number/verify/:{phone-number}`**
Send input phone number and dialing code to login.
- Input:
    - Path variable:
        - `phone-number`: should use the E.164 format
- Output:
    - Success: `{"status": 200, "message": "ok" }`
    - Error 400: `{"status: 400, "message": "data wrong format"}`
    
** API `GET /verify-code/:{code}`**
Send input verification code, which is received via phone.
- Input:
    - Path variable:
        - `code`: the verification code, consists of 6 numeric characters
- Output:
    - Success: `{"status": 200, "message": "ok", "data": "JWT_TOKEN" }`
    - Error 401: `{"status: 401, "message": "login failed"}`
    
    
** API `GET /account?token={jwt_token}`
Get user info by JWT token
- Input:
    - params:
        - `jwt_token`: token after login success.
- Output:
    - Success: `{"status": 200, "message": "ok", "data": "USER_INFO" }`
    - ERROR 400: `{"status: 404, "message": "resource not found"}`
    - ERROR 403: `{"status: 403, "message": "access denied"}`
    
** API `POST /account?token={jwt_token}`
POST create/update user info
- Input:
    - params:
        - `jwt_token`: token after login success.
    - payload:
        `{ USER_INFOMATION_CONTENT }`
     
- Output:
    - Success: `{"status": 200, "message": "ok" }`
    - ERROR 400: `{"status: 400, "message": "data wrong format"}`

## Implementation

**GetDialingCode**
Core:
- Get list master data in `dialing_code` table.

**VerifyPhoneNumber**
Core:
- Validate phone number and dialing code format:
    - If success: create request to `Notification Service` to send verification code by phone number.
    - If failed: response 400, data wrong format.
    
**VerifyPhoneNumber**
Core:
- Receive verification code and validate format:
    - If success: 
        - Create request to `Notification Service` to verify verification code by phone number.
        - Get response from `Notification service`, if success, create JWT token base on phone number, or else response 401.
    - If failed: response 400, data wrong format.
    
**GetInfoByToken**
Core:
- Receive token and parse it to get the phone number.
- Using phone number query `account` table to get the info.
- Response user info if found, or else response 404.

**CreateUser/UpdateUser**
Core:
- Receive payload and validate require fields.
- Create new record in `account` table if not existed, or else update old record.
