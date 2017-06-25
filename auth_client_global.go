package mAuthClient

import (
	"net/url"
	"golang.org/x/oauth2"
)

const (
	FACEBOOK = 1
	GOOGLE = 2
	PHONE = 3
)

type AuthUser struct {
	Name   string
	Email  string
	Mobile string
}

type FacebookUser struct {
	ID         string `json:"id"`
	Name   	   string `json:"name"`
	Email      string `json:"email"`
	ProfileUrl string `json:"link"`
}

type Email struct {
	Value string `json:"value" binding:"required"`
	Type  string `json:"type" binding:"required"`
}

type Image struct {
	URL       string `json:"url" binding:"required"`
	IsDefault string `json:"isDefault" binding:"required"`
}

type GoogleUser struct {
	ID         string  `json:"id"`
	Username   string  `json:"nickname"`
	Emails     []Email `json:"emails"`
	Name       string  `json:"displayName"`
	Image      Image   `json:"image"`
	ProfileUrl string  `json:"url"`
}

type provider_extras struct {
	req_url url.URL
	scopes []string
	endpoint oauth2.Endpoint
}

var facebook = provider_extras{
	req_url: url.URL{
		Scheme: "https",
		Host: "graph.facebook.com",
		Opaque: "//graph.facebook.com/me",
	},
	scopes: []string {
		"public_profile",
		"email",
	},
	endpoint: oauth2.Endpoint{
		AuthURL: "https://www.facebook.com/dialog/oauth",
		TokenURL: "https://graph.facebook.com/oauth/access_token",
	},
}

var google = provider_extras{
	req_url: url.URL{
		Scheme: "https",
		Host: "www.googleapis.com",
		Opaque: "//www.googleapis.com/plus/v1/people/me",
	},
	scopes: []string {
		"email",
		"profile",
		"https://www.googleapis.com/auth/plus.login",
		"https://www.googleapis.com/auth/plus.profile.emails.read",
	},
	endpoint: oauth2.Endpoint{
		AuthURL: "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://accounts.google.com/o/oauth2/token",
	},
}
