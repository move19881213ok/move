package win

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	//	"regexp"
	//	"runtime"
	"jsee/core"
	"strconv"
	"strings"
)

var (
	_                 = bytes.Compare
	_                 = os.Args
	_                 = io.Copy
	_                 = bufio.NewWriterSize
	PSLIST_THREAD_KEY = [...]string{"Tid", "Pri", "Cswtch", "State", "User Time", "Kernel Time", "Elapsed Time"}
	PSLIST_THREAD_KL  = len(PSLIST_THREAD_KEY)
)

const (
	WMIC    = "wmic"
	PSTOOLS = "win\\bin\\PSTools\\"
	PSLIST  = PSTOOLS + "pslist.exe"
)

/**
 *	获取进程信息（通过wmic）
 *	1、调用windows的wmic process [condition] list full
 *	2、进行数据拼装和排序
 */
func ProcessInfosList(continueStr, orderProperty string, order int) []map[string]string {
	//1、
	var params []string = make([]string, 0, 5)
	params = append(params, "process")
	if continueStr != "" {
		params = append(params, "where")
		//		params = append(params, "name=\""+continueStr+"\"")
		params = append(params, "name='"+continueStr+"'")
	}
	params = append(params, "list")
	params = append(params, "full")
	res := core.ExecCommand(WMIC, params...)

	//2、
	if res == nil {
		return nil
	}

	var list []map[string]string = make([]map[string]string, 0, 5)
	tmpMap := make(map[string]string)

	for _, v := range res {
		if len(strings.TrimSpace(v)) == 0 {
			if tmpMap != nil && len(tmpMap) > 0 && tmpMap["Name"] != "System Idle Process" {
				list = append(list, tmpMap)
				tmpMap = make(map[string]string)
			}
			continue
		}

		tokens := strings.Split(v, "=")
		if len(tokens) > 1 {
			tmpMap[tokens[0]] = tokens[1]
		}
	}

	if tmpMap != nil && len(tmpMap) > 0 && tmpMap["Name"] != "System Idle Process" {
		list = append(list, tmpMap)
	}

	for _, item := range list {
		var sumModeTime int64 = core.GetLong("KernelModeTime", item) + core.GetLong("UserModeTime", item)
		item["sumModeTime"] = strconv.FormatInt(sumModeTime, 10)
	}

	core.SelectSort(list, core.SortNumDesc)

	var tmp string
	for i, item := range list {
		if item["sumModeTime"] == "0" {
			tmp = "\t\t"
		} else {
			tmp = "\t"
		}
		fmt.Println(i+1, "\t", item["sumModeTime"], tmp, item["Name"], "\t", item["ExecutablePath"])
	}
	return list
}

/**
 *	获取指定线程信息（通过pslist）
 *	1、调用windows的pslist -dmx 6460
 *	2、进行数据拼装和排序
 */
func PslistGetThread(pid string, orderProperty string, order int) []map[string]string {
	//执行pslist -dmx pid
	var params []string = []string{"-dmx", pid}
	res := core.ExecCommand(PSLIST, params...)

	//2、
	if res == nil {
		return nil
	}

	var list []map[string]string = make([]map[string]string, 0, 5)
	var tmpMap map[string]string

	isStart := false
	for _, line := range res {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if !isStart && strings.HasPrefix(line, PSLIST_THREAD_KEY[0]) { //以线程的标题开头
			isStart = true
			continue
		}
		if isStart {
			values := strings.Fields(line)
			tmpMap = make(map[string]string)
			for i := 0; i < PSLIST_THREAD_KL; i++ {
				tmpMap[PSLIST_THREAD_KEY[i]] = values[i]
			}
			list = append(list, tmpMap)
		}
	}

	core.SelectSort(list, core.SortStringAsc)

	return list
}
