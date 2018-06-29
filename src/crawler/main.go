package main

/*
//包管理工具
go get -v github.com/gpmgo/gopm
gopm get -g -v golang.org/x/tools/cmd/goimports

go bulid golang.org/x/tools/cmd/goimports
go install  golang.org/x/tools/cmd/goimports

cd cmd
go install ./...
*/

// gopm get -g -v golang.org/x/text
// gopm get -g -v golang.org/x/net

import (
	"bufio"
	"crawler/engine"
	"crawler/zhenai/parser"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func main() {
	engine.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}

func ver1() {
	resp, err := http.Get(
		"http://www.zhenai.com/zhenghun")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}
	//utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	e := determineEncodeing(resp.Body)
	utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())
	all, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", all)
	printCityList(all)

}

func determineEncodeing(r io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

func printCityList2(contents []byte) {
	re := regexp.MustCompile(`<a href="http://www.zhenai.com/zhenghun/[\w]+"[^>]*>[^<]+</a>`)
	matches := re.FindAll(contents, -1)
	for _, m := range matches {
		fmt.Printf("%s\n", m)
	}
	fmt.Printf("Matches found:%d\n", len(matches))
}

func printCityList(contents []byte) {
	re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[\w]+)"[^>]*>([^<]+)</a>`)
	matches := re.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		fmt.Printf("City:%s,URL:%s\n", m[2], m[1])
	}
	fmt.Printf("Matches found:%d\n", len(matches))
}
