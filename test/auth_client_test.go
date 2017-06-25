package test

import (
	"testing"
	"bitbucket.org/aloknerurkar/mAuthClient"
	"fmt"
)

type CallBack struct {}

func (c CallBack) OnError(error_str string) {
	fmt.Printf("Got error callback %s\n", error_str)
}

func (c CallBack) OnSuccess(user mAuthClient.AuthUser) {
	fmt.Printf("Got Success callback %+v", user)
}

func TestNewAuthClient(t *testing.T)  {

	cb := CallBack{}

	// Test Invalid params to NewAuthClient
	client := mAuthClient.NewAuthClient(mAuthClient.FACEBOOK, "", "", "", cb)
	if client != nil {
		t.Errorf("NewAuthClient expected to fail but succeeded client:%v", client)
	}

	client = mAuthClient.NewAuthClient(mAuthClient.FACEBOOK, "DUMMY_ID", "", "DUMMY_REDIRECT", cb)
	if client != nil {
		t.Errorf("NewAuthClient expected to fail but succeeded client:%v", client)
	}

	client = mAuthClient.NewAuthClient(mAuthClient.FACEBOOK, "DUMMY_ID", "DUMMY_SECRET", "", cb)
	if client != nil {
		t.Errorf("NewAuthClient expected to fail but succeeded client:%v", client)
	}

	client = mAuthClient.NewAuthClient(100, "DUMMY_ID", "DUMMY_SECRET", "DUMMY_REDIRECT", cb)
	if client != nil {
		t.Errorf("NewAuthClient expected to fail but succeeded client:%v", client)
	}

	client = mAuthClient.NewAuthClient(mAuthClient.FACEBOOK, "DUMMY_ID", "DUMMY_SECRET", "DUMMY_REDIRECT", cb)
	if client == nil {
		t.Error("NewAuthClient failed")
	}

}
