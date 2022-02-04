package logger

import (
	"time"

	"github.com/rs/zerolog"
)

type (
	// AppLogger ..
	AppLogger struct {
		Event          string        `json:"event" bson:"event"`
		ConnectionType string        `json:"connection_type" bson:"connection_type"`
		Package        string        `json:"package" bson:"package"`
		TimerStart     time.Time     `json:"time_start" bson:"time_start"`
		ElapsedTime    time.Duration `json:"elapsed_time" bson:"elapsed_time"`
		Data           interface{}   `json:"data" bson:"data"`
		// Adapter Log
		RequestURL     string      `json:"request_url" bson:"request_url"`
		RequestHost    string      `json:"request_host" bson:"request_host"`
		RequestMethod  string      `json:"request_method" bson:"request_method"`
		RequestIP      string      `json:"request_ip" bson:"request_ip"`
		RequestTime    string      `json:"request_time" bson:"request_time"`
		Request        interface{} `json:"request" bson:"request"`
		RequestHeader  interface{} `json:"request_header" bson:"request_header"`
		ResponseTime   string      `json:"response_time" bson:"response_time"`
		Response       interface{} `json:"response" bson:"response"`
		ResponseHeader interface{} `json:"response_header" bson:"response_header"`
		HTTPCode       int         `json:"http_code" bson:"http_code"`
		RawData        string      `json:"rawdata" bson:"rawdata"`
		RawRequest     string      `json:"raw_request" bson:"raw_request"`
		RawResponse    string      `json:"raw_response" bson:"raw_response"`
		// Error Log
		DetailStr interface{} `json:"detail_str" bson:"detail_str"`
	}

	// StrObject ..
	StrObject struct {
		Raw string `json:"raw" bson:"raw"`
	}
)

// Info level : log default
func (l AppLogger) Info() *zerolog.Event {
	l.ElapsedTime = time.Since(l.TimerStart) / 1e6

	switch v := l.Data.(type) {
	case string:
		strObj := StrObject{
			Raw: string(v),
		}
		l.Data = strObj
	}

	switch v := l.Request.(type) {
	case string:
		strObj := StrObject{
			Raw: string(v),
		}
		l.Request = strObj
	}

	switch v := l.Response.(type) {
	case string:
		strObj := StrObject{
			Raw: string(v),
		}
		l.Response = strObj
	}

	return Info().
		// Base Log
		Str("event", l.Event).
		Interface("data", l.Data).
		Str("connection_type", l.ConnectionType).
		Str("package", l.Package).
		Interface("time_start", l.TimerStart).
		Interface("elapsed_time", l.ElapsedTime).
		Str("raw_data", l.RawData).
		Str("request_url", l.RequestURL).
		Str("request_host", l.RequestHost).
		Str("http_method", l.RequestMethod).
		Str("request_ip", l.RequestIP).
		Str("request_time", l.RequestTime).
		Interface("request", l.Request).
		Interface("request_header", l.RequestHeader).
		Str("response_time", l.ResponseTime).
		Interface("response", l.Response).
		Interface("response_header", l.ResponseHeader).
		Int("http_code", l.HTTPCode).
		Str("raw_request", l.RawRequest).
		Str("raw_response", l.RawResponse)
}

// Error level : log default
func (l AppLogger) Error() *zerolog.Event {
	l.ElapsedTime = time.Since(l.TimerStart) / 1e6

	switch v := l.Request.(type) {
	case string:
		strObj := StrObject{
			Raw: string(v),
		}
		l.Request = strObj
	}

	switch v := l.Response.(type) {
	case string:
		strObj := StrObject{
			Raw: string(v),
		}
		l.Response = strObj
	}

	return Info().
		// Base Log
		Str("event", l.Event).
		Interface("data", l.Data).
		Str("connection_type", l.ConnectionType).
		Str("package", l.Package).
		Interface("time_start", l.TimerStart).
		Interface("elapsed_time", l.ElapsedTime).
		Str("raw_data", l.RawData).
		Str("request_url", l.RequestURL).
		Str("request_host", l.RequestHost).
		Str("http_method", l.RequestMethod).
		Str("request_ip", l.RequestIP).
		Str("request_time", l.RequestTime).
		Interface("request", l.Request).
		Interface("request_header", l.RequestHeader).
		Str("response_time", l.ResponseTime).
		Interface("response", l.Response).
		Interface("response_header", l.ResponseHeader).
		Int("http_code", l.HTTPCode).
		Interface("detail_str", l.DetailStr).
		Str("raw_request", l.RawRequest).
		Str("raw_response", l.RawResponse)
}
