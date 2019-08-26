**Usage**                          

Import package into your code:                      

    import "github.com/vadiminshakov/exmo"

Call fabric function for api instance:            

    api := exmo.Api(key, secret)

*(you can find key and secret in your profile settings)*  
  
Now you can use api features, for example:      

    result, err := api.GetUserInfo()
    ...
    order, errOrder := api.Buy("BTC_RUB", "0.001", "664096.72")
    