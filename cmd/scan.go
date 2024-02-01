/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	scanIp    string
	scanPort  string
	scanDebug bool
)

type ScanData struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "多线程IP+端口扫描工具",
	Long:  `showme scan -i 127.0.0.1 -p 1-65535`,
	Run: func(cmd *cobra.Command, args []string) {
		if scanDebug {
			lvl.Set(slog.LevelDebug)
		} else {
			lvl.Set(slog.LevelInfo)
		}

		start := time.Now()

		data, _, err := ScanPort(scanIp, scanPort)
		if err != nil {
			slog.Error("端口扫描失败", "ERROR", err.Error())
			return
		}
		result := []ScanData{}
		for ips, value := range data {
			for status, ports := range value {
				if status == "open" {
					for _, port := range ports {
						tmp := ScanData{
							Ip:     ips,
							Status: status,
							Port:   port,
						}
						result = append(result, tmp)
					}
				}
			}
		}
		cost := time.Since(start)
		for index, x := range result {
			slog.Info("端口扫描结果", "序号", index, "IP", x.Ip, "PORT", x.Port, "状态", x.Status)
		}
		slog.Info("端口扫描完成", "耗时", cost.String(), "有效结果", len(result))
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scanCmd.Flags().BoolVarP(&scanDebug, "debug", "d", false, "Debug Log Info")
	scanCmd.Flags().StringVarP(&scanIp, "ip", "i", "127.0.0.1", "ip范围,例如:127.0.0.1")
	scanCmd.Flags().StringVarP(&scanPort, "port", "p", "1-65535", "端口范围，例如: 22,80,8080-9000")
}

func ScanPort(ip string, port string) (map[string]map[string][]int, []string, error) {
	ips := make(chan string, runtime.NumCPU())
	result := make(chan map[string]map[string][]int)
	resultFinally := make(map[string]map[string][]int)

	// 新建一个上下文
	ctx := context.Background()
	// 在初始上下文的基础上创建一个有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)
	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ipList, err := parseIps(ip)
	if err != nil {
		return nil, nil, err
	}
	slog.Debug("解析ip完成", "数量", len(ipList))

	// 生产者：生成ip主机扫描任务
	go func(ipList []string) {
		for _, i := range ipList {
			slog.Debug("放入扫描队列", "IP", i)
			ips <- i
		}
	}(ipList)

	// 消费者： 消费ip扫描任务
	for i := 0; i < cap(ips); i++ {
		slog.Debug("启动协程", "协程ID", i)
		go handleMultiIp(ctx, ips, port, result)
	}

	rsLength := len(ipList)
	for rsLength > 0 {
		resultIp := <-result
		for k, v := range resultIp {
			resultFinally[k] = v
		}
		rsLength--
	}

	// sort.Strings(ipList)

	// for _, i := range ipList {
	// 	for key, list := range resultFinally[i] {
	// 		if key == "open" {
	// 			for _, p := range list {
	// 				log.Infof("ip %s %s", i, key, p)
	// 			}
	// 		}
	// 	}
	// }
	return resultFinally, ipList, nil
}

func handleMultiIp(ctx context.Context, ips chan string, port string, result chan map[string]map[string][]int) {
	for {
		select {
		case i := <-ips:
			start := time.Now()
			open, close, err := tcpScanByGoroutineWithChannelAndSort(ctx, i, port)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			tmp := map[string][]int{
				"open":  open,
				"close": close,
			}
			cost := time.Since(start)
			slog.Debug("开始扫描", "IP", i, "PORT", port, "OPEN", open, "耗时", cost.String())

			result <- map[string]map[string][]int{i: tmp}
		case <-ctx.Done():
			slog.Debug("退出协程", "PORT", port)
			return
		}
	}
}

func ScanSingleIp(ip, port string) ([]int, []int, error) {
	return tcpScanByGoroutineWithChannelAndSort(context.Background(), ip, port)
}

// 80,3306,25-200
func parsePorts(ports string) ([]int, error) {
	rs := []int{}
	tmp := strings.Split(ports, ",")
	for _, x := range tmp {
		if strings.Contains(x, "-") {
			t1 := strings.Split(x, "-")
			start, err := strconv.Atoi(t1[0])
			if err != nil {
				return rs, err
			}
			end, err := strconv.Atoi(t1[1])
			if err != nil {
				return rs, err
			}

			if start > end {
				return rs, errors.New(fmt.Sprintf("%s start biger than end", x))
			}

			for i := start; i <= end; i++ {
				rs = append(rs, i)
			}
		} else {
			n, err := strconv.Atoi(x)
			if err != nil {
				return rs, err
			}
			rs = append(rs, n)
		}
	}
	return rs, nil
}

// 10.1.1.1
// 10-20.1.1.0
// 1.2.3.4-200
func parseIps(in string) ([]string, error) {
	rs := []string{}
	if strings.Contains(in, "-") {
		tmp_a := strings.Split(in, ".")
		if len(tmp_a) != 4 {
			fmt.Println(tmp_a)
			return nil, errors.New("ip地址不正确")
		}
		A := []string{}
		B := []string{}
		C := []string{}
		D := []string{}
		for m, n := range tmp_a {
			if strings.Contains(n, "-") {
				tmp := strings.Split(n, "-")
				a, err := strconv.Atoi(tmp[0])
				if err != nil {
					return rs, err
				}
				b, err := strconv.Atoi(tmp[1])
				if err != nil {
					return rs, err
				}
				for i := a; i <= b; i++ {
					if m == 0 {
						A = append(A, fmt.Sprintf("%d", i))
					} else if m == 1 {
						B = append(B, fmt.Sprintf("%d", i))
					} else if m == 2 {
						C = append(C, fmt.Sprintf("%d", i))
					} else if m == 3 {
						D = append(D, fmt.Sprintf("%d", i))
					}
				}
			} else {
				if m == 0 {
					A = append(A, n)
				} else if m == 1 {
					B = append(B, n)
				} else if m == 2 {
					C = append(C, n)
				} else if m == 3 {
					D = append(D, n)
				}
			}
		}
		for _, a1 := range A {
			for _, b1 := range B {
				for _, c1 := range C {
					for _, d1 := range D {
						rs = append(rs, fmt.Sprintf("%s.%s.%s.%s", a1, b1, c1, d1))
					}
				}
			}
		}
	} else {
		rs = append(rs, in)
	}
	return rs, nil
}

// The function handles checking if ports are open or closed for a given IP address.
func handleWorker(ctx context.Context, ip string, ports chan int, results chan int) {
	for p := range ports {
		select {
		default:
			address := fmt.Sprintf("%s:%d", ip, p)
			conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
			if err != nil {
				// fmt.Printf("[debug] ip %s Close \n", address)
				results <- (-p)
				continue
			}
			// fmt.Printf("[debug] ip %s Open \n", address)
			conn.Close()
			results <- p
		case <-ctx.Done():
			results <- (-1)
			return
		}
	}
}

func tcpScanByGoroutineWithChannelAndSort(ctx context.Context, ip string, port string) ([]int, []int, error) {
	start := time.Now()

	portList, err := parsePorts(port)
	if err != nil {
		return nil, nil, err
	}

	// 参数校验
	if err := verifyParam(ip, portList); err != nil {
		slog.Error("[Exit]\n")
		return nil, nil, err
	}

	ports := make(chan int, runtime.NumCPU())
	results := make(chan int)
	var openSlice []int
	var closeSlice []int

	// 任务生产者-分发任务 (新起一个 goroutinue ，进行分发数据)
	go func(data []int) {
		for _, i := range data {
			select {
			case <-ctx.Done():
				return
			default:
				ports <- i
			}
		}
	}(portList)

	// 任务消费者-处理任务  (每一个端口号都分配一个 goroutinue ，进行扫描)
	// 结果生产者-每次得到结果 再写入 结果 chan 中
	for i := 0; i < cap(ports); i++ {
		go handleWorker(ctx, ip, ports, results)
	}

	// 结果消费者-等待收集结果 (main中的 goroutinue 不断从 chan 中阻塞式读取数据)
	rsLength := len(portList)
	for rsLength > 0 {
		resPort := <-results
		if resPort > 0 {
			openSlice = append(openSlice, resPort)
		} else {
			closeSlice = append(closeSlice, -resPort)
		}
		rsLength--
	}

	// 关闭 chan
	close(ports)
	close(results)

	// 排序
	sort.Ints(openSlice)
	sort.Ints(closeSlice)

	// 输出
	// for _, p := range openSlice {
	// 	slog.Info(fmt.Sprintf("[info] %s:%-8d Open", ip, p))
	// }
	// for _, p := range closeSlice {
	//     fmt.Printf("[info] %s:%-8d Close\n", ip, p)
	// }

	cost := time.Since(start)
	slog.Debug(fmt.Sprintf("[tcpScanByGoroutineWithChannelAndSort] cost %s second \n", cost))
	return openSlice, closeSlice, nil
}

func verifyParam(ip string, ports []int) error {
	netip := net.ParseIP(ip)
	if netip == nil {
		slog.Error("[Error] ip type is must net.ip")
		return fmt.Errorf("[Error] ip type is must net.ip")
	}
	// log.Infof("[Info] ip=%s | ip type is: %T | ip size is: %d \n", netip, netip, unsafe.Sizeof(netip))

	for _, port := range ports {
		if port < 1 || port > 65535 {
			slog.Error("[Error] port is must in the range of 1~65535")
			return fmt.Errorf("[Error] port is must in the range of 1~65535")
		}
	}
	// log.Infof("[Info] port start:%d - end:%d \n", ports[0], ports[len(ports)-1])

	return nil
}
