package goutils

import (
	"fmt"
	"net/url"
)

func ExamplePost() {
	uri := ""
	data := url.Values{}
	data.Add("user_pwd", "asdf1234")

	resp, _ := HTTPPost(uri, &data)

	fmt.Println(resp["message"])
	fmt.Println(resp["result"])

	// Output:
	// OK
	// 000
}

func ExampleGet() {
	uri := ""
	data := url.Values{}
	data.Add("ts", "1451282836296")
	data.Add("ticket", "0eeeb553c5861a699aa77d4c37cda513")
	data.Add("ip", "127.0.0.1")

	resp, _ := HTTPGet(uri, &data)

	fmt.Println(resp["message"])
	fmt.Println(resp["result"])

	// Output:
	// OK
	// 000
}

func ExampleGetToStruct() {
	uri := ""
	data := url.Values{}
	data.Add("ts", "1451282836296")
	data.Add("ticket", "0eeeb553c5861a699aa77d4c37cda513")
	data.Add("ip", "127.0.0.1")

	type session struct {
		User     string `json:"user"`
		Server   string `json:"server"`
		Resource string `json:"resource"`
	}

	type sessionList struct {
		Message string `json:"message"`
		Result  string `json:"result"`

		Value []session `json:"value"`
	}

	obj := sessionList{}

	err := HTTPGetToStruct(uri, &data, &obj)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(obj.Message)
	fmt.Println(obj.Result)

	// Output:
	// OK
	// 000
}
