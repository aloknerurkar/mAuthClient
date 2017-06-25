AuthClient is designed to be an easy-to-use and portable OAuth2 client library. It can be used as is in a web client
or it can be compiled into android/ios libraries using the gomobile tool. The basic OAuth2 functionality has been
abstracted out into this library. As it is completely written in go, it is easily portable.

Flow of Operations:

1. Create a new AuthClient object with the following
    a. Provider ID: Check auth_client_global for supported providers. Currently its Facebook and Google. Plans to add
                    more providers in future. Pull requests welcome.
    b. Client ID: This is your client ID for the OAuth Service.You can get your client ID (or App ID) on your apps page
                  on the respective developers page.
    c. Client Secret: This is your client secret key. You are expected to know what these mean.
    d. Redirect Url: This is the redirect URL relative to client path. So typically in web servers it would be a
                     path relative to http://localhost:8080/auth/x. In case of android, this could be intent URI scheme.
    e. AuthCallback: This is the callback you need to implement in your app. Check the example client for more info.
    
2. AuthClient.GetUserRedirect
    This will return the constructed URL which can be used to initiate Oauth2. This URL will display the login
    page of the provider. On successful login, the redirect URL specified will be called with the Authorization
    code.
    
3. AuthClient.Authenticate (Authorization Code)
    After parsing the above redirect from the provider, we need to call Authenticate with the code. This will
    asynchronously begin the authentication. It will retrieve the token and also make a request to the provider to get
    user info. Currently only {Name, Email, Mobile} is returned. Typically on client side, these fields are enough
    to create a profile etc.
    
    
In the example directory, you will find a small Facebook login example. You need to supply your Client ID and Secret
to run it. It opens a browser with the login page and sets the redirect to localhost:1234/.

    auth_client = mAuthClient.NewAuthClient(mAuthClient.FACEBOOK, *client_id, *client_secret,
						                    *redirect_uri, auth_cb)
	fmt.Println("Go to URL " + auth_client.GetUserRedirect())
	openBrowser(auth_client.GetUserRedirect())
	
In the handler we will parse the code and then begin the authentication.

    err := req.ParseForm()
    if err != nil {
    	fmt.Printf("Failed parsing form %v", err)
    	return
    }
    code := req.Form.Get("code")
    fmt.Println("Authorization code: " + code)
    auth_client.Authenticate(code)
    
You should see the callback with the user information. In this case, we will have name and email.

The main reason, in having a two-step process is, so that we can handle different platforms. In different platforms,
the only thing that changes is how to display the URL to the user. In case of android, this could be a web-view.

This library is gomobile compatible. Which means, you can build it as android (.jar) / ios library. You can even
have a wrapper around this in android to abstract the part of displaying the web-page and handling the redirect.

Tests are very limited at the moment. Testing will require a considerable amount of mock interfaces. This is definitely
planned in the future. Any contributions are welcome.