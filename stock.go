// stock
package main

import (
	"encoding/csv"
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Stock struct {
	Name       string
	Code       string
	Area       string
	HighPrice  float32
	LowPrice   float32
	OpenPrice  float32
	ClosePrice float32
	VolNum     int64
	Time       time.Time
}

func getStockInfo(code string) (s *Stock, err error) {
	var qturl = "http://qt.gtimg.cn/q=s_"
	resp, _ := http.Get(qturl + code)
	ctn, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	enc := mahonia.NewDecoder("gbk")
	gbct := enc.ConvertString(string(ctn))
	fmt.Println(gbct)
	if strings.Contains(gbct, "none_match") {
		return s, nil
	}

	inx1 := strings.Index(gbct, "\"")
	inx2 := strings.LastIndex(gbct, "\"")

	fmt.Println(gbct[inx1+1 : inx2])

	ctns := strings.Split(gbct[inx1+1:inx2], "~")
	s = &Stock{}
	s.Name = ctns[1]
	s.Code = code[2:]
	s.Area = code[:2]

	return
}

func getStockHistoryPrice(code string, stype string) (stks []Stock, err error) {
	s, _ := getStockInfo(code)

	var ycode string
	if strings.Contains(code, "sh") {
		ycode = code[2:] + ".SS"
	} else if strings.Contains(code, "sz") {
		ycode = code[2:] + ".SZ"
	} else {
		ycode = code
	}

	var url = "http://table.finance.yahoo.com/table.csv?s="
	resp, err := http.Get(url + ycode + "&g=" + stype)
	if err != nil {
		fmt.Println("Get stock failed")
		return
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)

	record, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return stks, err
	}

	stks = make([]Stock, len(record)-1)
	for i := 1; i < len(record); i++ {
		stks[i-1] = Stock{
			Name:       s.Name,
			Code:       s.Code,
			Area:       s.Area,
			OpenPrice:  Atofloat((record[i][1]))*1.1 + 1.1,
			HighPrice:  Atofloat(record[i][2])*1.1 + 1.1,
			LowPrice:   Atofloat(record[i][3])*1.1 + 1.1,
			ClosePrice: Atofloat(record[i][4])*1.1 + 1.1,
			Time:       Atotime(record[i][0]),
		}
	}
	return
}

func Atotime(str string) time.Time {
	t, _ := time.Parse("yyyy-MM-dd", str)
	return t
}

func Atofloat(str string) float32 {
	val, _ := strconv.ParseFloat(strings.TrimSpace(str), 32)
	return float32(val) + 0.00001
}
func Atoi64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}
