package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	//	"regexp"
	//	"runtime"
	"strconv"
	"strings"
)

var (
	_ = io.Copy
	_ = bufio.NewWriterSize
	_ = os.Args
)

func SortStringAsc(item1, item2 map[string]string) bool {
	return strings.Compare(item1["User Time"], item2["User Time"]) < 0
}

func SortNumAsc(item1, item2 map[string]string) bool {
	return GetLong("sumModeTime", item1) > GetLong("sumModeTime", item2)
}
func SortNumDesc(item1, item2 map[string]string) bool {
	return GetLong("sumModeTime", item1) < GetLong("sumModeTime", item2)
}

/**
 *	从map中获取指定的int64
 */
func GetLong(key string, item map[string]string) int64 {
	if v, ok := item[key]; ok && v != "" {
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return num
		} else {
			panic(err)
		}
	} else if ok {
		return 0
	} else {
		panic("not key:" + key)
	}
}

/**
 *	执行控制台，并返回信息
 */
func ExecCommand(exeStr string, params ...string) []string {
	var cmd *exec.Cmd
	fmt.Println("exeStr:", exeStr, "params:", params)
	if len(params) > 0 {
		cmd = exec.Command(exeStr, params...)
	} else {
		cmd = exec.Command(exeStr)
	}

	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stdout

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
		return nil
	}

	var res []string = make([]string, 0, 10)

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if line != "" {
			res = append(res, strings.Replace(strings.Replace(line, "\n", "", -1), "\r", "", -1))
			//			res = append(res, strings.TrimSuffix())
		}
	}

	return res
}

//通用------------------------
//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//截取字符串 start 起点下标 end 终点下标(不包括)
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong,start:" + strconv.Itoa(start) + ",length:" + strconv.Itoa(length))
	}

	if end < 0 || end > length {
		panic("end is wrong,end:" + strconv.Itoa(end) + ",length:" + strconv.Itoa(length))
	}

	return string(rs[start:end])
}

/**
 *	选择排序
 */
func SelectSort(list []map[string]string, compFunc func(map[string]string, map[string]string) bool) {
	length := len(list)
	for i := 0; i < length; i++ {
		index := 0
		//保存索引值
		for j := 1; j < length-i; j++ {
			if compFunc(list[j], list[index]) {
				index = j
			}
		}
		list[length-i-1], list[index] = list[index], list[length-i-1]
	}
}
