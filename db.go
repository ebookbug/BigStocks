// db
package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
)

const (
	DBName = "stocks"
)

func Connection(dbHost string, dbPort int64) (conn *client.Client, err error) {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", dbHost, dbPort))
	if err != nil {
		log.Fatal(err)
		return conn, err
	}

	conf := &client.Config{URL: *u}

	conn, err = client.NewClient(*conf)

	if err != nil {
		log.Fatal(err)
		return conn, err
	}

	dur, ver, err := conn.Ping()

	if err != nil {
		return conn, err
	} else {
		log.Printf("Success make a influxdb connection,! %v, %s .\n", dur, ver)
	}
	return conn, nil
}

func QueryDB(conn *client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: DBName,
	}
	response, err := conn.Query(q)

	if err != nil {
		if response.Error() != nil {
			return res, response.Error()
		} else {
			return res, err
		}
	}
	res = response.Results
	return
}

func WriteStockBatch(conn *client.Client, stks []Stock, batchSize int) error {
	size := len(stks) / batchSize
	var nstks []Stock
	for i := 0; i < size; i++ {
		if (i+1)*batchSize > len(stks) {
			nstks = stks[i*batchSize : len(stks)]
		} else {
			nstks = stks[i*batchSize : (i+1)*batchSize]
		}
		WriteStock(conn, nstks)
	}
	return nil
}

func WriteStock(conn *client.Client, stks []Stock) error {
	pts := make([]client.Point, len(stks))
	for i := 0; i < len(stks); i++ {
		pts[i] = client.Point{
			Measurement: "stock_daily",
			Tags: map[string]string{
				"Name": stks[i].Name,
				"Code": stks[i].Code,
				"Area": stks[i].Area,
			},
			Fields: map[string]interface{}{
				"HighPrice":  stks[i].HighPrice,
				"LowPrice":   stks[i].LowPrice,
				"OpenPrice":  stks[i].OpenPrice,
				"ClosePrice": stks[i].ClosePrice,
			},
			Time: stks[i].Time,
		}
	}

	bth := client.BatchPoints{
		Points:          pts,
		Database:        DBName,
		RetentionPolicy: "default",
	}

	_, err := conn.Write(bth)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
