package main

import (
	"AdidasTracker/src/cli"
	"AdidasTracker/src/pkg/logger"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/osext"
	//"AdidasTracker/src/pkg/logger"
	//"io/ioutil"
	//"os"
	//"github.com/kardianos/osext"
)

type Order struct {
	OrderNumber string
	Email       string
	Invoice     string
}

type Tracked struct {
	OrderNumber string
	Email       string

	Status            string
	Name              string
	Size              string
	SKU               string
	EstimatedDelivery string

	TrackingURL    string
	TrackingNumber string
	Carrier        string
}

func main() {
	//getting the executable path
	executablePath, _ := osext.ExecutableFolder()

	//setting the os path to the executable path
	os.Chdir(executablePath)

	//read proxies.txt
	proxiesTxt, err := ioutil.ReadFile("proxies.txt")
	if err != nil {
		logger.Error("INIT", "Error reading proxies.txt")
		return
	}

	//split after each line
	proxies := strings.Split(string(proxiesTxt), "\n")
	if len(proxies) == 0 {
		proxies[0] = ""
	}

	//read orders.csv
	ordersCsv, err := ioutil.ReadFile("orders.csv")
	if err != nil {
		logger.Error("INIT", "Error reading orders.csv")
		return
	}

	reader := csv.NewReader(strings.NewReader(string(ordersCsv)))

	// Ignore headers
	_, err = reader.Read()
	if err != nil {
		logger.Error("INIT", "Error reading orders.csv")
	}

	orders := make([]Order, 0)

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		orders = append(orders, Order{
			OrderNumber: line[0],
			Email:       line[1],
			Invoice:     line[2],
		})
	}

	trackedOrders := make([]Tracked, 0)

	counter := 0

	for _, order := range orders {
		if counter >= len(proxies) {
			counter = 0
		}

		proxy := proxies[counter]

		tracker := cli.Tracker{
			OrderId: order.OrderNumber,
			Email:   order.Email,
			Proxy:   proxy,
			Invoice: order.Invoice,
		}

		tracker.Track()

		if tracker.Error != nil {
			logger.Error("TRACKER", fmt.Sprintf("Error tracking order %s: %s", order.OrderNumber, tracker.Error))
			continue
		}

		trackedOrders = append(trackedOrders, Tracked{
			OrderNumber:       tracker.Result.OrderNo,
			Email:             tracker.Result.CustomerEmail,
			Status:            tracker.Result.Status,
			Name:              tracker.Result.ProductLineItems[0].ProductName,
			Size:              tracker.Result.ProductLineItems[0].LiteralSize,
			SKU:               tracker.Result.ProductLineItems[0].ArticleNumber,
			EstimatedDelivery: fmt.Sprintf("%s - %s", tracker.Result.ProductLineItems[0].EstimatedDeliveryPeriod.From, tracker.Result.ProductLineItems[0].EstimatedDeliveryPeriod.To),
			TrackingURL:       tracker.Result.Shipments[0].TrackingURL,
			TrackingNumber:    tracker.Result.Shipments[0].TrackingNo,
			Carrier:           tracker.Result.Shipments[0].SCAC,
		})
	}
	csvFilePath := filepath.Join(executablePath, "trackedOrders.csv")

	file, err := os.Create(csvFilePath)
	if err != nil {
		logger.Error("CSV", fmt.Sprintf("Error creating trackedOrders.csv: %s", err))
	}

	if err != nil {
		logger.Error("CSV", fmt.Sprintf("Error creating trackedOrders.csv: %s", err))
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	err = writer.Write([]string{"OrderNumber", "Email", "Status", "Name", "Size", "SKU", "EstimatedDelivery", "TrackingURL", "TrackingNumber", "Carrier"})
	if err != nil {
		logger.Error("CSV", fmt.Sprintf("Error writing headers to trackedOrders.csv: %s", err))
	}

	// Write data
	for _, order := range trackedOrders {
		err := writer.Write([]string{order.OrderNumber, order.Email, order.Status, order.Name, order.Size, order.SKU, order.EstimatedDelivery, order.TrackingURL, order.TrackingNumber, order.Carrier})
		if err != nil {
			logger.Error("CSV", fmt.Sprintf("Error writing order %s to trackedOrders.csv: %s", order.OrderNumber, err))
		}
	}

	logger.Info("CREDIT", "Made with ❤️  by @prizzle")
}
