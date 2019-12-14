## Golang lib for EXMO.com (exmo.me) cryptocurrency exchange

[Usage](https://github.com/VadimInshakov/exmo#usage)      
[Testing](https://github.com/VadimInshakov/exmo#testing)     
[API](https://github.com/VadimInshakov/exmo#api)

<br/>

### **Usage**                          

Import package into your code:                      

    import "github.com/vadiminshakov/exmo"

Call fabric function for api instance:            

    var api = exmo.Api(key, secret)

*(you can find key and secret in your profile settings)*  
  
Now you can use api features, for example:  
```golang
package main
    import (
        "github.com/vadiminshakov/exmo" 
        ...
        )   
  func main(){
    var api = exmo.Api("K-92fds9df9ew0sfg9df9sf", "S-293r9dfsjvnef3n31lmr")
    
    // Getting information about user's account
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
    
    // Buy BTC for RUB
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
  }
 ```                                                                 
   
<br/>

### **Testing**
___

<br/>

Set environment variables:
    
    export EXMO_PUBLIC="your public key"
    export EXMO_SECRET="your secret key"

<br/>

**Test specific method:**

    go test -run <method name>

<br/>

**Example:**

    go test -run GetPairSettings

<br/>

**Or run all tests:**
 
_(**ATTENTION!** Some test tasks will create buy and sell orders using your account)_

    go test

                                         
<br/>

### **API**

---

You can check the [official EXMO API](https://exmo.me/ru/api) for details

<br/>


**GetTrades(arg string)**

_List of the deals on currency pairs_

**arg** - one or various currency pairs separated by commas (example: BTC_USD,BTC_EUR)
```golang
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
``` 
  	
 <br>   
 
 **GetOrderBook(pair string, limit int)**
 
 _The book of current orders on the currency pair_
 
 **pair** - one or various currency pairs separated by commas (example: BTC_USD,BTC_EUR)
 
 **limit** - the number of returned deals (default: 100, мmaximum: 10 000)
 
 ```golang
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
```

 <br>   
     	
**Ticker()**

_Statistics on prices and volume of trades by currency pairs_

```golang
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
  ```
  	
<br>   
    
**GetPairSettings()**

_Currency pairs settings_

```golang
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
```

<br>   

**GetCurrency()**

_Currencies list_
```golang
    resultCurrency, err := api.GetCurrency()
    	if err != nil {
    		fmt.Printf("api error: %s\n", err)
    	} else {
    		fmt.Println("\nCurrencies:")
    		for _, pair := range resultCurrency {
    			fmt.Println(pair)
    		}
    	}
```

<br>   

**GetUserInfo()**

_Getting information about user's account_

```golang
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
``` 
   	
   <br>
    	
**GetUserTrades(pair string)**

_Getting the list of user’s deals_

**pair** - one or various currency pairs separated by commas (example: BTC_USD,BTC_EUR)
**limit** - limit the number of displayed positions (default: 100, max: 1000)

```golang
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
```
    	
<br>

**Buy(pair string, quantity string, price string)**

_Creation of an order to buy the currency_

**pair** - currency pair

**quantity** - quantity for the order

**price** - price for the order

```golang
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
 ```
   	
 <br>
 
**MarketBuy(pair string, quantity string)**

_Creation of an order to buy the currency at a market price_

**pair** - currency pair

**quantity** - quantity for the order

```golang
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
 ```
   	
<br>

**Sell(pair string, quantity string, price string)**

_Creation of an order to sell the currency_

**pair** - currency pair

**quantity** - quantity for the order

**price** - price for the order

```golang
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
```
    	
<br>

**MarketSell(pair string, quantity string)**

_Creation of an order to sell the currency at a market price_

**pair** - currency pair

**quantity** - quantity for the order

```golang
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
```
	
<br>

**OrderCancel(orderId string)**

_Cancels the open order_

**orderId** - id of the order to cancel
```golang
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
```
    	
<br>

**GetUserOpenOrders()**

_Getting the list of user’s active orders_

```golang
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
  ```
  	
<br>

**GetUserCancelledOrders(offset uint, limit uint)**

_Getting the list of user’s cancelled orders_

**offset** - last deal offset (default: 0)

**limit** - the number of returned deals (default: 100, мmaximum: 10 000)

```golang
    resultUserCancelledOrders, err := api.GetUserCancelledOrders(0, 100)
    	if err != nil {
    		fmt.Errorf("api error: %s\n", err)
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
  ```
  	
**GetOrderTrades(orderId string)**

_Getting the history of deals with the order_

**orderId** - order identifier

```golang
    resultOrderTrades, err := api.GetOrderTrades(orderId)
    	if err != nil {
    		fmt.Errorf("api error: %s\n", err)
    	} else {
    		for k, v := range resultOrderTrades {
    			fmt.Println(k, v)
    		}
    	}
```

<br>

**GetRequiredAmount(pair string, quantity string)**

_Calculating the sum of buying a certain amount of currency for the particular currency pair_
    
**pair** - currency pair

**quantity** - quantity to buy

```golang
    resultRequiredAmount, err := api.GetRequiredAmount("BTC_RUB", "0.01")
    	if err != nil {
    		fmt.Errorf("api error: %s\n", err)
    	} else {
    		for k, v := range resultRequiredAmount {
    			fmt.Println(k, v)
    		}
    	}
```

<br>

**GetDepositAddress()**

_Getting the list of addresses for cryptocurrency deposit_

```golang
    resultDepositAddress, err := api.GetDepositAddress()
    	if err != nil {
    		fmt.Errorf("api error: %s\n", err)
    	} else {
    		for k, v := range resultDepositAddress {
    			fmt.Println(k, v)
    		}
    	}
```

<br>
    	
**GetWalletHistory(date time.Time)**

_Get history of wallet_

**date** - timestamp of the day (if empty got current day)

```golang
    date := time.Date(2019, 10, 4, 0, 0, 0, 0, time.UTC)
    	subdate := 10*time.Hour
    
    	resultWalletHistory, err := api.GetWalletHistory(date.Truncate(subdate))
    
    	if err != nil {
    		fmt.Errorf("api error: %s\n", err)
    	} else {
    		for k, v := range resultWalletHistory {
    			if k == "history" {
    				fmt.Println(k, v)
    				for key, val := range v.([]interface{}) {
    					fmt.Println(key, val)
    				}
    			}
    		}
    	}
```
