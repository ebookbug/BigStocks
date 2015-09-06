// code
package main

import (
	"encoding/csv"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/crufter/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func FetchAllStockCode(fname string) {
	var url = "http://quote.eastmoney.com/stocklist.html"
	htm, err := goquery.ParseUrl(url)
	if err != nil {
		fmt.Println("Error:", err)
	}

	quote := htm.Find("div.quotebody")

	links := quote.Find("a")

	lns := links.Attrs("href")

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer f.Close()
	w := csv.NewWriter(f)

	for i := range lns {
		herf := lns[i]
		if strings.HasSuffix(herf, ".html") && strings.Contains(herf, "") {
			fmt.Println(herf)
			ix1 := strings.LastIndex(herf, "/")
			ix2 := strings.Index(herf, ".html")
			s := &Stock{}
			s.Code = herf[ix1+1 : ix2][2:]
			s.Area = herf[ix1+1 : ix2][:2]
			getStockInfoFromSina(s)

			w.Write([]string{s.Code, s.Area, s.Name})

			//			f.WriteString(code+"\r\n")
			//			fmt.Println(herf[27:35])
		}
	}

	w.Flush()
	f.Close()
}

func ReadAllStockCode(fname string) (stks []Stock, err error) {
	f, err := os.OpenFile(fname, os.O_RDONLY, os.ModePerm)

	if err != nil {
		log.Fatal(err)
		return stks, err
	}

	defer f.Close()
	reader := csv.NewReader(f)
	record, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return stks, err
	}

	stks = make([]Stock, len(record))

	for i := 0; i < len(record); i++ {
		stks[i] = Stock{
			Name: record[i][2],
			Code: record[i][0],
			Area: record[i][1],
		}
	}

	return
}

func getStockInfoFromSina(s *Stock) {
	var qturl = "http://hq.sinajs.cn/list="
	resp, _ := http.Get(qturl + s.Area + s.Code)
	ctn, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	enc := mahonia.NewDecoder("gbk")
	gbct := enc.ConvertString(string(ctn))

	if strings.Contains(gbct, "none_match") {
		return
	}

	fmt.Println(s.Area, s.Code, ":", gbct)
	inx1 := strings.Index(gbct, "\"")
	inx2 := strings.LastIndex(gbct, "\"")

	fmt.Println(gbct[inx1+1 : inx2])

	ctns := strings.Split(gbct[inx1+1:inx2], ",")
	s.Name = ctns[0]

	fmt.Println(ctns)

}
