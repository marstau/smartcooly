package api

import (
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/log"
)

// OKCoinCn : the exchange struct of okcoin.cn
type OKCoinCn struct {
	stockMap map[string]string
	host     string
	log      log.Logger
	option   Option
}

// NewOKCoinCn : create an exchange struct of okcoin.cn
func NewOKCoinCn(opt Option) *OKCoinCn {
	e := OKCoinCn{
		stockMap: map[string]string{"BTC": "btc", "LTC": "ltc"},
		host:     "https://www.okcoin.cn/api/v1/",
		log:      log.New(opt.Type),
		option:   opt,
	}
	if _, ok := e.stockMap[e.option.MainStock]; !ok {
		e.option.MainStock = "BTC"
	}
	return &e
}

// Log : print something to console
func (e *OKCoinCn) Log(msgs ...interface{}) {
	e.log.Do("info", 0.0, 0.0, msgs...)
}

// GetMainStock : get the MainStock of this exchange
func (e *OKCoinCn) GetMainStock() string {
	return e.option.MainStock
}

// SetMainStock : set the MainStock of this exchange
func (e *OKCoinCn) SetMainStock(stock string) string {
	if _, ok := e.stockMap[stock]; ok {
		e.option.MainStock = stock
	}
	return e.option.MainStock
}

// GetAccount : get the account detail of this exchange
func (e *OKCoinCn) GetAccount() interface{} {
	account := make(map[string]float64)
	params := []string{
		"api_key=" + e.option.AccessKey,
		"secret_key=" + e.option.SecretKey,
	}
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(e.host+"userinfo.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, err)
		return nil
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, err)
		return nil
	}

	if result := json.Get("result").MustBool(); !result {
		err = fmt.Errorf("GetAccount() error, the error number is %v", json.Get("error_code").MustInt())
		e.log.Do("error", 0.0, 0.0, err)
		return nil
	}
	account["Total"] = conver.Float64Must(json.GetPath("info", "funds", "asset", "total").Interface())
	account["Net"] = conver.Float64Must(json.GetPath("info", "funds", "asset", "net").Interface())
	account["Balance"] = conver.Float64Must(json.GetPath("info", "funds", "free", "cny").Interface())
	account["FrozenBalance"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "cny").Interface())
	account["BTC"] = conver.Float64Must(json.GetPath("info", "funds", "free", "btc").Interface())
	account["FrozenBTC"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "btc").Interface())
	account["LTC"] = conver.Float64Must(json.GetPath("info", "funds", "free", "ltc").Interface())
	account["FrozenLTC"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "ltc").Interface())
	account["Stocks"] = account[e.option.MainStock]
	account["FrozenStocks"] = account["Frozen"+e.option.MainStock]
	return account
}

// Buy ...
func (e *OKCoinCn) Buy(stockType string, price, amount float64, msgs ...interface{}) (id int) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "Buy() error, unrecognized stockType")
		return
	}
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[stockType] + "_cny",
	}
	typeParam := "type=buy_market"
	amountParam := fmt.Sprint("price=", amount)
	if price > 0 {
		typeParam = "type=buy"
		amountParam = fmt.Sprint("amount=", amount)
		params = append(params, fmt.Sprint("price=", price))
	}
	params = append(params, typeParam, amountParam)
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	fmt.Println(params)
	resp, err := post(e.host+"trade.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, err)
		return
	}
	if result := json.Get("result").MustBool(); !result {
		e.log.Do("error", 0.0, 0.0, "Buy() error, the error number is ", json.Get("error_code").MustInt())
		return
	}
	e.log.Do("buy", price, amount, msgs...)
	id = json.Get("order_id").MustInt()
	return
}