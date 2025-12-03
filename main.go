package main

import (
	"bufio"
	"database/compute"
	"database/storage"
	"database/storage/engine"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting KV database")

	eng := engine.NewInMemoryEngine(logger)
	stor := storage.NewStorage(eng, logger)
	comp := compute.NewCompute(stor, logger)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("KV Database v1.0")
	fmt.Println("Enter commands (SET key value, GET key, DEL key):")
	fmt.Print("> ")

	for scanner.Scan() {
		query := scanner.Text()
		query = strings.TrimSpace(query)

		if query == "" {
			fmt.Print("> ")
			continue
		}

		if strings.ToUpper(query) == "EXIT" || strings.ToUpper(query) == "QUIT" {
			logger.Info("Shutting down database")
			fmt.Println("Goodbye!")
			break
		}

		result, err := comp.Execute(query)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println(result)
		}

		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Scanner error", zap.Error(err))
	}
}
