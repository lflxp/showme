package pkg

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/tatsushid/go-fastping"
)

type PingResult struct {
	Ip  string
	Rtt string
}

// 80,3306,25-200
func ParsePorts(ports string) ([]int, error) {
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
func ParseIps(in string) ([]string, error) {
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

func Pings2(ips []string, stop chan string, w io.Writer) error {
	active := 0
	p := fastping.NewPinger()
	t1 := time.Now()
	for _, x := range ips {
		ra, err := net.ResolveIPAddr("ip4:icmp", x)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Fprintln(w, fmt.Sprintf("[%s] - scaning:  %s", time.Now().Format("2006-01-02 15:04:05"), ra.String()))
		p.AddIPAddr(ra)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		active++
		fmt.Fprintln(w, fmt.Sprintf("IP Addr: %s receive, RTT: %v", addr.String(), rtt))
		stop <- fmt.Sprintf("%s:%v", addr.String(), rtt)
	}
	p.OnIdle = func() {
		elapsed := time.Since(t1)
		fmt.Fprintln(w, Colorize(fmt.Sprintf("[%s] - Finished:  Count - %d | Online - %d | Elapsed - %s", time.Now().Format("2006-01-02 15:04:05"), len(ips), active, elapsed.String()), "red", "", false, true))
	}
	err := p.Run()
	if err != nil {
		// return err
		fmt.Fprintln(w, err.Error())
		stop <- err.Error()
	}
	defer p.Stop()
	// defer close(stop)
	return nil
}

func Pings3(ips []string, stop chan string) error {
	p := fastping.NewPinger()

	for _, x := range ips {
		ra, err := net.ResolveIPAddr("ip4:icmp", x)
		if err != nil {
			fmt.Println(err)
			return err
		}
		p.AddIPAddr(ra)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		for i := 0; i <= 3306; i++ {
			if ScanPort(addr.String(), string(i)) {
				stop <- fmt.Sprintf("IP Addr:port %s:%d  receive, RTT: %v\n", addr.String(), i, rtt)
			}
		}
	}
	// p.OnIdle = func() {
	// 	fmt.Println("finish")
	// }
	err := p.Run()
	if err != nil {
		// return err
		stop <- err.Error()
	}
	defer p.Stop()
	return nil
}

// 通过chan获取ip
func Pings(ips []string, rs chan PingResult) error {
	p := fastping.NewPinger()

	for _, x := range ips {
		ra, err := net.ResolveIPAddr("ip4:icmp", x)
		if err != nil {
			fmt.Println(err)
			return err
		}
		p.AddIPAddr(ra)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		tmp := PingResult{
			Ip:  addr.String(),
			Rtt: fmt.Sprintf("%v", rtt),
		}
		rs <- tmp
	}
	// p.OnIdle = func() {
	// 	fmt.Println("finish")
	// }
	err := p.Run()
	if err != nil {
		return err
	}
	defer p.Stop()
	return nil
}

func ScanPort(host, port string) bool {
	remote := fmt.Sprintf("%s:%s", host, port)
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", remote) //转换IP格式 // 根据域名查找ip
	//fmt.Printf("%s", tcpAddr)
	// conn, err := net.DialTCP("tcp", nil, tcpAddr) //查看是否连接成功
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), 500*time.Microsecond) //建立tcp连接
	if err != nil {
		// fmt.Printf("no==%s:%s\r\n", host, port)
		return false
	}
	defer conn.Close()
	// fmt.Printf("ok==%s:%s\r\n", host, port)
	return true
}

func ScanPort2H(ip string, ports string, stop chan string, w io.Writer) error {
	// ip = strings.Split(ip, "|")[1]
	// stop <- fmt.Sprintf("ip %s  port range %s", ip, ports)
	active := 0
	t1 := time.Now()
	pports, err := ParsePorts(ports)
	if err != nil {
		return err
	}
	for _, port := range pports {
		if ScanPort(ip, fmt.Sprintf("%d", port)) {
			active++
			fmt.Fprintln(w, fmt.Sprintf("[%s] - scaning: %s:%d Active", time.Now().Format("2006-01-02 15:04:05"), ip, port))
			stop <- fmt.Sprintf("%s:%d", ip, port)
		} else {
			fmt.Fprintln(w, fmt.Sprintf("[%s] - scaning: %s:%d Failed", time.Now().Format("2006-01-02 15:04:05"), ip, port))
		}
	}
	elapsed := time.Since(t1)
	// fmt.Fprintln(w, Colorize(fmt.Sprintf("[%s] - Count:  %d", time.Now().Format("2006-01-02 15:04:05"), len(pports)), "red", "", false, true))
	// fmt.Fprintln(w, Colorize(fmt.Sprintf("[%s] - Elapsed:  %s", time.Now().Format("2006-01-02 15:04:05"), elapsed.String()), "red", "", false, true))
	// fmt.Fprintln(w, Colorize(fmt.Sprintf("[%s] - Finished:  DONE", time.Now().Format("2006-01-02 15:04:05")), "red", "", false, true))
	fmt.Fprintln(w, Colorize(fmt.Sprintf("[%s] - Finished:  Count - %d | Online - %d | Elapsed - %s", time.Now().Format("2006-01-02 15:04:05"), len(pports), active, elapsed.String()), "red", "", false, true))

	return nil
}
