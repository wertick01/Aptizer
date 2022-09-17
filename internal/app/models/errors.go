package models

import "errors"

var ErrNoRecord = errors.New("models: подходящей записи не найдено")
var ExpiredToken = errors.New("Token expired.")
var DecodeErr = errors.New("Decode data error. Please check that the data is entered correctly.")
var DatabaseQueryError = errors.New("Database query error.")
var WrongPassword = errors.New("Incorrect password.")
var RefreshTokenError = errors.New("Error while creation refresh-token")
var JWTTokenError = errors.New("Error while creation refresh-token")
var ErrorCookieToken = errors.New("No JWT-Token in cookie. Check the authorization.")
var ErrorParseToken = errors.New("Error of parsing JWT-Token.")
var NonValidToken = errors.New("Non valid token.")
var ExpiredSessionTime = errors.New("Expired session time.")
var UpdateJWTTokenError = errors.New("Error while updating JWT-Token.")
var ErrorWrongJWTToken = errors.New("JWT-Token extracting error.")
var ErrorParsingID = errors.New("Error while parsing id (check the correction of argument).")
var ExtractingParameterError = errors.New("Error While exctracting the parameter.")
var UserNotFound = errors.New("User not found")
