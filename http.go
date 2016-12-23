package goutils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// HTTPPost 함수는 x-www-form-urlencoded 파라미터를 입력으로 하며,
// application/json 타입의 응답을 받아,
// map[string]interface{} 형식의 결과로 반환한다.
func HTTPPost(uri string, data *url.Values) (map[string]interface{}, error) {

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// HTTPGet 함수는,
// application/json 타입의 응답을 받아,
// map[string]interface{} 형식의 결과로 반환한다.
func HTTPGet(uri string, data *url.Values) (map[string]interface{}, error) {

	var m map[string]interface{}

	err := HTTPGetToStruct(uri, data, &m)

	return m, err
}

// HTTPGetToStruct 함수는,
// application/json 타입의 응답을 받아,
// 특정 형식의 구조체로 반환한다.
func HTTPGetToStruct(uri string, data *url.Values, obj interface{}) error {

	req, err := http.NewRequest("GET", uri+"?"+data.Encode(), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&obj)
	if err != nil {
		return err
	}

	return nil
}
