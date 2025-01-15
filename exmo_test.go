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
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestApi_query(t *testing.T) {

	var orderId string // global variable for testing order cancelling after buying

	// ATTENTION!
	key := os.Getenv("EXMO_PUBLIC")
	secret := os.Getenv("EXMO_SECRET")

	api := Api(key, secret)

	t.Run("GetTrades", func(t *testing.T) {
		result, err := api.GetTrades("BTC_RUB")
		require.NoError(t, err)

		for _, v := range result {
			for _, val := range v.([]interface{}) {
				for key, value := range val.(map[string]interface{}) {
					if key == "trade_id" || key == "date" {
						check, ok := value.(float64)
						require.True(t, ok)

						if check < 0 {
							t.Errorf("%s could not be less 0, got %d", key, value)
						}
					}
				}
			}
		}
	})

	t.Run("GetOrderBook", func(t *testing.T) {
		result, err := api.GetOrderBook("BTC_RUB", 200)
		require.NoError(t, err)

		for _, v := range result {
			for key, value := range v.(map[string]interface{}) {
				if key == "bid" || key == "ask" {
					for _, val := range value.([]interface{}) {
						for _, valnested := range val.([]interface{}) {
							check, err := strconv.ParseFloat(valnested.(string), 64)
							if err != nil {
								t.Errorf("Could not convert %s to float64", key)
							}
							if check < 0 {
								t.Errorf("%s could not be less 0, got %d", key, valnested)
							}
						}
					}
				} else {
					check, err := strconv.ParseFloat(value.(string), 64)
					if err != nil {
						t.Errorf("Could not convert %s to float64", key)
					}
					if check < 0 {
						t.Errorf("%s could not be less 0, got %d", key, value)
					}
				}
			}
		}

	})

	t.Run("Ticker", func(t *testing.T) {
		ticker, errTicker := api.Ticker()
		require.NoError(t, errTicker)

		for _, pairvalue := range ticker {
			for key, value := range pairvalue.(map[string]interface{}) {
				if key == "updated" {
					check, ok := value.(float64)
					if !ok {
						t.Errorf("Could not convert %s to float64", key)
					}
					if check < 0 {
						t.Errorf("%s could not be less 0, got %d", key, value)
					}
				} else {
					check, err := strconv.ParseFloat(value.(string), 64)
					if err != nil {
						t.Errorf("Could not convert %s to float64", key)
					}
					if check < 0 {
						t.Errorf("%s could not be less 0, got %d", key, value)
					}
				}
			}
		}
	})

	t.Run("GetPairSettings", func(t *testing.T) {
		resultPairSettings, err := api.GetPairSettings()
		require.NoError(t, err)

		for _, pairvalue := range resultPairSettings {
			for key, value := range pairvalue.(map[string]interface{}) {
				if reflect.TypeOf(key).Name() != "string" {
					t.Errorf("response item %#v not a string", key)
				}

				_, err := strconv.ParseFloat(value.(string), 64)
				if err != nil {
					t.Errorf("Can't cast %#v to float64, error: %s", value, err)
				}
			}
		}
	})

	t.Run("GetCurrency", func(t *testing.T) {
		result, err := api.GetCurrency()
		require.NoError(t, err)

		for _, pair := range result {
			if reflect.TypeOf(pair).Name() != "string" {
				t.Errorf("response item %#v not a string", pair)
			}
		}
	})

	t.Run("GetUserInfo", func(t *testing.T) {
		result, err := api.GetUserInfo()
		require.NoError(t, err)

		for key, value := range result {
			if key == "balances" || key == "reserved" {
				for k, v := range value.(map[string]interface{}) {
					check, err := strconv.ParseFloat(v.(string), 64)
					if err != nil {
						t.Errorf("Could not convert %s to float64", k)
					}
					if check < 0 {
						t.Errorf("%s could not be less 0, got %d", k, v)
					}
				}
			} else {
				check, ok := value.(float64)
				if !ok {
					t.Errorf("Could not convert %s to float64", key)
				}
				if check < 0 {
					t.Errorf("%s could not be less 0, got %d", key, value)
				}
			}
		}
	})

	t.Run("GetUserTrades", func(t *testing.T) {
		usertrades, err := api.GetUserTrades("BTC_RUB", 0, 1000)
		require.NoError(t, err)

		for _, val := range usertrades {
			for _, interfacevalue := range val.([]interface{}) {
				for k, v := range interfacevalue.(map[string]interface{}) {
					if k == "trade_id" || k == "date" || k == "order_id" {
						check, ok := v.(float64)
						if !ok {
							t.Errorf("Could not convert %s to float64", k)
						}
						if check < 0 {
							t.Errorf("%s could not be less 0, got %d", k, v)
						}
					} else if k == "quantity" || k == "price" || k == "amount" {
						check, err := strconv.ParseFloat(v.(string), 64)
						if err != nil {
							t.Errorf("Could not convert %s to float64", k)
						}
						if check < 0 {
							t.Errorf("%s could not be less 0, got %d", k, v)
						}
					} else {
						if reflect.TypeOf(v).Name() != "string" {
							t.Errorf("response item %s (value %#v) not a string, but %T", k, v, v)
						}
					}
				}
			}
		}
	})

	//t.Run("Buy", func(t *testing.T) {
	//	order, err := api.Buy("BTC_RUB", "0.001", "50096.72")
	//	if err != nil {
	//		t.Errorf("api error: %s\n", err)
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				val := strconv.Itoa(int(value.(float64)))
	//				orderId = val
	//				fmt.Printf("Order id: %s", orderId)
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("MarketBuy", func(t *testing.T) {
	//	order, err := api.MarketBuy("BTC_RUB", "0.001")
	//	if err != nil {
	//		t.Errorf("api error: %s\n", err)
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				fmt.Printf("Order id: %f", value.(float64))
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("Sell", func(t *testing.T) {
	//	order, err := api.Sell("BTC_RUB", "0.001", "800000")
	//	if err != nil {
	//		t.Errorf("api error: %s\n", err)
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				fmt.Printf("Order id: %f", value.(float64))
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("MarketSell", func(t *testing.T) {
	//	order, err := api.MarketSell("BTC_RUB", "0.001")
	//	if err != nil {
	//		t.Errorf("api error: %s\n", err)
	//	} else {
	//		fmt.Println("Creating order...")
	//		for key, value := range order {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//			if key == "order_id" && value != nil {
	//				fmt.Printf("Order id: %f", value.(float64))
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("OrderCancel", func(t *testing.T) {
	//	orderCancel, err := api.OrderCancel(orderId)
	//	if err != nil {
	//		t.Errorf("api error: %s\n", err)
	//	} else {
	//		fmt.Printf("\nCancel order %s \n", orderId)
	//		for key, value := range orderCancel {
	//			if key == "result" && value != true {
	//				fmt.Println("\nError")
	//			}
	//			if key == "error" && value != "" {
	//				fmt.Println(value)
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("GetUserOpenOrders", func(t *testing.T) {
	//	resultUserOpenOrders, err := api.GetUserOpenOrders()
	//	if err != nil {
	//		fmt.Errorf("api error: %s\n", err)
	//	} else {
	//		for _, v := range resultUserOpenOrders {
	//			if v != nil {
	//				for _, val := range v.([]interface{}) {
	//					for key, value := range val.(map[string]interface{}) {
	//						if key == "quantity" {
	//							check, err := strconv.ParseFloat(value.(string), 64)
	//							if err != nil {
	//								t.Errorf("Could not convert %s to float64", key)
	//							}
	//							if check < 0 {
	//								t.Errorf("%s could not be less 0, got %d", key, value)
	//							}
	//						}
	//						if key == "price" {
	//							check, err := strconv.Atoi(value.(string))
	//							if err != nil {
	//								t.Errorf("Could not convert %s to int", key)
	//							}
	//							if check < 0 {
	//								t.Errorf("%s could not be less 0, got %d", key, value)
	//							}
	//						}
	//						if key == "amount" {
	//							check, err := strconv.ParseFloat(value.(string), 64)
	//							if err != nil {
	//								t.Errorf("Could not convert %s to float64", key)
	//							}
	//							if check < 0 {
	//								t.Errorf("%s could not be less 0, got %d", key, value)
	//							}
	//						}
	//					}
	//				}
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("GetUserCancelledOrders", func(t *testing.T) {
	//	resultUserCancelledOrders, err := api.GetUserCancelledOrders(0, 100)
	//	if err != nil {
	//		fmt.Errorf("api error: %s\n", err)
	//	} else {
	//		for _, v := range resultUserCancelledOrders {
	//
	//			if v != nil {
	//				for key, value := range v.(map[string]interface{}) {
	//
	//					if key == "quantity" || key == "price" || key == "amount" {
	//						check, ok := value.(float64)
	//						if ok != true {
	//							t.Errorf("Could not convert %s to float64", key)
	//						}
	//						if check < 0 {
	//							t.Errorf("%s could not be less 0, got %d", key, value)
	//						}
	//					}
	//				}
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("GetRequiredAmount", func(t *testing.T) {
	//	resultRequiredAmount, err := api.GetRequiredAmount("BTC_RUB", "0.01")
	//	if err != nil {
	//		fmt.Errorf("api error: %s\n", err)
	//	} else {
	//		for k, v := range resultRequiredAmount {
	//			check, err := strconv.ParseFloat(v.(string), 64)
	//			if err != nil {
	//				t.Errorf("Could not convert %s to float64", k)
	//			}
	//			if check < 0 {
	//				t.Errorf("%s could not be less 0, got %d", k, v)
	//			}
	//		}
	//	}
	//})
	//
	//t.Run("GetDepositAddress", func(t *testing.T) {
	//	resultDepositAddress, err := api.GetDepositAddress()
	//	if err != nil {
	//		fmt.Errorf("api error: %s\n", err)
	//	} else {
	//		for _, v := range resultDepositAddress {
	//			_, ok := v.(string)
	//			if ok != true {
	//				t.Errorf("Could not convert %s address to string", key)
	//			}
	//		}
	//	}
	//})
	_ = orderId
}
