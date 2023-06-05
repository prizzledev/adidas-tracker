package cli

import (
	"AdidasTracker/src/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls "github.com/bogdanfinn/tls-client"
)

type Tracker struct {
	OrderId string
	Email   string
	Proxy   string
	Invoice string

	client     tls.HttpClient
	trackingId string
	region     string

	Result      TrackingResponse
	InvoiceList []InvoiceListResponse
	Error       error
}

type SubmitStuct struct {
	Id string `json:"id"`
}

type Shipment struct {
	ShipmentNo           string `json:"shipmentNo"`
	TrackingURL          string `json:"trackingUrl"`
	TrackingNo           string `json:"trackingNo"`
	CarrierCode          string `json:"carrierCode"`
	SCAC                 string `json:"SCAC"`
	ShipmentDate         string `json:"shipmentDate"`
	ExpectedDeliveryDate string `json:"expectedDeliveryDate"`
	Status               string `json:"status"`
}

type TrackingResponse struct {
	AdjustedMerchandizeTotal float64 `json:"adjustedMerchandizeTotal"`
	MerchandizeTotal         float64 `json:"merchandizeTotal"`
	TotalAmount              float64 `json:"totalAmount"`
	TotalTax                 float64 `json:"totalTax"`
	IsCancellable            bool    `json:"isCancellable"`
	IsReturnable             bool    `json:"isReturnable"`
	IsExchangeable           bool    `json:"isExchangeable"`
	DeliveryStatus           string  `json:"deliveryStatus"`
	Status                   string  `json:"status"`
	OrderTrackData           string  `json:"orderTrackData"`
	ID                       string  `json:"id"`
	Currency                 string  `json:"currency"`
	IsPaid                   bool    `json:"isPaid"`
	OrderNo                  string  `json:"orderNo"`
	CreationDate             string  `json:"creationDate"`
	Type                     string  `json:"type"`
	ProductLineItems         []struct {
		ID                              string  `json:"id"`
		LineKey                         string  `json:"lineKey"`
		LineItemPosition                float64 `json:"lineItemPosition"`
		ReturnableQuantity              float64 `json:"returnableQuantity"`
		ProductID                       string  `json:"productId"`
		ProductName                     string  `json:"productName"`
		ProductType                     string  `json:"productType"`
		Color                           string  `json:"color"`
		Size                            string  `json:"size"`
		LiteralSize                     string  `json:"literalSize"`
		ModelNumber                     string  `json:"modelNumber"`
		ArticleNumber                   string  `json:"articleNumber"`
		Status                          string  `json:"status"`
		StatusDate                      string  `json:"statusDate"`
		IsCancellable                   bool    `json:"isCancellable"`
		Price                           float64 `json:"price"`
		AdjustedPrice                   float64 `json:"adjustedPrice"`
		AdjustedPriceAfterOrderDiscount float64 `json:"adjustedPriceAfterOrderDiscount"`
		IsGiftCard                      bool    `json:"isGiftCard"`
		IsHype                          bool    `json:"isHype"`
		IsFlash                         bool    `json:"isFlash"`
		IsMTBR                          bool    `json:"isMTBR"`
		IsMCBP                          bool    `json:"isMCBP"`
		IsRecyclable                    bool    `json:"isRecyclable"`
		IsReturnable                    bool    `json:"isReturnable"`
		IsInline                        bool    `json:"isInline"`
		IsPreOrder                      bool    `json:"isPreOrder"`
		IsBackOrder                     bool    `json:"isBackOrder"`
		IsKnownBackOrder                bool    `json:"isKnownBackOrder"`
		Gender                          string  `json:"gender"`
		ExpectedShipmentDate            string  `json:"expectedShipmentDate"`
		ExpectedDeliveryDate            string  `json:"expectedDeliveryDate"`
		EstimatedDeliveryPeriod         struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"estimatedDeliveryPeriod"`
		CancellationReason  interface{} `json:"cancellationReason"`
		CustomizationRecipe interface{} `json:"customizationRecipe"`
		ImageURL            string      `json:"imageUrl"`
		IsPersonalized      bool        `json:"isPersonalized"`
		ShipmentNo          string      `json:"shipmentNo"`
		ActualShipmentDate  string      `json:"actualShipmentDate"`
		ShipmentStatus      string      `json:"shipmentStatus"`
		IsBonusProduct      bool        `json:"isBonusProduct"`
		PriceAdjustments    interface{} `json:"priceAdjustments"`
	} `json:"productLineItems"`
	Shipments      []Shipment `json:"shipments"`
	PaymentMethods []struct {
		PaymentType string `json:"paymentType"`
	} `json:"paymentMethods"`
	OrderSource struct {
		OmniHub   bool `json:"omniHub"`
		BasketAPI bool `json:"basketApi"`
	} `json:"orderSource"`
	CustomerEmail  string `json:"customerEmail"`
	BillingAddress struct {
		Email        string `json:"email"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Phone        string `json:"phone"`
		AddressLine1 string `json:"addressLine1"`
		AddressLine2 string `json:"addressLine2"`
		City         string `json:"city"`
		State        string `json:"state"`
		Country      string `json:"country"`
		PostalCode   string `json:"postalCode"`
	} `json:"billingAddress"`
	CustomerEUCI     string `json:"customerEUCI"`
	IsCashOnDelivery bool   `json:"isCashOnDelivery"`
	TimeStamp        int64  `json:"timeStamp"`
	IsExchange       bool   `json:"isExchange"`
	IsGoodWillCredit bool   `json:"isGoodWillCredit"`
	Shipping         struct {
		DeliveryInformation []struct {
			CarrierName        string `json:"carrierName"`
			ShippingMethodName string `json:"shippingMethodName"`
		} `json:"deliveryInformation"`
		EstimatedDeliveryDates []string `json:"estimatedDeliveryDates"`
		ShippingAddress        struct {
			Email        string `json:"email"`
			FirstName    string `json:"firstName"`
			LastName     string `json:"lastName"`
			Phone        string `json:"phone"`
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			City         string `json:"city"`
			State        string `json:"state"`
			Country      string `json:"country"`
			PostalCode   string `json:"postalCode"`
		} `json:"shippingAddress"`
		Prices struct {
			ShippingTotal         float64 `json:"shippingTotal"`
			AdjustedShippingTotal float64 `json:"adjustedShippingTotal"`
			ShippingTotalTax      float64 `json:"shippingTotalTax"`
			ShipmentPrices        []struct {
				Reference     string  `json:"reference"`
				AdjustedTotal float64 `json:"adjustedTotal"`
				Total         float64 `json:"total"`
			} `json:"shipmentPrices"`
		} `json:"prices"`
	} `json:"shipping"`
	Exchanges struct {
		Items []interface{} `json:"items"`
	} `json:"exchanges"`
	ExternalOrderNo string `json:"externalOrderNo"`
	DisplayStatus   string `json:"displayStatus"`
	InvoiceListID   string `json:"invoiceListId"`
}

type InvoiceListResponse struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	InvoicedOn string  `json:"invoicedOn"`
	OrderNo    string  `json:"orderNo"`
	Key        string  `json:"key"`
	Amount     float64 `json:"amount"`
}

func (t *Tracker) Track() {
	logger.Info("Adidas-Tracker", "Tracking order "+t.OrderId+" for "+t.Email)

	//parse the region from the order id (2nd and 3rd characters)
	t.region = strings.ToLower(t.OrderId[1:3])

	if strings.ToLower(t.region) == "d9" {
		t.region = "com"
	}

	options := []tls.HttpClientOption{
		tls.WithClientProfile(tls.Chrome_111),
		tls.WithTimeoutMilliseconds(20000),
		tls.WithCookieJar(tls.NewCookieJar()),
		//tls.WithCharlesProxy("localhost", "8889"),
	}

	if t.Proxy != "" {
		logger.Info("Adidas-Tracker", "Using proxy "+t.Proxy)
		proxyParts := strings.Split(t.Proxy, ":")
		options = append(options, tls.WithProxyUrl(fmt.Sprintf("http://%s:%s@%s:%s", proxyParts[2], proxyParts[3], proxyParts[0], proxyParts[1])))
	}

	client, err := tls.NewHttpClient(nil, options...)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error creating client: "+err.Error())
		t.Error = err
		return
	}

	t.client = client

	// Get tracking page
	// This is needed to 'gen' some cookies to avoid akamai blocks

	err = t.getTrackingPage()
	if err != nil {
		t.Error = err
		return
	}

	logger.Info("Adidas-Tracker", "Got tracking page")

	// Submit orderid and email

	err = t.submitTracking()
	if err != nil {
		t.Error = err
		return
	}

	logger.Info("Adidas-Tracker", "TrackingId retrieved")

	// Get tracking data

	err = t.getTrackingData()
	if err != nil {
		t.Error = err
		return
	}

	logger.Success("Adidas-Tracker", "Tracking data retrieved")

	if t.Result.Status == "COMPLETED" || t.Result.Status == "DELIVERED" && strings.ToLower(t.Invoice) == "true" {
		logger.Info("Adidas-Tracker", "Getting invoice list")

		err = t.getInvoiceList()
		if err != nil {
			t.Error = err
			return
		}

		logger.Info("Adidas-Tracker", "Invoice list retrieved")

		for i, invoice := range t.InvoiceList {
			logger.Info("Adidas-Tracker", "Getting invoice "+invoice.ID)

			err = t.getInvoice(invoice.ID, i)
			if err != nil {
				t.Error = err
				return
			}

			logger.Success("Adidas-Tracker", "Invoice "+invoice.OrderNo+" saved")
		}

		logger.Success("Adidas-Tracker", "Invoices retrieved")
	} else {
		logger.Warning("Adidas-Tracker", "Invcoice list not available")
	}

	return
}

func (t *Tracker) getTrackingPage() error {
	homeReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.adidas.%s/", t.region), nil)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error creating request: "+err.Error())
		return nil
	}

	homeReq.Header = http.Header{
		"sec-ch-ua":                 []string{`"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`},
		"sec-ch-ua-mobile":          []string{"?0"},
		"sec-ch-ua-platform":        []string{`"macOS"`},
		"upgrade-insecure-requests": []string{"1"},
		"user-agent":                []string{`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36`},
		"accept":                    []string{`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`},
		"sec-fetch-site":            []string{"none"},
		"sec-fetch-mode":            []string{"navigate"},
		"sec-fetch-user":            []string{"?1"},
		"sec-fetch-dest":            []string{"document"},
		"accept-encoding":           []string{"gzip, deflate, br"},
		"accept-language":           []string{"en-GB,en;q=0.9"},

		http.HeaderOrderKey: []string{
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
		},
	}

	homeResp, err := t.client.Do(homeReq)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error getting tracking page: "+err.Error())
		return errors.New("Error getting tracking page: " + err.Error())
	}

	if homeResp.StatusCode != http.StatusOK {
		logger.Error("Adidas-Tracker", "Error getting tracking page: "+homeResp.Status)
		return errors.New("Error getting tracking page: " + homeResp.Status)
	}

	return nil
}

func (t *Tracker) submitTracking() error {
	submitorderBody := map[string]interface{}{
		"email":     t.Email,
		"recaptcha": "03AL8dmw__DGeBroNMagCZnBO4iGxble40onqADtLKMshydW2mIfMt6OnkYwMIus0n5q6eyi1PJkcIEVmZA576bSFIBv2upMYRlwKJZ9akbUts38yO9ISAxbPYb1fF_sTzHvSVlMQc2Jlf-Z0g-SrqKzkWzNLLFB0-TgCfiv5n12iMHwRZbEFvHSbadUhdTueiEVaDXTzXAWJeBJjaqunhzTVsOM9t8xSU_TEpmOg_hL-H_uw4_GwDUlDk4CyF1qNlJ8eFIryk4YEXQyxEhLc9jF-q8bR0ozUAkun04lU18_5_g7vCcP0eG5lyVWbR8z1EiOSXd2DTtAivxE_wMq6fIvm4g6q-EFoW7p3jzV6AYDVq59T1IO-XO2p2MvuS_qgd1IUQPz8yUS7SnChHEx9membGRDS8L29uDGg_vB6__hgywrGvxTgehrVd4BsPoLGW_bY20U8shfqRKpm1wxbf6DN9bs-ku53oNqegs18EEUmwuAXbLwWAmREFh4O5bWfLf59PVw7uIBM7xq3V99KzyPl5mqSl2HM1yaCXxvwd6ixBmuvsDebtbgS39Mx22CAIgGq_6bh_nEIKetUbB2xDZnihbRbEQvP35g",
		"orderNo":   t.OrderId,
		"returnHub": false,
	}

	submitorderPayload, err := json.Marshal(submitorderBody)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error marshalling submitorder payload: "+err.Error())
		return errors.New("Error marshalling submitorder payload: " + err.Error())
	}

	submitOrderReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://www.adidas.%s/api/orders/search", t.region), strings.NewReader(string(submitorderPayload)))
	if err != nil {
		logger.Error("Adidas-Tracker", "Error creating submitorder request: "+err.Error())
		return errors.New("Error creating submitorder request: " + err.Error())
	}

	submitOrderReq.Header = http.Header{
		"x-instana-t":        []string{"7b55e312cb11596c"},
		"sec-ch-ua":          []string{`"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`},
		"sec-ch-ua-mobile":   []string{"?0"},
		"x-instana-l":        []string{"1,correlationType=web;correlationId=7b55e312cb11596c"},
		"user-agent":         []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"content-type":       []string{"application/json"},
		"x-instana-s":        []string{"7b55e312cb11596c"},
		"sec-ch-ua-platform": []string{`"macOS"`},
		"accept":             []string{"*/*"},
		"origin":             []string{"https://www.adidas.de"},
		"sec-fetch-site":     []string{"same-origin"},
		"sec-fetch-mode":     []string{"cors"},
		"sec-fetch-dest":     []string{"empty"},
		"referer":            []string{"https://www.adidas.de/bestellverfolgung"},
		"accept-encoding":    []string{"gzip, deflate, br"},
		"accept-language":    []string{"en-GB,en;q=0.9"},

		http.HeaderOrderKey: []string{
			"content-length",
			"x-instana-t",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"x-instana-l",
			"user-agent",
			"content-type",
			"x-instana-s",
			"sec-ch-ua-platform",
			"accept",
			"origin",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
		},
	}

	submitOrderResp, err := t.client.Do(submitOrderReq)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error submitting orderId: "+err.Error())
		return errors.New("Error submitting orderId: " + err.Error())
	}

	defer submitOrderResp.Body.Close()

	if submitOrderResp.StatusCode != http.StatusOK {
		logger.Error("Adidas-Tracker", "Error submitting orderId: "+submitOrderResp.Status)
		return errors.New("Error submitting orderId: " + submitOrderResp.Status)
	}

	var submitOrderRespBody SubmitStuct

	submitOrderBody, err := ioutil.ReadAll(submitOrderResp.Body)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error reading submitorder response body: "+err.Error())
		return errors.New("Error reading submitorder response body: " + err.Error())
	}

	err = json.Unmarshal(submitOrderBody, &submitOrderRespBody)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error unmarshalling submitorder response body: "+err.Error())
		return errors.New("Error unmarshalling submitorder response body: " + err.Error())
	}

	t.trackingId = submitOrderRespBody.Id

	return nil
}

func (t *Tracker) getTrackingData() error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.adidas.%s/api/orders/%s", t.region, url.QueryEscape(t.trackingId)), nil)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error creating getTrackingData request: "+err.Error())
		return errors.New("Error creating getTrackingData request: " + err.Error())
	}

	req.Header = http.Header{
		"x-instana-t":        []string{"399d813590722063"},
		"sec-ch-ua":          []string{`"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`},
		"sec-ch-ua-mobile":   []string{"?0"},
		"x-instana-l":        []string{"1,correlationType=web;correlationId=399d813590722063"},
		"user-agent":         []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"content-type":       []string{"application/json"},
		"x-instana-s":        []string{"399d813590722063"},
		"sec-ch-ua-platform": []string{`"macOS"`},
		"accept":             []string{"*/*"},
		"sec-fetch-site":     []string{"same-origin"},
		"sec-fetch-mode":     []string{"cors"},
		"sec-fetch-dest":     []string{"empty"},
		"referer":            []string{"https://www.adidas.de/order-details?data=" + t.trackingId},
		"accept-encoding":    []string{"gzip, deflate, br"},
		"accept-language":    []string{"en-GB,en;q=0.9"},

		http.HeaderOrderKey: []string{
			"x-instana-t",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"x-instana-l",
			"user-agent",
			"content-type",
			"x-instana-s",
			"sec-ch-ua-platform",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error getting tracking data: "+err.Error())
		return errors.New("Error getting tracking data: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Adidas-Tracker", "Error getting tracking data: "+resp.Status)
		return errors.New("Error getting tracking data: " + resp.Status)
	}

	var trackingData TrackingResponse

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error reading tracking data response body: "+err.Error())
		return errors.New("Error reading tracking data response body: " + err.Error())
	}

	err = json.Unmarshal(body, &trackingData)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error unmarshalling tracking data response body: "+err.Error())
		return errors.New("Error unmarshalling tracking data response body: " + err.Error())
	}

	if len(trackingData.Shipments) == 0 {
		trackingData.Shipments = append(trackingData.Shipments, Shipment{})
	}

	t.Result = trackingData

	return nil
}

func (t *Tracker) getInvoiceList() error {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.adidas.%s/api/orders/invoice/list/%s", t.region, url.QueryEscape(t.Result.InvoiceListID)), nil)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error creating getInvoiceList request: "+err.Error())
	}

	req.Header = http.Header{
		"content-type": []string{"application/json"},
		//"x-timestamp":     []string{"1685909215348"},
		"x-timestamp":     []string{strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)},
		"accept":          []string{"*/*"},
		"x-instana-l":     []string{"1,correlationType=web;correlationId=7b55e312cb11596c"},
		"sec-fetch-site":  []string{"same-origin"},
		"accept-language": []string{"en-GB,en;q=0.9"},
		"accept-encoding": []string{"gzip, deflate, br"},
		"sec-fetch-mode":  []string{"cors"},
		"user-agent":      []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"x-instana-s":     []string{"7b55e312cb11596c"},
		"referer":         []string{"https://www.adidas.de/order-details?data=" + url.QueryEscape(t.trackingId)},
		"x-instana-t":     []string{"7b55e312cb11596c"},
		"sec-fetch-dest":  []string{"empty"},

		http.HeaderOrderKey: []string{
			"content-type",
			"x-timestamp",
			"accept",
			"x-instana-l",
			"sec-fetch-site",
			"accept-language",
			"accept-encoding",
			"sec-fetch-mode",
			"user-agent",
			"x-instana-s",
			"referer",
			"x-instana-t",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error getting invoice list: "+err.Error())
		return errors.New("Error getting invoice list: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			logger.Info("Adidas-Tracker", "No invoice list found, adidas is weird :/")
			return errors.New("Error getting invoice list:  adidas is weird :/")
		} else {
			logger.Error("Adidas-Tracker", "Error getting invoice list: "+resp.Status)
			return errors.New("Error getting invoice list: " + resp.Status)
		}
	}

	var invoiceList []InvoiceListResponse

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error reading invoice list response body: "+err.Error())
		return errors.New("Error reading invoice list response body: " + err.Error())
	}

	err = json.Unmarshal(body, &invoiceList)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error unmarshalling invoice list response body: "+err.Error())
		return errors.New("Error unmarshalling invoice list response body: " + err.Error())
	}

	t.InvoiceList = invoiceList

	return nil
}

func (t *Tracker) getInvoice(invoiceId string, i int) error {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.adidas.%s/api/orders/invoice/%s", t.region, url.QueryEscape(invoiceId)), nil)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error creating getInvoice request: "+err.Error())
	}

	req.Header = http.Header{
		"accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site":  []string{"none"},
		"accept-encoding": []string{"gzip, deflate, br"},
		"sec-fetch-mode":  []string{"navigate"},
		"user-agent":      []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"accept-language": []string{"en-GB,en;q=0.9"},
		"sec-fetch-dest":  []string{"document"},

		http.HeaderOrderKey: []string{
			"accept",
			"sec-fetch-site",
			"cookie",
			"accept-encoding",
			"sec-fetch-mode",
			"user-agent",
			"accept-language",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error getting invoice: "+err.Error())
		return errors.New("Error getting invoice: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Adidas-Tracker", "Error getting invoice: "+resp.Status)
		return errors.New("Error getting invoice: " + resp.Status)
	}

	dir := "invoices"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error("Adidas-Tracker", "Error creating directory: "+err.Error())
			return errors.New("Error creating directory: " + err.Error())
		}
	}

	file, err := os.Create(fmt.Sprintf("%s/%s-(%d).pdf", dir, t.OrderId, i))
	if err != nil {
		logger.Error("Adidas-Tracker", "Error creating file: "+err.Error())
		return errors.New("Error creating file: " + err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		logger.Error("Adidas-Tracker", "Error saving invoice: "+err.Error())
		return errors.New("Error saving invoice: " + err.Error())
	}

	return nil
}
