package hwpush

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/golang/glog"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

var (
	MaxRetryTimes = 3
)

type HuaweiPushClient struct {
	clientID, clientSecret string
	oauthResponse          *OauthResponse
}

type NspCtx struct {
	Version string `json:"ver"`
	AppId   string `json:"appId"`
}

func NewClient(clientID, clientSecret string) (*HuaweiPushClient, error) {
	client := &HuaweiPushClient{
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	oauthResponse, err := client.oauth()
	if err != nil {
		log.Warning(err)
		return nil, err
	}

	client.oauthResponse = oauthResponse
	go func() {
		for {
			sleepDuration := time.Duration(oauthResponse.ExpiresIn-1000) * time.Second
			time.Sleep(sleepDuration)
			retry := 0
		auth:
			oauthResponse, err = client.oauth()
			if err != nil {
				log.Warning(err)
				if retry < MaxRetryTimes {
					retry++
					time.Sleep(time.Duration(retry*2) * time.Second)
					goto auth
				} else {
					log.Fatal("please check network ,oauth failed too much", err)
				}
			}
			client.oauthResponse = oauthResponse
		}
	}()
	return client, nil
}

func (this HuaweiPushClient) oauth() (*OauthResponse, error) {
	params := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_secret=%s&client_id=%s", this.clientSecret, this.clientID))
	body, err := this.doPost(context.Background(), OAUTHURL, params)

	if err != nil {
		log.Warning(err)
		return nil, err
	}
	oauthResponse := &OauthResponse{}
	err = json.Unmarshal(body, oauthResponse)
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	if oauthResponse.AccessToken == "" {
		log.Warning("oauth failed")
		errs := &Error{}
		err = json.Unmarshal(body, errs)
		if err != nil {
			log.Warning(err)
			return nil, err
		}
		return nil, errs
	}
	return oauthResponse, nil
}

func (this HuaweiPushClient) SendPush(context context.Context, tokens []string, message Notification) (*Result, error) {
	nspTs := time.Now().Unix()
	expireTime := time.Now().Add(time.Hour * 2).Format("2006-01-02T15:04")
	payLoad, err := json.Marshal(message)
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	tokenString, err := json.Marshal(tokens)
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	bys, _ := json.Marshal(message)
	fmt.Println(string(bys))
	params := strings.NewReader(fmt.Sprintf("access_token=%s&nsp_svc=openpush.message.api.send&nsp_ts=%d&expire_time=%s&device_token_list=%s&payload=%s", url.QueryEscape(this.oauthResponse.AccessToken), nspTs, url.QueryEscape(expireTime), url.QueryEscape(string(tokenString)), url.QueryEscape(string(payLoad))))
	response, err := this.doPost(context, fmt.Sprintf(SENDPUSHURL, this.clientID), params)

	if err != nil {
		log.Warning(err)
		return nil, err
	}

	result := &Result{}
	err = json.Unmarshal(response, result)

	if err != nil {
		log.Warning(err)
		return nil, err
	}
	if result.Code != SUCC {
		log.Warning(string(response))
		return nil, Error{
			Description: result.Msg,
		}
	}
	return result, nil
}

func (this HuaweiPushClient) doPost(ctx context.Context, url string, body io.Reader) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	client := &http.Client{}
	tryTime := 0
tryAgain:
	res, err = ctxhttp.Do(ctx, client, req)
	//log.Info(res)
	if err != nil {
		log.Info("xiaomi push post err:", err, tryTime, res)
		select {
		case <-ctx.Done():
			return nil, err
		default:
		}
		tryTime += 1
		if tryTime < MaxRetryTimes {
			goto tryAgain
		}
		return nil, err
	}
	if res.Body == nil {
		return nil, errors.New("response is nil")
	}
	defer res.Body.Close()
	//fmt.Println("res.StatusCode=", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func push(oauthResponse *OauthResponse, clientId string, tokens []string, message Notification) error {
	nspTs := time.Now().Unix()
	expireTime := time.Now().Add(time.Hour * 2).Format("2006-01-02T15:04")
	payLoad, err := json.Marshal(message)
	if err != nil {
		log.Warning(err)
		return nil
	}
	tokenString, err := json.Marshal(tokens)
	if err != nil {
		log.Warning(err)
		return nil
	}
	params := strings.NewReader(fmt.Sprintf("access_token=%s&nsp_svc=openpush.message.api.send&nsp_ts=%d&expire_time=%s&device_token_list=%s&payload=%s", url.QueryEscape(oauthResponse.AccessToken), nspTs, url.QueryEscape(expireTime), url.QueryEscape(string(tokenString)), url.QueryEscape(string(payLoad))))
	resp, err := http.Post(fmt.Sprintf(SENDPUSHURL, clientId), "application/x-www-form-urlencoded", params)
	if err != nil {
		log.Warning(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warning(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
