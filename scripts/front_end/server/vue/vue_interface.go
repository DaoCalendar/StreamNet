// Copyright by StreamNet team
// 拼装对noderank访问的参数，
// 该类是为了适配main.go 和 noderank.go 而开发的。无特殊逻辑
package vue

import (
	"fmt"
	"github.com/caryxiao/go-zlog"
	nr "github.com/triasteam/noderank"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Message ...
type Message struct {
	Code      int64       `json:"code"`
	Timestamp int64       `json:"timestamp"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

// DataTee ...
type DataTee struct {
	DataScore interface{} `json:"dataScore"`
	DataCtx   interface{} `json:"dataCtx"`
}

// AddNodeRequest ...
type AddNodeRequest struct {
	Attester string `json:"attester,omitempty"`
	Attestee string `json:"attestee,omitempty"`
	Score    int    `json:"score,omitempty"`
	Time     string `json:"time,omitempty"`
	Nonce    int    `json:"nonce,omitempty"`
	Address  string `json:"address,omitempty"`
	Host     string `json:"host,omitempty"`
	Sign     string `json:"sign,omitempty"`
	OriData  string `json:"oriData,omitempty"`
}

// QueryNodesRequest ...
type QueryNodesRequest struct {
	Period  int64  `json:"period"`
	NumRank int64  `json:"numRank"`
	Url     string `json:"url,omitempty"`
	Address string `json:"address,omitempty"`
	Sign    string `json:"sign,omitempty"`
	OriData string `json:"oriData,omitempty"`
}

// NodeDetailRequest ...
type NodeDetailRequest struct {
	RequestUrl    string `json:"requestUrl,omitempty"`
	RequestData   string `json:"requestData,omitempty"`
	RequestMethod string `json:"requestMethod,omitempty"`
}

// OCli ...
type OCli struct {
}

// AddAtInfo ...
type AddAtInfo interface {
	AddAttestationInfoFunction(_data []byte) Message
	GetRankFunction(_data []byte) Message
}

// AddAttestationInfoFunction 新增证实数据包装类，负责组装参数
func (o *OCli) AddAttestationInfoFunction(request *AddNodeRequest) Message {
	mess := Message{}
	newReq := new(AddNodeRequest)
	newReq.Attester = request.Attester
	newReq.Attestee = request.Attestee
	newReq.Score = request.Score
	newReq.Time = request.Time
	newReq.Nonce = request.Nonce
	//newReq.Address = request.Address
	//newReq.Sign = request.Sign
	info := make([]string, 5)
	info[0] = newReq.Attester
	info[1] = newReq.Attestee
	info[2] = strconv.Itoa(newReq.Score)
	//info[3] = newReq.Address
	info[3] = strconv.Itoa(newReq.Nonce)
	info[4] = newReq.Time
	//info[6] = newReq.Sign
	zlog.Logger.Info("vue info split content is ", info)
	err1 := nr.AddAttestationInfo("", request.Host, info)
	if err1 != nil {
		mess = Message{Code: 1, Timestamp: time.Now().Unix(), Message: "Failed to add node"}
		return mess
	}
	mess = Message{Code: 0, Timestamp: time.Now().Unix(), Message: "Node added successfully"}
	return mess
}

// GetRankFunction 排名查询包装类
func (o *OCli) GetRankFunction(request *QueryNodesRequest) Message {
	mess := Message{}
	teescore, teectx, err1 := nr.GetRank(request.Url, request.Period, request.NumRank)
	if teectx == nil || err1 != nil || teescore == nil {
		mess = Message{Code: 1, Timestamp: time.Now().Unix(), Message: "Failed to query node data"}
		return mess
	}
	data := DataTee{teescore, teectx}
	mess = Message{Code: 0, Timestamp: time.Now().Unix(), Message: "Query node data successfully", Data: data}
	return mess
}

// QueryNodeDetail ...
func (o *OCli) QueryNodeDetail(request *NodeDetailRequest) Message {
	if request.RequestUrl == "" {
		return Message{Code: 0, Message: "RequestUrl is empty"}
	}
	result, err := httpSend(request.RequestUrl, request.RequestData, request.RequestMethod)
	if err == nil {
		return Message{Code: 1, Message: "Success!", Data: result}
	} else {
		fmt.Println(err)
		return Message{Code: 0, Message: "Query node's details failed!"}
	}
}

func httpSend(url string, param string, method string) (string, error) {
	payload := strings.NewReader(param)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	body, _ := ioutil.ReadAll(res.Body)

	return string(body), nil
}
