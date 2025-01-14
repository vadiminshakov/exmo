/*
   Copyright 2019 Vadim Inshakov

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package exmo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ApiResponse is a map for API responses.
type ApiResponse map[string]interface{}

// ApiParams is a map for API calls' params.
type ApiParams map[string]string

// Exmo holds client-specific info.
type Exmo struct {
	key    string // public key
	secret string // secret key
	client *http.Client
}

// Api creates Exmo instance with specified credentials.
func Api(key string, secret string) Exmo {
	var netTransport = &http.Transport{
		MaxIdleConns:        30,
		MaxConnsPerHost:     1,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	return Exmo{key, secret, client}
}

// Api_query is a general query method for API calls.
func (ex *Exmo) Api_query(mode string, method string, params ApiParams) (ApiResponse, error) {

	post_params := url.Values{}
	if mode == "authenticated" {
		post_params.Add("nonce", nonce())
	}
	if params != nil {
		for key, value := range params {
			post_params.Add(key, value)
		}
	}
	post_content := post_params.Encode()

	sign := ex.Do_sign(post_content)

	req, _ := http.NewRequest("POST", "https://api.exmo.com/v1/"+method, bytes.NewBuffer([]byte(post_content)))
	req.Header.Set("Key", ex.key)
	req.Header.Set("Sign", sign)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(post_content)))

	resp, err := ex.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, errors.New("http status: " + resp.Status)
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}

	var dat map[string]interface{}
	err2 := json.Unmarshal([]byte(body), &dat)
	if err2 != nil {
		return nil, err2
	}

	if result, ok := dat["result"]; ok && result.(bool) != true {
		return nil, errors.New(dat["error"].(string))
	}

	return dat, nil
}

// nonce generates request parameter ‘nonce’ with incremental numerical value (>0). The incremental numerical value should never reiterate or decrease.
func nonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Do_sign encrypts POST data (param=val&param1=val1) with method HMAC-SHA512 using secret key; the secret key also can be found in user’s profile settings.
func (ex *Exmo) Do_sign(message string) string {
	mac := hmac.New(sha512.New, []byte(ex.secret))
	mac.Write([]byte(message))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

/*
   Public API
*/

// GetTrades return list of the deals on currency pairs.
func (ex *Exmo) GetTrades(pair string) (ApiResponse, error) {
	return ex.Api_query("public", "trades", ApiParams{"pair": pair})
}

// GetOrderBook return the book of current orders on the currency pair.
func (ex *Exmo) GetOrderBook(pair string, limit int) (ApiResponse, error) {
	if limit < 100 || limit > 1000 {
		return nil, errors.New("limit param must be in range of 100-1000")
	}

	return ex.Api_query("public", "order_book", ApiParams{"pair": pair, "limit": string(limit)})
}

// Ticker return statistics on prices and volume of trades by currency pairs.
func (ex *Exmo) Ticker() (ApiResponse, error) {
	return ex.Api_query("public", "ticker", ApiParams{})
}

// GetPairSettings return currency pairs settings.
func (ex *Exmo) GetPairSettings() (ApiResponse, error) {
	return ex.Api_query("public", "pair_settings", ApiParams{})
}

// GetCurrency return currencies list.
func (ex *Exmo) GetCurrency() ([]string, error) {
	resp, err := http.Get("https://api.exmo.com/v1/currency")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return nil, errors.New("http status: " + resp.Status)
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}

	var dat []string
	err2 := json.Unmarshal(body, &dat)
	if err2 != nil {
		return nil, err2
	}

	return dat, nil
}

/*
   Authenticated API
*/

// GetUserInfo return information about user's account.
func (ex *Exmo) GetUserInfo() (ApiResponse, error) {
	return ex.Api_query("authenticated", "user_info", nil)
}

// GetUserTrades return the list of user’s deals.
func (ex *Exmo) GetUserTrades(pair string, offset, limit int) (ApiResponse, error) {
	return ex.Api_query("authenticated", "user_trades", ApiParams{"pair": pair, "limit": string(limit), "offset": string(offset)})
}

// OrderCreate creates order
func (ex *Exmo) OrderCreate(pair string, quantity string, price string, typeOrder string) (ApiResponse, error) {
	return ex.Api_query("authenticated", "order_create", ApiParams{"pair": pair, "quantity": quantity, "price": price, "type": typeOrder})
}

// Buy creates buy order
func (ex *Exmo) Buy(pair string, quantity string, price string) (ApiResponse, error) {
	return ex.OrderCreate(pair, quantity, price, "buy")
}

// Buy creates sell order
func (ex *Exmo) Sell(pair string, quantity string, price string) (ApiResponse, error) {
	return ex.OrderCreate(pair, quantity, price, "sell")
}

// MarketBuy creates market buy-order
func (ex *Exmo) MarketBuy(pair string, quantity string) (ApiResponse, error) {
	return ex.OrderCreate(pair, quantity, "0", "market_buy")
}

// MarketBuyTotal creates market buy-order for a certain amount (quantity parameter)
func (ex *Exmo) MarketBuyTotal(pair string, quantity string) (ApiResponse, error) {
	return ex.OrderCreate(pair, quantity, "0", "market_buy_total")
}

// MarketSell creates market sell-order
func (ex *Exmo) MarketSell(pair string, quantity string) (ApiResponse, error) {
	return ex.OrderCreate(pair, quantity, "0", "market_sell")
}

// MarketSellTotal creates market sell-order for a certain amount (quantity parameter)
func (ex *Exmo) MarketSellTotal(pair string, quantity string) (ApiResponse, error) {
	return ex.OrderCreate(pair, quantity, "0", "market_sell_total")
}

// OrderCancel cancels order
func (ex *Exmo) OrderCancel(orderId string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "order_cancel", ApiParams{"order_id": orderId})
	CheckErr(err)
	return
}

// GetUserOpenOrders returns the list of user’s active orders
func (ex *Exmo) GetUserOpenOrders() (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "user_open_orders", ApiParams{})
	CheckErr(err)
	return
}

// GetUserCancelledOrders returns the list of user’s deals
// This method almost completely copies Api_query method, but it returns array of interfaces, not map
func (ex *Exmo) GetUserCancelledOrders(offset uint, limit uint) (response ApiResponse, err error) {
	if limit < 100 || limit > 1000 {
		fmt.Printf("limit param must be in range of 100-1000")
		response = nil
		err = errors.New("limit param must be in range of 100-1000")
	} else {
		response, err = ex.Api_query("authenticated", "order_cancel", ApiParams{"offset": string(offset), "limit": string(limit)})
		CheckErr(err)
		return
	}
	return
}

// GetOrderTrades returns the list of user’s cancelled orders
func (ex *Exmo) GetOrderTrades(orderId string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "order_trades", ApiParams{"order_id": orderId})
	CheckErr(err)
	return
}

// GetRequiredAmount calculating and returns the sum of buying a certain amount of currency for the particular currency pair
func (ex *Exmo) GetRequiredAmount(pair string, quantity string) (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "required_amount", ApiParams{"pair": pair, "quantity": quantity})
	CheckErr(err)
	return
}

// GetDepositAddress returns the list of addresses for cryptocurrency deposit
func (ex *Exmo) GetDepositAddress() (response ApiResponse, err error) {
	response, err = ex.Api_query("authenticated", "deposit_address", ApiParams{})
	CheckErr(err)
	return
}

/*
   Wallet API
*/

// GetWalletHistory returns history of wallet
func (ex *Exmo) GetWalletHistory(date time.Time) (response ApiResponse, err error) {
	dateUnix := date.Unix()

	dateConverted := strconv.Itoa(int(dateUnix))

	if date.IsZero() {
		response, err = ex.Api_query("authenticated", "wallet_history", ApiParams{})
	} else {
		response, err = ex.Api_query("authenticated", "wallet_history", ApiParams{"date": dateConverted})
	}

	CheckErr(err)
	return
}
