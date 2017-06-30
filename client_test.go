package hwpush

import (
	"testing"

	"golang.org/x/net/context"
)

func TestOauth(t *testing.T) {
	client, err := NewClient("10378558", "ld7uexh68ufi9x6dr74x0aq7wqncjftu")
	if err != nil {
		t.Error(err)
	}
	if client == nil {
		t.Fatal("client should not be null")
	}
}

func TestOauthFail(t *testing.T) {
	client, err := NewClient("10378558", "ld7uexh68ufi9x6dr74x0aq7wqncjftu1")

	if err == nil {
		t.Error("err should not be null")
	}

	if client != nil {
		t.Fatal("client should be null")
	}
}

func TestSendPush(t *testing.T) {
	client, err := NewClient("10378558", "ld7uexh68ufi9x6dr74x0aq7wqncjftu")

	if err != nil {
		t.Error(err)
	}

	message := NewNotification("1111", "1111")
	result, err := client.SendPush(context.Background(), []string{"0865297036030316300000558700CN01"}, message)
	if err != nil {
		t.Error(err)
	}
	if result.Code != SUCC {
		t.Fatalf("response code should be %s", SUCC)
	}
}

func TestSendPushFail(t *testing.T) {
	client, err := NewClient("10378558", "ld7uexh68ufi9x6dr74x0aq7wqncjftu")

	if err != nil {
		t.Error(err)
	}

	message := NewNotification("1111", "1111")
	result, err := client.SendPush(context.Background(), []string{"086195803749019730"}, message)
	if err == nil {
		t.Error(err)
	}
	if result != nil {
		t.Fatal("result should be nil")
	}
}
