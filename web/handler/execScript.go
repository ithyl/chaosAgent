package handler

import (
	"chaosAgent/pkg/bash"
	"chaosAgent/transport"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
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

// generateScriptFile creates a script file based on the execution ID and command.
func generateScriptFile(caseId, execID, cmd string) (string, error) {
	// Create a directory with the execution ID
	dirPath := filepath.Join("/path/to/scripts", execID)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Create the script file with the command
	scriptPath := filepath.Join(dirPath, "script.sh")
	file, err := os.Create(scriptPath)
	if err != nil {
		return "", fmt.Errorf("failed to create script file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to write to script file: %v", err)
	}

	return scriptPath, nil
}

func (ch *ExecScriptHandler) createCase(caseId, execId, cmd string) *transport.Response {
	start := time.Now()
	fields := strings.Fields(cmd)

	if len(fields) == 0 {
		logrus.Warningf("less command parameters")
		return transport.ReturnFail(transport.ParameterLess, "command")
	}

	scriptPath, err := generateScriptFile(caseId, execId, cmd)
	if err != nil {
		logrus.Errorf("failed to generate script file: %v", err)
		return transport.ReturnFail(transport.ScriptCreateFailed, err.Error())
	}

	// Execute the script
	resultCmd, errMsg, ok := bash.AsyncExecScript(context.Background(), scriptPath, cmd)
	diffTime := time.Since(start)
	logrus.Infof("execute chaosblade result, result: %s, errMsg: %s, ok: %t, duration time: %v, cmd : %v", "已执行", errMsg, ok, diffTime, cmd)
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
func NewExecScriptHandler(transportClient *transport.TransportClient) *ExecScriptHandler {
	return &ExecScriptHandler{
		running:         make(map[string]string, 0),
		mutex:           sync.Mutex{},
		transportClient: transportClient,
	}
}

func (ch *ExecScriptHandler) Handle(request *transport.Request) *transport.Response {
	logrus.Infof("chaosblade: %+v", request)
	caseId := request.Params["caseId"]
	execId := request.Params["execId"]
	cmd := request.Params["cmd"]

	if cmd == "" {
		return transport.ReturnFail(transport.ParameterEmpty, "cmd")
	}
	return ch.createCase(caseId, execId, cmd)

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
	//time.Sleep(10 * time.Second)
	//resultCmd.Process.Kill()
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
