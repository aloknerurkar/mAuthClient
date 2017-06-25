package main

import (
	"bitbucket.org/aloknerurkar/mAuthClient"
	"net/http"
	"flag"
	"fmt"
	_ "strings"
	"runtime"
	"os/exec"
)

var (
	client_id      = flag.String("id", "CLIENT_ID", "OAuth Client ID.")
	client_secret  = flag.String("secret", "CLIENT_SECRET", "OAuth Client Secret.")
	redirect_uri   = flag.String("redirect", "http://localhost:1234/", "OAuth Client Redirect.")
)

type CallBack struct {}

func (c CallBack) OnError(error_str string) {
	fmt.Printf("Got error callback %s\n", error_str)
}

func (c CallBack) OnSuccess(user mAuthClient.AuthUser) {
	fmt.Printf("Got Success callback %+v", user)
}

var auth_client *mAuthClient.AuthClient
var auth_cb CallBack

func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Printf("Failed parsing form %v", err)
		return
	}
	code := req.Form.Get("code")
	fmt.Println("Authorization code: " + code)
	auth_client.Authenticate(code)
}

func main ()  {
	flag.Parse()
	auth_client = mAuthClient.NewAuthClient(mAuthClient.FACEBOOK, *client_id, *client_secret,
						*redirect_uri, auth_cb)
	fmt.Println("Go to URL " + auth_client.GetUserRedirect())
	openBrowser(auth_client.GetUserRedirect())

	http.HandleFunc("/", handler)
	http.ListenAndServe(":1234", nil)
}
