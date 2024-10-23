package tests

import (
	"github.com/gavv/httpexpect/v2"
	"net/url"
	"testing"
)

const (
	addr     = "localhost:8090"
	admToken = "randomToken"
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		testName     string
		login        string
		pswd         string
		token        string
		codeResponse int
	}{
		{
			testName:     "Access register",
			login:        "longUser3",
			pswd:         "Password3&",
			token:        admToken,
			codeResponse: 200,
		},
		{
			testName:     "Already Exists User",
			login:        "longUser3",
			pswd:         "Password3&",
			token:        admToken,
			codeResponse: 400,
		},
		{
			testName:     "Short login",
			login:        "short",
			pswd:         "Password3&",
			token:        admToken,
			codeResponse: 400,
		},
		{
			testName:     "Login not digit",
			login:        "userNotDigit",
			pswd:         "Password3&",
			token:        admToken,
			codeResponse: 400,
		},
		{
			testName:     "Password not special symbol",
			login:        "userCool7",
			pswd:         "Password3",
			token:        admToken,
			codeResponse: 400,
		},
		{
			testName:     "Password not digit",
			login:        "userCool8",
			pswd:         "Password&",
			token:        admToken,
			codeResponse: 400,
		},
		{
			testName:     "Password short",
			login:        "userCool9",
			pswd:         "Short7&",
			token:        admToken,
			codeResponse: 400,
		},
		{
			testName:     "Password small symbols",
			login:        "userCool1",
			pswd:         "smallsymbols7&",
			token:        admToken,
			codeResponse: 400,
		},
		{
			testName:     "Token for register no valid",
			login:        "longUser4",
			pswd:         "password3&",
			token:        "tokenUser",
			codeResponse: 403,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   addr,
			}

			e := httpexpect.Default(t, u.String())

			_ = e.POST("/api/register").
				WithJSON(map[string]string{
					"login": tc.login,
					"pswd":  tc.pswd,
					"token": tc.token,
				},
				).Expect().
				Status(tc.codeResponse).
				JSON().
				Object()
		})
	}

}

func TestAuth(t *testing.T) {
	testCases := []struct {
		testName     string
		login        string
		pswd         string
		codeResponse int
	}{
		{
			testName:     "Access auth",
			login:        "longUser3",
			pswd:         "Password3&",
			codeResponse: 200,
		},
		{
			testName:     "No user",
			login:        "noUser",
			pswd:         "Password3&",
			codeResponse: 404,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   addr,
			}

			e := httpexpect.Default(t, u.String())

			_ = e.POST("/api/auth").
				WithJSON(map[string]string{
					"login": tc.login,
					"pswd":  tc.pswd,
				},
				).Expect().
				Status(tc.codeResponse).
				JSON().
				Object()
		})
	}
}
