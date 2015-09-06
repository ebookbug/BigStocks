// BigStocks project main.go
package main

import (
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"path/filepath"
)

var CfgFile *string = flag.String("config", "./conf/config.toml", "the stock code file")

func main() {
	flag.Parse()
	absCfg, _ := filepath.Abs(*CfgFile)
	fmt.Println("Using config file:", absCfg)
	cfg, err := toml.LoadFile(absCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	CodeFile := cfg.Get("stock.path").(string)
	dbHost := cfg.Get("database.server").(string)
	dbPort := cfg.Get("database.port").(int64)

	if flag.Arg(0) == "fetch" {
		fmt.Println("fetch stock code")
		FetchAllStockCode(CodeFile)
	} else {
		fmt.Println("Using the stock code file:", CodeFile)
		stks, _ := ReadAllStockCode(CodeFile)
		fmt.Println(stks[0].Code)
	}

	//	getStockInfo("sh600000")

	//	getStockHistoryPrice("sh600000","m")

	conn, err := Connection(dbHost, dbPort)
	if err != nil {
		fmt.Println("Connection Error:", err)
	}

	result, err := QueryDB(conn, "select * from stock_daily")
	if err != nil {
		fmt.Println("Query Error:", err)
	}

	fmt.Println("Result length:", len(result))
	//	for i := 0; i < len(result); i++ {
	//		fmt.Println(result[i])
	//	}

	stks, err := getStockHistoryPrice("sh600000", "d")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("StockSize:", len(stks))
	err = WriteStockBatch(conn, stks, 10)
	if err != nil {
		log.Fatal("Write Failed:", err)
	}

	record, err := QueryDB(conn, "select * from stock_daily")
	if err != nil {
		log.Fatal("Query Faild:", err)
	}

	fmt.Println("Count", len(record[0].Series[0].Values))

}
