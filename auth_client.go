package mAuthClient

import (
	"net/url"
	"golang.org/x/oauth2"
	ctx "context"
	"encoding/json"
	"io/ioutil"
	"log"
	"io"
	"os"
	"errors"
	"fmt"
)

type AuthClient struct {
	ProviderId int
	ClientId string
	ClientSecret string
	RedirectUrl string
	Callback AuthCallback
	elog *log.Logger
	ilog *log.Logger
}

type AuthCallback interface {
	OnError(msg string)
	OnSuccess(user AuthUser)
}

func NewAuthClient(provider_id int, client_id string, client_secret string, redirect_url string,
		   callback AuthCallback) *AuthClient {
	if len(client_id) == 0 || len(client_secret) == 0 || len(redirect_url) == 0 {
		fmt.Errorf("Invalid args given. One of ID:%s, SECRET:%s, REDIRECT_URL:%s", client_id, client_secret,
			   redirect_url)
		return nil
	}
	if provider_id != FACEBOOK && provider_id != GOOGLE {
		fmt.Errorf("Invalid provider %d", provider_id)
		return nil
	}
	client := new(AuthClient)
	client.init(provider_id, client_id, client_secret, redirect_url, callback)
	client.initLogger(os.Stdout)
	return client
}

func (a *AuthClient) init(provider_id int, client_id string, client_secret string, redirect_url string,
			  callback AuthCallback) {
	a.ProviderId = provider_id
	a.Callback = callback
	a.ClientSecret = client_secret
	a.ClientId = client_id
	a.RedirectUrl = redirect_url
}

func (a *AuthClient) initLogger(logger_op io.Writer) {
	a.elog = log.New(logger_op, "AuthClient\tERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	a.ilog = log.New(logger_op, "AuthClient\tINFO\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func (a *AuthClient) getOauthConfig() *oauth2.Config {
	var conf *oauth2.Config
	switch a.ProviderId {
	case FACEBOOK:
		conf = &oauth2.Config{
			ClientID:     a.ClientId,
			ClientSecret: a.ClientSecret,
			RedirectURL:  a.RedirectUrl,
			Scopes:       facebook.scopes,
			Endpoint:     facebook.endpoint,
		}
		break
	case GOOGLE:
		conf = &oauth2.Config{
			ClientID:     a.ClientId,
			ClientSecret: a.ClientSecret,
			RedirectURL:  a.RedirectUrl,
			Scopes:       google.scopes,
			Endpoint:     google.endpoint,
		}
		break
	}
	return conf
}

func (a *AuthClient) getReqURL() *url.URL {
	var ret_url *url.URL
	switch a.ProviderId {
	case FACEBOOK:
		ret_url = &facebook.req_url
		break
	case GOOGLE:
		ret_url = &google.req_url
		break
	}
	return ret_url
}

func (a *AuthClient) getAuthUser(val interface{}) AuthUser {
	ret_user := AuthUser{}
	switch a.ProviderId {
	case FACEBOOK:
		ret_user.Email = val.(*FacebookUser).Email
		ret_user.Name = val.(*FacebookUser).Name
		break
	case GOOGLE:
		ret_user.Email = val.(*GoogleUser).Emails[0].Value
		ret_user.Name = val.(*GoogleUser).Name
		break
	default:
		break
	}
	return ret_user
}

func (a *AuthClient) getUser(body io.ReadCloser) (AuthUser, error) {
	byte_body, err := ioutil.ReadAll(body)
	if err != nil {
		return AuthUser{}, err
	}

	var user interface{}
	switch a.ProviderId {
	case FACEBOOK:
		user = &FacebookUser{}
		break
	case GOOGLE:
		user = &GoogleUser{}
		break
	default:
		return AuthUser{}, errors.New("Invalid Provider")
	}

	err = json.Unmarshal(byte_body, &user)
	if err != nil {
		return AuthUser{}, err
	}

	return a.getAuthUser(user), nil
}

func (a *AuthClient) GetUserRedirect() string {
	a.ilog.Println("Getting Redirect URL.")
	return a.getOauthConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (a *AuthClient) Authenticate(authorization_code string) {
	a.ilog.Println("Authenticaticating user.")

	go func(c *AuthClient, code string) {
		conf := c.getOauthConfig()
		// Get Token
		token, err := conf.Exchange(ctx.Background(), code)
		if err != nil {
			c.elog.Printf("Error getting token. Err:%s Code:%s\n", err.Error(), "")
			c.Callback.OnError(err.Error())
			return
		}
		c.ilog.Printf("Got token %s.", token.AccessToken)

		// Get User info.
		client := conf.Client(oauth2.NoContext, token)
		resp, err := client.Get(c.getReqURL().String())
		if err != nil {
			c.elog.Printf("Error getting user. Err:%s Code:%s\n", err.Error(), "")
			c.Callback.OnError(err.Error())
			return
		}

		defer resp.Body.Close()

		// Parse User info.
		ret_user, err := c.getUser(resp.Body)
		if err != nil {
			c.elog.Printf("Error reading user response. Err:%s Code:%s\n", err.Error(), code)
			c.Callback.OnError(err.Error())
			return
		}

		c.ilog.Printf("Auth successful. USER %+v.", ret_user)
		c.Callback.OnSuccess(ret_user)

	}(a, authorization_code)
}
