package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"net/http"
	"github.com/marstau/GoEx"
	"github.com/marstau/GoEx/huobi"
	"github.com/marstau/conver"
	"github.com/marstau/smartcooly/constant"
	"github.com/marstau/smartcooly/model"
)

// HuobiPro the exchange struct of huobipro.com
type HuobiPro struct {
	stockTypeMap     map[string]string
	tradeTypeMap     map[int]string
	recordsPeriodMap map[string]string
	minAmountMap     map[string]float64
	records          map[string][]Record
	host             string
	logger           model.Logger
	api              *huobi.HuobiPro

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewHuobiPro create an exchange struct of huobipro.com
func NewHuobiPro(opt Option) Exchange {
	hbpro := huobi.NewHuobiPro(http.DefaultClient, opt.AccessKey, opt.SecretKey, "")
	// if err != nil {
	// 	fmt.Println("get ticker failed:", err.Error())
	// 	return
	// }

	return &HuobiPro{
		stockTypeMap: map[string]string{
			"BTC/CNY": "1",
			"LTC/CNY": "2",
		},
		tradeTypeMap: map[int]string{
			1: constant.TradeTypeBuy,
			2: constant.TradeTypeSell,
			3: constant.TradeTypeBuy,
			4: constant.TradeTypeSell,
		},
		recordsPeriodMap: map[string]string{
			"M":   "001",
			"M5":  "005",
			"M15": "015",
			"M30": "030",
			"H":   "060",
			"D":   "100",
			"W":   "200",
		},
		minAmountMap: map[string]float64{
			"BTC/CNY": 0.001,
			"LTC/CNY": 0.01,
		},
		records: make(map[string][]Record),
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		api:  hbpro,

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}




// Log print something to console
func (e *HuobiPro) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *HuobiPro) GetType() string {
	// return e.api.Type
	return ""
}

// GetName get the name of this exchange
func (e *HuobiPro) GetName() string {
	// return e.api.Name
	name, err := e.api.GetAccountId()
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetName() error, ", err)
		return ""
	}
	return name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *HuobiPro) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *HuobiPro) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *HuobiPro) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// func (e *HuobiPro) getAuthJSON(url string, params []string, optionals ...string) (json *simplejson.Json, err error) {
// 	e.lastTimes++
// 	params = append(params, []string{
// 		"access_key=" + e.api.AccessKey,
// 		"secret_key=" + e.api.SecretKey,
// 		fmt.Sprint("created=", time.Now().Unix()),
// 	}...)
// 	sort.Strings(params)
// 	params = append(params, "sign="+signMd5(params))
// 	resp, err := post(url, append(params, optionals...))
// 	if err != nil {
// 		return
// 	}
// 	return simplejson.NewJson(resp)
// }

// GetAccount get the account detail of this exchange
func (e *HuobiPro) GetAccount() interface{} {
	// params := []string{
	// 	"method=get_account_info",
	// }

	// json, err := e.getAuthJSON(e.host, params, "market=cny")
	// if err != nil {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
	// 	return false
	// }

	accountid, err := e.api.GetAccountId()
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return nil
	}
	
	// return map[string]float64{
	// 	"CNY":       conver.Float64Must(json.Get("available_cny_display").Interface()),
	// 	"FrozenCNY": conver.Float64Must(json.Get("frozen_cny_display").Interface()),
	// 	"BTC":       conver.Float64Must(json.Get("available_btc_display").Interface()),
	// 	"FrozenBTC": conver.Float64Must(json.Get("frozen_btc_display").Interface()),
	// 	"LTC":       conver.Float64Must(json.Get("available_ltc_display").Interface()),
	// 	"FrozenLTC": conver.Float64Must(json.Get("frozen_ltc_display").Interface()),
	// }
	return accountid
}


// Trade place an order
func (e *HuobiPro) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	// stockType = strings.ToUpper(stockType)
	// tradeType = strings.ToUpper(tradeType)
	// price := conver.Float64Must(_price)
	// amount := conver.Float64Must(_amount)
	// if _, ok := e.stockTypeMap[stockType]; !ok {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized stockType: ", stockType)
	// 	return false
	// }
	// switch tradeType {
	// case constant.TradeTypeBuy:
	// 	return e.buy(stockType, price, amount, msgs...)
	// case constant.TradeTypeSell:
	// 	return e.sell(stockType, price, amount, msgs...)
	// default:
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized tradeType: ", tradeType)
	// 	return false
	// }
	return nil
}

func (e *HuobiPro) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	// params := []string{
	// 	"coin_type=" + e.stockTypeMap[stockType],
	// 	fmt.Sprintf("amount=%v", amount),
	// }
	// methodParam := "method=buy_market"
	// if price > 0 {
	// 	methodParam = "method=buy"
	// 	params = append(params, fmt.Sprintf("price=%v", price))
	// }
	// params = append(params, methodParam)
	// json, err := e.getAuthJSON(e.host, params, "market=cny")
	// if err != nil {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
	// 	return false
	// }
	// if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", strings.TrimSpace(json.Get("msg").MustString()))
	// 	return false
	// }
	// e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	// return fmt.Sprint(json.Get("id").Interface())
	return nil
}

func (e *HuobiPro) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	// params := []string{
	// 	"coin_type=" + e.stockTypeMap[stockType],
	// 	fmt.Sprintf("amount=%v", amount),
	// }
	// methodParam := "method=sell_market"
	// if price > 0 {
	// 	methodParam = "method=sell"
	// 	params = append(params, fmt.Sprintf("price=%v", price))
	// }
	// params = append(params, methodParam)
	// json, err := e.getAuthJSON(e.host, params, "market=cny")
	// if err != nil {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
	// 	return false
	// }
	// if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", strings.TrimSpace(json.Get("msg").MustString()))
	// 	return false
	// }
	// e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	// return fmt.Sprint(json.Get("id").Interface())
	return nil
}

// GetOrder get details of an order
func (e *HuobiPro) GetOrder(stockType, id string) interface{} {
	// stockType = strings.ToUpper(stockType)
	// if _, ok := e.stockTypeMap[stockType]; !ok {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
	// 	return false
	// }
	// params := []string{
	// 	"method=order_info",
	// 	"coin_type=" + e.stockTypeMap[stockType],
	// 	"id=" + id,
	// }
	// json, err := e.getAuthJSON(e.host, params, "market=cny")
	// if err != nil {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
	// 	return false
	// }
	// if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", strings.TrimSpace(json.Get("msg").MustString()))
	// 	return false
	// }
	// return Order{
	// 	ID:         fmt.Sprint(json.Get("id").Interface()),
	// 	Price:      conver.Float64Must(json.Get("order_price").Interface()),
	// 	Amount:     conver.Float64Must(json.Get("order_amount").Interface()),
	// 	DealAmount: conver.Float64Must(json.Get("processed_amount").Interface()),
	// 	TradeType:  e.tradeTypeMap[json.Get("type").MustInt()],
	// 	StockType:  stockType,
	// }
	return nil
}

// GetOrders get all unfilled orders
func (e *HuobiPro) GetOrders(stockType string) interface{} {
	// stockType = strings.ToUpper(stockType)
	// orders := []Order{}
	// if _, ok := e.stockTypeMap[stockType]; !ok {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
	// 	return false
	// }
	// params := []string{
	// 	"method=get_orders",
	// 	"coin_type=" + e.stockTypeMap[stockType],
	// }
	// json, err := e.getAuthJSON(e.host+"order_info.do", params)
	// if err != nil {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
	// 	return false
	// }
	// if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", strings.TrimSpace(json.Get("msg").MustString()))
	// 	return false
	// }
	// count := len(json.MustArray())
	// for i := 0; i < count; i++ {
	// 	orderJSON := json.GetIndex(i)
	// 	orders = append(orders, Order{
	// 		ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
	// 		Price:      conver.Float64Must(orderJSON.Get("order_price").Interface()),
	// 		Amount:     conver.Float64Must(orderJSON.Get("order_amount").Interface()),
	// 		DealAmount: conver.Float64Must(orderJSON.Get("processed_amount").Interface()),
	// 		TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustInt()],
	// 		StockType:  stockType,
	// 	})
	// }
	// return orders
	return nil
}

// GetTrades get all filled orders recently
func (e *HuobiPro) GetTrades(stockType string) interface{} {
	// stockType = strings.ToUpper(stockType)
	// orders := []Order{}
	// if _, ok := e.stockTypeMap[stockType]; !ok {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
	// 	return false
	// }
	// params := []string{
	// 	"method=get_new_deal_orders",
	// 	"coin_type=" + e.stockTypeMap[stockType],
	// }
	// json, err := e.getAuthJSON(e.host+"order_history.do", params)
	// if err != nil {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", err)
	// 	return false
	// }
	// if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", strings.TrimSpace(json.Get("msg").MustString()))
	// 	return false
	// }
	// count := len(json.MustArray())
	// for i := 0; i < count; i++ {
	// 	orderJSON := json.GetIndex(i)
	// 	orders = append(orders, Order{
	// 		ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
	// 		Price:      conver.Float64Must(orderJSON.Get("order_price").Interface()),
	// 		Amount:     conver.Float64Must(orderJSON.Get("order_amount").Interface()),
	// 		DealAmount: conver.Float64Must(orderJSON.Get("processed_amount").Interface()),
	// 		TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustInt()],
	// 		StockType:  stockType,
	// 	})
	// }
	// return orders
	return nil
}

// CancelOrder cancel an order
func (e *HuobiPro) CancelOrder(order Order) bool {
	// params := []string{
	// 	"method=cancel_order",
	// 	"coin_type=" + e.stockTypeMap[order.StockType],
	// 	"id=" + order.ID,
	// }
	// json, err := e.getAuthJSON(e.host, params, "market=cny")
	// if err != nil {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
	// 	return false
	// }
	// if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
	// 	e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", strings.TrimSpace(json.Get("msg").MustString()))
	// 	return false
	// }
	// if json.Get("result").MustString() == "success" {
	// 	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	// 	return true
	// }
	// e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", json.Get("msg").Interface())
	// return false
	return false
}

// // getTicker get market ticker & depth
// func (e *HuobiPro) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
// 	stockType = strings.ToUpper(stockType)
// 	if _, ok := e.stockTypeMap[stockType]; !ok {
// 		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
// 		return
// 	}
// 	size := 20
// 	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
// 		size = conver.IntMust(sizes[0])
// 	}
// 	resp, err := get(fmt.Sprintf("http://api.huobi.com/staticmarket/depth_%v_%v.js", strings.ToLower(strings.TrimSuffix(stockType, "/CNY")), size))
// 	if err != nil {
// 		err = fmt.Errorf("GetTicker() error, %+v", err)
// 		return
// 	}
// 	json, err := simplejson.NewJson(resp)
// 	if err != nil {
// 		err = fmt.Errorf("GetTicker() error, %+v", err)
// 		return
// 	}
// 	depthsJSON := json.Get("bids")
// 	for i := 0; i < len(depthsJSON.MustArray()); i++ {
// 		depthJSON := depthsJSON.GetIndex(i)
// 		ticker.Bids = append(ticker.Bids, OrderBook{
// 			Price:  depthJSON.GetIndex(0).MustFloat64(),
// 			Amount: depthJSON.GetIndex(1).MustFloat64(),
// 		})
// 	}
// 	depthsJSON = json.Get("asks")
// 	for i := 0; i < len(depthsJSON.MustArray()); i++ {
// 		depthJSON := depthsJSON.GetIndex(i)
// 		ticker.Asks = append(ticker.Asks, OrderBook{
// 			Price:  depthJSON.GetIndex(0).MustFloat64(),
// 			Amount: depthJSON.GetIndex(1).MustFloat64(),
// 		})
// 	}
// 	if len(ticker.Bids) < 1 || len(ticker.Asks) < 1 {
// 		err = fmt.Errorf("GetTicker() error, can not get enough Bids or Asks")
// 		return
// 	}
// 	ticker.Buy = ticker.Bids[0].Price
// 	ticker.Sell = ticker.Asks[0].Price
// 	ticker.Mid = (ticker.Buy + ticker.Sell) / 2
// 	return
// }

// GetTicker get market ticker & depth
func (e *HuobiPro) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.api.GetTicker(goex.BTC_USDT)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *HuobiPro) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized stockType: ", stockType)
		return false
	}
	if _, ok := e.recordsPeriodMap[period]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return false
	}
	size := 200
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("http://api.huobi.com/staticmarket/%v_kline_%v_json.js", strings.ToLower(strings.TrimSuffix(stockType, "/CNY")), e.recordsPeriodMap[period]))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	timeLast := int64(0)
	if len(e.records[period]) > 0 {
		timeLast = e.records[period][len(e.records[period])-1].Time
	}
	recordsNew := []Record{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		t, _ := time.Parse("20060102150405000", recordJSON.GetIndex(0).MustString("19700101000000000"))
		recordTime := t.Unix()
		if recordTime > timeLast {
			recordsNew = append([]Record{{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			e.records[period][len(e.records[period])-1] = Record{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}
		} else {
			break
		}
	}
	e.records[period] = append(e.records[period], recordsNew...)
	if len(e.records[period]) > size {
		e.records[period] = e.records[period][len(e.records[period])-size : len(e.records[period])]
	}
	return e.records[period]
}

func (e *HuobiPro) ExchangeRate(count string,stockType string) string {
	return ""
}
