package handler

import (
	"chaosAgent/pkg/bash"
	"chaosAgent/transport"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var execMap = make(map[string]*exec.Cmd)

type ExecScriptHandler struct {
	mutex   sync.Mutex
	running map[string]string

	transportClient *transport.TransportClient
}

func NewExecScriptHandler(transportClient *transport.TransportClient) *ExecScriptHandler {
	return &ExecScriptHandler{
		running:         make(map[string]string, 0),
		mutex:           sync.Mutex{},
		transportClient: transportClient,
	}
}

func (ch *ExecScriptHandler) Handle(request *transport.Request) *transport.Response {
	logrus.Infof("chaosblade: %+v", request)

	//todo 版本不一致时，需要update,这里是判断是否升级完成
	//if handler.blade.upgrade.NeedWait() {
	//	return transport.ReturnFail(transport.Code[transport.Upgrading], "agent is in upgrading")
	//}
	cmd := request.Params["cmd"]
	if cmd == "" {
		return transport.ReturnFail(transport.ParameterEmpty, "cmd")
	}
	return ch.exec(cmd)
}

func (ch *ExecScriptHandler) exec(cmd string) *transport.Response {
	start := time.Now()
	fields := strings.Fields(cmd)

	if len(fields) == 0 {
		logrus.Warningf("less command parameters")
		return transport.ReturnFail(transport.ParameterLess, "command")
	}
	//// 判断 chaosblade 是否存在
	//if !tools.IsExist(options.BladeBinPath) {
	//	logrus.Warningf(transport.Errors[transport.ChaosbladeFileNotFound])
	//	return transport.ReturnFail(transport.ChaosbladeFileNotFound)
	//}
	//command := fields[0]

	// 执行 命令
	result, errMsg, ok := bash.ExecScript(context.Background(), "/home/project/monitor/chaos/chaosAgent/test.sh", cmd)
	diffTime := time.Since(start)
	logrus.Infof("execute chaosblade result, result: %s, errMsg: %s, ok: %t, duration time: %v, cmd : %v", result, errMsg, ok, diffTime, cmd)
	if ok {
		// 解析返回结果
		response := parseResult(result)
		if !response.Success {
			logrus.Warningf("execute chaos failed, result: %s", result)
			return response
		}
		// 安全点处理
		//ch.handleCacheAndSafePoint(cmd, command, fields[1], response)
		return response
	} else {
		var response transport.Response
		err := json.Unmarshal([]byte(result), &response)
		if err != nil {
			logrus.Warningf("Unmarshal chaosblade error message err: %s, result: %s", err.Error(), result)
			return transport.ReturnFail(transport.ResultUnmarshalFailed, result, errMsg)
		} else {
			return &response
		}
	}
}
func (ch *ExecScriptHandler) asyncExec(cmd string) *transport.Response {
	//start := time.Now()
	fields := strings.Fields(cmd)

	if len(fields) == 0 {
		logrus.Warningf("less command parameters")
		return transport.ReturnFail(transport.ParameterLess, "command")
	}
	//// 判断 chaosblade 是否存在
	//if !tools.IsExist(options.BladeBinPath) {
	//	logrus.Warningf(transport.Errors[transport.ChaosbladeFileNotFound])
	//	return transport.ReturnFail(transport.ChaosbladeFileNotFound)
	//}
	//command := fields[0]

	// 执行 命令
	resultCmd, errMsg, ok := bash.AsyncExecScript(context.Background(), "/home/project/monitor/chaos/chaosAgent/test.sh", cmd)
	//diffTime := time.Since(start)
	//logrus.Infof("execute chaosblade result, result: %s, errMsg: %s, ok: %t, duration time: %v, cmd : %v", result, errMsg, ok, diffTime, cmd)
	time.Sleep(10 * time.Second)
	resultCmd.Process.Kill()
	if ok {
		execMap["aaa"] = resultCmd
		// 解析返回结果
		response := parseResult("执行中")
		//if !response.Success {
		//	logrus.Warningf("execute chaos failed, result: %s", result)
		//	return response
		//}
		// 安全点处理
		//ch.handleCacheAndSafePoint(cmd, command, fields[1], response)
		return response
	} else {
		var response transport.Response
		err := json.Unmarshal([]byte("执行中"), &response)
		if err != nil {
			logrus.Warningf("Unmarshal chaosblade error message err: %s, result: %s", err.Error(), "执行中")
			return transport.ReturnFail(transport.ResultUnmarshalFailed, "执行中", errMsg)
		} else {
			return &response
		}
	}
}
