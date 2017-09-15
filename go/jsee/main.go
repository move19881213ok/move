package main

import (
	"bufio"
	"fmt"
	"io"
	"jsee/core"
	"jsee/win"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	_ = io.Copy
	_ = bufio.NewWriterSize
)

const (
	JSTACK = "jstack"
	JPS    = "jps"
)

/**
 *	查看javacpu及内存使用情况，用来运维分析
 */
func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("无参数")
		return
	}

	/*
		s, e := exec.LookPath("java")
		if e != nil {
			fmt.Println("e:", e)
		}
		fmt.Println("ping:", s)
	*/

	osStr := runtime.GOOS
	if strings.HasPrefix(osStr, "win") {
		switch args[1] {
		case "pslist":
			{
				if len(args) < 3 {
					fmt.Println("请输入pid")
					return
				}
				match, err := regexp.MatchString("^[0-9]*$", args[2])
				if err != nil {
					fmt.Println("请输入合法数字", err.Error())
					return
				}
				if !match {
					fmt.Println("请输入合法数字")
					return
				}

				list := win.PslistGetThread(args[2], "", 0)

				if len(list) > 0 {
					stackTraceMap := getStackTrace(args[2])
					for _, item := range list {
						//把线程id转换为16进制
						tid := item["Tid"]
						num, err := strconv.ParseInt(tid, 10, 32)
						if err != nil {
							panic("str:" + tid + ",err:" + err.Error())
						}
						tid2 := strconv.FormatInt(num, 16)
						fmt.Println("tid2:", tid2, ":", tid, "User Time:", item["User Time"])
						for _, line := range stackTraceMap["0x"+tid2] {
							fmt.Println(line)
						}
						fmt.Println()
					}
				} else {
					fmt.Println("为查找到进程" + args[2] + "的信息")
				}
			}
		case "jps":
			{

			}
		case "ps": //进程信息
			{
				var name, order string = "", ""
				desc := 0
				for i := 2; i < len(args); i++ {
					switch args[i] {
					case "-name":
						{
							i++
							if len(args) == i || strings.HasPrefix(args[i], "-") {
								fmt.Println("缺少参数")
								return
							}
							name = args[i]
						}
					case "-o":
						{
							i++
							if len(args) == i || strings.HasPrefix(args[i], "-") {
								fmt.Println("缺少参数")
								return
							}
							order = args[i]
						}
					case "-r":
						desc = 1
					}
				}
				win.ProcessInfosList(name, order, desc)
			}
		}
	} else if strings.HasPrefix(osStr, "linux") {
		fmt.Println("linux os,osStr:" + osStr)
	} else {
		fmt.Println("unknow os,osStr:" + osStr)
	}
}

/**
 *	使用jstack，通过pid和线程id，获取指定线程的调用栈
 */
func getStackTrace(pid string) map[string][]string {
	var params []string = []string{pid}
	res := core.ExecCommand(JSTACK, params...)

	//2、
	if res == nil {
		fmt.Println(JSTACK + "未能查询到信息,pid:" + pid)
		return nil
	}

	stackTraceMap := make(map[string][]string)
	var nid string
	stackTraceArr := make([]string, 0, 5)
	for _, line := range res[2:] {
		if strings.TrimSpace(line) == "" {
			if nid != "" {
				stackTraceMap[nid] = stackTraceArr
				nid = ""
				stackTraceArr = make([]string, 0, 5)
			}
			continue
		}

		if strings.HasPrefix(line, "JNI") {
			continue
		}

		if nid == "" {
			nidIndex := strings.Index(line, "nid=")
			nid = core.Substr2(line, nidIndex, len(line))
			nidIndex = strings.Index(nid, " ")
			nid = core.Substr2(nid, 4, nidIndex)
		}
		stackTraceArr = append(stackTraceArr, line)
	}
	if nid != "" {
		stackTraceMap[nid] = stackTraceArr
	}

	return stackTraceMap
}
