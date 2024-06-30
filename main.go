package main

import (
	"bufio"
	"design_patterns/stock_broker/models"
	"design_patterns/stock_broker/supports"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// Create a channel to communicate between the CLI and the background task
	stopChan := make(chan bool)
	tradeService := supports.NewTradeService()

	orderProcessor := supports.NewOrderProcessor(tradeService)
	orderService := supports.NewOrderService(orderProcessor)

	go tradeService.Show()

	go orderProcessor.GetMatches()

	// Start the CLI loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		text, _ := reader.ReadString('\n')
		texts := strings.Split(text, " ")

		// Handle different commands
		switch texts[0] {
		case "exit":
			fmt.Println("Exiting...")
			stopChan <- true // Signal the background task to stop
			return
		case "place_order":
			var orderType models.OrderType
			switch texts[1] {
			case "buy":
				orderType = models.Buy
			case "sell":
				orderType = models.Sell
			default:
				fmt.Println("unsupported order type")
				continue
			}
			price, _ := strconv.Atoi(texts[3])
			quantity, _ := strconv.Atoi(texts[4])
			orderService.AddOrder("1", orderType, texts[2], quantity, float64(price))
		default:
			fmt.Println("Unknown command:", text)
		}
		spew.Printf("buyHeap: %+v\n", orderProcessor.Buy)
		spew.Printf("sellHeap: %v\n", orderProcessor.Sell)
	}
}
