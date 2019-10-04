package main

import (
	"fmt"
	"math/big"
	_ "math/big"
	"strconv"
	_ "strconv"
	"time"

	"github.com/vadiminshakov/exmo"
)

func main() {

	var orderId string

	key := ""    // TODO replace with your api key from profile page
	secret := "" // TODO replace with your api secret from profile page

	api := exmo.Api(key, secret)

	resultTrades, err := api.GetTrades("BTC_RUB")
	if err != nil {
		fmt.Errorf("api error: %s\n", err)
	} else {
		for _, v := range resultTrades {
			for k, val := range v.([]interface{}) {
				tmpindex := 0
				for key, value := range val.(map[string]interface{}) {
					if tmpindex != k {
						fmt.Printf("\n\nindex: %d \n", k)
						tmpindex = k
					}
					if key == "trade_id" {
						fmt.Println(key, big.NewFloat(value.(float64)).String())
					} else if key == "date" {
						fmt.Println(key, time.Unix(int64(value.(float64)), 0))
					} else {
						fmt.Println(key, value)
					}
				}
			}
		}
	}

	resultBook, err := api.GetOrderBook("BTC_RUB", 200)
	if err != nil {
		fmt.Errorf("api error: %s\n", err)
	} else {
		for _, v := range resultBook {
			for key, value := range v.(map[string]interface{}) {
				if key == "bid" || key == "ask" {
					for _, val := range value.([]interface{}) {
						fmt.Printf("%s: ", key)
						for index, valnested := range val.([]interface{}) {
							switch index {
							case 0:
								fmt.Printf("price %s, ", valnested.(string))

							case 1:
								fmt.Printf("quantity %s, ", valnested.(string))
							case 2:
								fmt.Printf("total %s \n", valnested.(string))
							}
						}
					}
				} else {
					fmt.Println(key, value)
				}
			}

		}
	}

	ticker, err := api.Ticker()
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		for pair, pairvalue := range ticker {
			fmt.Printf("\n\n%s:\n", pair)
			for key, value := range pairvalue.(map[string]interface{}) {
				fmt.Println(key, value)
			}
		}
	}

	resultPairSettings, err := api.GetPairSettings()
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		for pair, pairvalue := range resultPairSettings {
			fmt.Printf("\n\n%s:\n", pair)
			for key, value := range pairvalue.(map[string]interface{}) {
				fmt.Println(key, value)
			}
		}
	}

	resultCurrency, err := api.GetCurrency()
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		fmt.Println("\nCurrencies:")
		for _, pair := range resultCurrency {
			fmt.Println(pair)
		}
	}

	resultUserInfo, err := api.GetUserInfo()
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		for key, value := range resultUserInfo {
			if key == "balances" {
				fmt.Println("\n-- balances:")
				for k, v := range value.(map[string]interface{}) {
					fmt.Println(k, v)
				}
			}
			if key == "reserved" {
				fmt.Println("\n-- reserved:")
				for k, v := range value.(map[string]interface{}) {
					fmt.Println(k, v)
				}
			}
		}

	}

	fmt.Printf("-------------\n")

	usertrades, err := api.GetUserTrades("BTC_RUB")
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		fmt.Println("User trades")
		for pair, val := range usertrades {
			fmt.Printf("\n\n %s", pair)
			for _, interfacevalue := range val.([]interface{}) {
				fmt.Printf("\n\n***\n")
				for k, v := range interfacevalue.(map[string]interface{}) {
					fmt.Println(k, v)
				}
			}
		}
	}

	order, err := api.Buy("BTC_RUB", "0.001", "50096")
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		fmt.Println("Creating order...")
		for key, value := range order {
			if key == "result" && value != true {
				fmt.Println("\nError")
			}
			if key == "error" && value != "" {
				fmt.Println(value)
			}
			if key == "order_id" && value != nil {
				fmt.Printf("Order id: %d\n", int(value.(float64)))
				val := strconv.Itoa(int(value.(float64)))
				orderId = val
				fmt.Printf("Order id: %s\n", orderId)
			}
		}
	}

	marketOrder, err := api.MarketBuy("BTC_RUB", "0.001")
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		fmt.Println("Creating order...")
		for key, value := range marketOrder {
			if key == "result" && value != true {
				fmt.Println("\nError")
			}
			if key == "error" && value != "" {
				fmt.Println(value)
			}
			if key == "order_id" && value != nil {
				val := strconv.Itoa(int(value.(float64)))
				orderId = val
				fmt.Printf("Order id: %s", orderId)
			}
		}
	}

	orderSell, err := api.Sell("BTC_RUB", "0.001", "800000")
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		fmt.Println("Creating order...")
		for key, value := range orderSell {
			if key == "result" && value != true {
				fmt.Println("\nError")
			}
			if key == "error" && value != "" {
				fmt.Println(value)
			}
			if key == "order_id" && value != nil {
				val := strconv.Itoa(int(value.(float64)))
				orderId = val
				fmt.Printf("Order id: %f", orderId)
			}
		}
	}

	orderSellMarket, err := api.MarketSell("BTC_RUB", "0.0005")
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		fmt.Println("Creating order...")
		for key, value := range orderSellMarket {
			if key == "result" && value != true {
				fmt.Println("\nError")
			}
			if key == "error" && value != "" {
				fmt.Println(value)
			}
			if key == "order_id" && value != nil {
				val := strconv.Itoa(int(value.(float64)))
				orderId = val
				fmt.Printf("Order id: %s", orderId)
			}
		}
	}

	orderCancel, err := api.OrderCancel(orderId)
	if err != nil {
		fmt.Printf("api error: %s\n", err)
	} else {
		fmt.Printf("\nCancel order %s \n", orderId)
		for key, value := range orderCancel {
			if key == "result" && value != true {
				fmt.Println("\nError")
			}
			if key == "error" && value != "" {
				fmt.Println(value)
			}
		}
	}

	resultUserOpenOrders, err := api.GetUserOpenOrders()
	if err != nil {
		fmt.Errorf("api error: %s\n", err)
	} else {
		for _, v := range resultUserOpenOrders {
			for _, val := range v.([]interface{}) {
				for key, value := range val.(map[string]interface{}) {
					fmt.Println(key, value)
				}
			}
		}
	}

	resultUserCancelledOrders, err := api.GetUserCancelledOrders(0, 100)
	if err != nil {
		fmt.Errorf("api error: %s\n", err.Error())
	} else {
		for _, v := range resultUserCancelledOrders {
			for key, val := range v.(map[string]interface{}) {
				if key == "pair" {
					fmt.Printf("\n%s\n", val)
				} else {
					fmt.Println(key, val)
				}
			}
		}
	}

	time.Sleep(10000 * time.Millisecond)

	resultOrderTrades, err := api.GetOrderTrades(orderId)
	if err != nil {
		fmt.Errorf("api error: %s\n", err.Error())
	} else {
		for k, v := range resultOrderTrades {
			fmt.Println(k, v)

		}
	}

	resultRequiredAmount, errRequiredAmount := api.GetRequiredAmount("BTC_RUB", "0.01")
	if errRequiredAmount != nil {
		fmt.Errorf("api error: %s\n", errRequiredAmount.Error())
	} else {
		for k, v := range resultRequiredAmount {
			fmt.Println(k, v)
		}
	}

	resultDepositAddress, errDepositAddress := api.GetDepositAddress()
	if errDepositAddress != nil {
		fmt.Errorf("api error: %s\n", errDepositAddress.Error())
	} else {
		for k, v := range resultDepositAddress {
			fmt.Println(k, v)
		}
	}

	/*
	   WALLET API
	*/

	resultWalletHistory, errWalletHistory := api.GetWalletHistory(time.Now())
	if errWalletHistory != nil {
		fmt.Errorf("api error: %s\n", errWalletHistory.Error())
	} else {
		for k, v := range resultWalletHistory {
			if k == "history" {
				for key, val := range v.([]interface{}) {
					fmt.Println(key, val)
				}
			}
			fmt.Println(k, v)
		}
	}
}
