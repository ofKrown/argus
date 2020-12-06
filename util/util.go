package util

import (
	"log"
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"bytes"
	"../configuration"
)

// LogError : logs error ¯\_(ツ)_/¯
func LogError(err error) {
	if err != nil {
		log.Fatalln(err);
	}
}

// // Query :
// func QueryFactory(urlString string) (query Query) {
// 	hostURL, err := url.Parse(urlString)
// 	LogError(err)

// 	query := hostURL.Query()
// 	query.Add("stuff", "value")
//     hostURL.RawQuery = query.Encode()
	
// 	queryString = hostURL.String()
	
// 	return 
// }

// DoHTTPRequest : 
func DoHTTPRequest(httpMethod string, urlString string, payload []byte) (responseBody []byte) {
	config := configuration.GetConfig()
	url, err := url.Parse(urlString)
	LogError(err)

	if (config.Debug == true) {
		fmt.Println("REQUEST")
		fmt.Println(fmt.Sprintf("%s %s",httpMethod, url))
	}

	request := &http.Request {
		Method: httpMethod,
		URL: url,
	}

	if (payload != nil) {
		LogError(err)

		if (config.Debug == true) {
			fmt.Println("REQUEST BODY")
			fmt.Println(string(payload))
		}

		payloadBuffer := bytes.NewBuffer(payload)
		request.Header = map[string][]string{
				"Content-Type": { "application/json; charset=UTF-8" },
				"Client": { "Argus" },
		}
		request.Body = ioutil.NopCloser(payloadBuffer)
	}

	response, err := http.DefaultClient.Do(request)

	LogError(err)
	defer response.Body.Close()

	LogError(err)
	defer response.Body.Close()
	responseBody, err = ioutil.ReadAll(response.Body)
	LogError(err)

	if  (config.Debug == true) {
		fmt.Println("RESPONSE")
		fmt.Println(response)
	}

	return
}

// DebugResponseData :
func DebugResponseData(responseDataString string) {
	config := configuration.GetConfig()
	if  (config.Debug == true) {
		fmt.Println("RESPONSE DATA")
		fmt.Println(responseDataString)
	}
}
