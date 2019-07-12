package collector

/*
硬件监控
采集物理机静态属性
*/
import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/lflxp/showme/utils"
	log "github.com/sirupsen/logrus"
)

var HardwareStatic *Hardware
var stopChan chan bool

func init() {
	ticker := time.NewTicker(3 * time.Second)
	stopChan = make(chan bool)
	HardwareStatic = NewHardware()
	err := HardwareStatic.GetAll()
	if err != nil {
		log.Error(err.Error())
	}
	go func(ticker *time.Ticker) {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := HardwareStatic.GetAll()
				if err != nil {
					log.Error(err.Error())
				}
			case stop := <-stopChan:
				if stop {
					log.Debug("Ticker Stop")
					return
				}
			}
		}
	}(ticker)
}

// machine 物理机信息
type Machine struct {
	Manufacturer  string `json:"manufacturer"`
	Type          string `json:"type"` // Product Name: PowerEdge R730
	Sn            string `json:"sn"`
	Os            string `json:"os"`            // 系统版本 CentOs Linux
	OsCoreVersion string `json:"osCoreVersion"` // 内核版本
	OsPlatm       string `json:"osPlatm"`       // 系统架构 x86_64
	OsVersion     string `json:"osVersion"`     // 系统爸爸号 7.2.1511
}

// # 总核数 = 物理CPU个数 X 每颗物理CPU的核数
// # 总逻辑CPU数 = 物理CPU个数 X 每颗物理CPU的核数 X 超线程数
// # 查看物理CPU个数
// cat /proc/cpuinfo| grep "physical id"| sort| uniq| wc -l
// # 查看每个物理CPU中core的个数(即核数)
// cat /proc/cpuinfo| grep "cpu cores"| uniq
// # 查看逻辑CPU的个数
// cat /proc/cpuinfo| grep "processor"| wc -l
type Cpu struct {
	Type        string `json:"type"`        // Central Processor
	CacheSize   string `json:"cacheSize"`   // 缓存
	PhysicalNum string `json:"physicalNum"` //  物理cpu个数
	CoreNum     string `json:"coreNum"`     // 核数
	Processor   string `json:"processor"`   // 逻辑cpu个数
}

type Memory struct {
	Capacity     string `json:"capacity"`     // 单卡容量
	Manufacturer string `json:"manufacturer"` // 厂商
	Type         string `json:"type"`         // 型号
	SerialNumber string `json:"serialNumber"` // 序列号
	Speed        string `json:"speed"`        // 2400 MHz
	PartNumber   string `json:"partNumber"`
}

type Disk struct {
	MountPath string `json:"mountPath"` // 挂载路径
	Capacity  string `json:"capacity"`  // 容量（默认为GB）
}

type Gpu struct {
	Type string `json:"type"`
}

// 网卡信息 主要取Mac地址
type NetInterface struct {
	Name    string `json:"name"` // 网卡名
	Ip      string `json:"ip"`
	Macaddr string `json:"macaddr"`
}

// 硬件信息
type Hardware struct {
	Hostname    string         `json:"hostname"`
	Net         []NetInterface `json:"net"`
	TotalMemory string         `json:"totalMemory"`
	TotalDisk   string         `json:"totalDisk"`
	TotalGpu    string         `json:"totalGpu"`
	Cpu         Cpu            `json:"cpu"`
	Machine     Machine        `json:"machine"`
	Memory      []Memory       `json:"memory"`
	Disk        []Disk         `json:"disk"`
	Gpu         []Gpu          `json:"gpu"`
	Status      []string       `json:"status"`     // 结果状态返回
	Prometheus  string         `json:"prometheus"` // 监控端口记录
	SampleLoad  LoadStatus     `json:"sampleLoad"` // 简单监控
}

type LoadStatus struct {
	Load       string `json:"load"`       // 5分钟内负载
	CpuUsed    string `json:"cpuUsed"`    // cpu使用率
	MemoryUsed string `json:"memoryUsed"` // 内存使用率
	DiskIoWait string `json:"diskIoWait"` // 磁盘iowait
}

func NewHardware() *Hardware {
	return &Hardware{}
}

func (this *Hardware) SetProm(add string) {
	this.Prometheus = add
}

func (this *Hardware) GetAll() error {
	this.Status = []string{}
	this.GetNet()
	err := this.GetHostname()
	if err != nil {
		log.Error(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	err = this.GetMachine()
	if err != nil {
		log.Error(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	log.Debug(1)
	err = this.GetCpu()
	if err != nil {
		log.Error(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	log.Debug(2)
	err = this.GetTotal()
	if err != nil {
		log.Error(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	log.Debug(3)
	err = this.GetMemory()
	if err != nil {
		log.Error(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	log.Debug(4)
	err = this.GetDisk()
	if err != nil {
		log.Error(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	log.Debug(5)
	err = this.GetGpu()
	if err != nil {
		log.Debug(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	log.Debug(6)
	err = this.GetSampleLoad()
	if err != nil {
		log.Debug(err)
		if this.Status == nil {
			this.Status = []string{}
		}
		this.Status = append(this.Status, err.Error())
	}
	log.Debug(7)
	return nil
}

func (this *Hardware) GetHostname() error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}

	this.Hostname = host
	return nil
}

var gpuTypes map[string]string = map[string]string{
	"1b06": "GeForce GTX 1080 Ti",
	"1b80": "GeForce GTX 1080",
	"1b38": "Tesla P40",
	"1bb3": "Tesla P4",
	"15f7": "Tesla P100 PCIe 12GB",
	"15f8": "Tesla P100 PCIe 16GB",
	"15f9": "Tesla P100 SXM2 16GB",
	"1e04": "GeForce RTX 2080 Ti",
	"1d81": "TITAN V",
	"1db1": "Tesla V100 SXM2",
	"1db4": "Tesla V100 PCIe",
	"17c2": "GeForce GTX TITAN X",
	"1b00": "TITAN X",
	"1b02": "TITAN Xp",
}

func (this *Hardware) GetGpu() error {
	log.Debug("GetGpu start")
	gpus, err := ExecCommandString("lspci|grep NVIDIA|grep VGA|wc -l")
	if err != nil {
		return err
	}
	log.Debug(gpus, len(gpus), len(strings.TrimSpace(gpus)))
	this.TotalGpu = strings.TrimSpace(gpus)
	if this.TotalGpu == "0" {
		log.Debug("gpu num is 0")
		return errors.New("gpu num is 0")
	} else {
		this.Gpu = []Gpu{}
	}

	gpuDetail, err := ExecCommandString("lspci |grep NVIDIA|grep VGA|cut -d ' ' -f8-|cut -d '(' -f1")
	if err != nil {
		return err
	}

	for _, x := range strings.Split(gpuDetail, "\n") {
		if x != "" && x != " " {
			var tmp Gpu
			if value, ok := gpuTypes[strings.TrimSpace(x)]; ok {
				tmp = Gpu{Type: value}
			} else {
				tmp = Gpu{Type: strings.Replace(strings.Replace(strings.TrimSpace(x), "[", "", 1), "]", "", 1)}
			}
			this.Gpu = append(this.Gpu, tmp)
		}
	}
	log.Debug("GetGpu end")
	return nil
}

func (this *Hardware) GetDisk() error {
	log.Debug("GetDisk start")
	this.Disk = []Disk{}

	list, err := ExecCommandString("fdisk -l|grep -E 'Disk /dev/sd|Disk /dev/vd|磁盘 /dev/sd|磁盘 /dev/vd'|sed 's/：/: /g'|sed 's/，/, /g'|awk '{print $2,$5/1024/1024/1024}'|sed 's/://g'")
	if err != nil {
		return err
	}

	log.Debug(list, len(strings.Split(list, "\n")), strings.Split(list, "\n"))
	for _, x := range strings.Split(list, "\n") {
		if len(x) > 4 {
			tmp := strings.Split(x, " ")
			this.Disk = append(this.Disk, Disk{
				MountPath: tmp[0],
				Capacity:  tmp[1],
			})
		}
	}
	log.Debug("GetDisk end")
	return nil
}

// top - 17:38:06 up 1 day, 23:42,  5 users,  load average: 2.41, 2.46, 3.02
// %Cpu(s): 18.2 us, 15.6 sy,  0.0 ni, 66.2 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
// MiB Mem :  15923.9 total,    522.6 free,   9283.5 used,   6117.9 buff/cache
func (this *Hardware) GetSampleLoad() error {
	log.Debug("GetSample start")
	this.SampleLoad = LoadStatus{}

	infos, err := ExecCommandString("top -bn 1 -i -c|sed -n '1p;3p;4p;'|sed 's/+/ /g'")
	if err != nil {
		return err
	}

	log.Debug(infos, len(strings.Split(infos, "\n")))
	for _, x := range strings.Split(infos, "\n") {
		log.Debug("debug", x)
		if strings.Contains(x, "load average") {
			this.SampleLoad.Load = strings.TrimSpace(strings.Split(x, ",")[4])
		} else if strings.Contains(x, "%Cpu") {
			t1 := strings.Split(x, ",")
			// cpu使用率
			id, err := strconv.ParseFloat(strings.Split(strings.TrimSpace(t1[3]), " ")[0], 64)
			if err != nil {
				break
				return err
			}
			this.SampleLoad.CpuUsed = fmt.Sprintf("%.1f", 100.0-id)

			// iowait
			this.SampleLoad.DiskIoWait = strings.Split(strings.TrimSpace(t1[4]), " ")[0]
		}
		// else if strings.Contains(x, "Mem") {
		// 	t2 := strings.Split(x, ",")
		// 	total, err := strconv.ParseFloat(strings.TrimSpace(strings.Split(strings.TrimSpace(strings.Split(t2[0], ":")[1]), " ")[0]), 64)
		// 	if err != nil {
		// 		log.Error("total", total, err.Error())
		// 		break
		// 		return err
		// 	}

		// 	used, err := strconv.ParseFloat(strings.Split(strings.TrimSpace(t2[2]), " ")[0], 64)
		// 	if err != nil {
		// 		log.Error("used", used, err.Error())
		// 		break
		// 		return err
		// 	}
		// 	log.Debug("mem user", used, total)
		// 	this.SampleLoad.MemoryUsed = fmt.Sprintf("%.1f", used/total*100)
		// }
	}

	mems, err := ExecCommandString("free -m|grep -E 'Mem|内存'|awk '{printf (\"%.1f\\n\",(1-$7/$2)*100)}'")
	if err != nil {
		log.Error(err.Error(), mems)
		return err
	}
	this.SampleLoad.MemoryUsed = strings.TrimSpace(strings.Replace(mems, "\n", "", -1))
	return nil
}

func (this *Hardware) GetTotal() error {
	log.Debug("GetDisk total start")
	tdisk, err := ExecCommandString("fdisk -l|grep -E 'Disk /dev/sd|Disk /dev/vd|磁盘 /dev/sd|磁盘 /dev/vd'|sed 's/：/: /g'|sed 's/，/, /g'|awk '{sum+=$5} END {print sum/1024/1024/1024}'")
	if err != nil {
		return err
	}
	this.TotalDisk = strings.TrimSpace(tdisk)

	tmem, err := ExecCommandString("free -g|sed -n '2p'|awk '{print $2}'")
	if err != nil {
		return err
	}
	this.TotalMemory = strings.TrimSpace(tmem)
	log.Debug("GetDisk total end")
	return nil
}

// 获取mac和ip
func (this *Hardware) GetNet() {
	log.Debug("GetNet start")
	ips := GetIPs()
	addrs := GetMacAddrs()

	this.Net = []NetInterface{}

	// fmt.Println("addrs", addrs, ips)
	// Warning: mac is error if machine has virtual ip
	for n, ip := range ips {
		if n <= len(addrs)-1 {
			t := strings.Split(addrs[n], ",")
			tmp := NetInterface{
				Ip:      ip,
				Macaddr: t[0],
				Name:    t[1],
			}
			this.Net = append(this.Net, tmp)
		}
	}

	log.Debug("GetNet end")
}

// 获取基础硬件信息
// Manufacturer string
// Type         string // Product Name: PowerEdge R730
// Sn           string
// Os           string // 系统版本 CentOs Linux
// OsPlatm      string // 系统架构 x86_64
// OsVersion    string // 系统爸爸号 7.2.1511
func (this *Hardware) GetMachine() error {
	log.Debug("GetMachine start")
	// m t s
	rs, err := ExecCommandString("dmidecode -t system|grep -E 'Manufacturer|Product|Serial'|cut -d ':' -f2")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	rsTmp := strings.Split(rs, "\n")
	log.Debug(rs, rsTmp, len(rsTmp))
	if len(rsTmp) < 3 {
		return errors.New("确认是否有权限执行dmidecode")
	}
	this.Machine.Manufacturer = strings.TrimSpace(rsTmp[0])
	this.Machine.Type = strings.TrimSpace(rsTmp[1])
	this.Machine.Sn = strings.TrimSpace(rsTmp[2])

	plat, err := ExecCommandString("uname -s -p -r")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	rsTmp1 := strings.Split(plat, " ")
	this.Machine.OsPlatm = fmt.Sprintf("%s %s", rsTmp1[0], strings.TrimSpace(rsTmp1[2]))
	this.Machine.OsCoreVersion = strings.TrimSpace(rsTmp1[1])

	oss, err := ExecCommandString("cat /etc/os-release|sed 's/\"//g'|grep -E \"^NAME=|^VERSION=\"")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	rsTmp2 := strings.Split(oss, "\n")
	this.Machine.Os = strings.Split(rsTmp2[0], "=")[1]
	this.Machine.OsVersion = strings.Split(rsTmp2[1], "=")[1]

	log.Debug("GetMachine end")
	return nil
}

// 获取CPU配置
/*
// # 总核数 = 物理CPU个数 X 每颗物理CPU的核数
// # 总逻辑CPU数 = 物理CPU个数 X 每颗物理CPU的核数 X 超线程数
// # 查看物理CPU个数
// cat /proc/cpuinfo| grep "physical id"| sort| uniq| wc -l
// # 查看每个物理CPU中core的个数(即核数)
// cat /proc/cpuinfo| grep "cpu cores"| uniq
// # 查看逻辑CPU的个数
// cat /proc/cpuinfo| grep "processor"| wc -l
type Cpu struct {
	Type        string // 型号
	CacheSize   string // 缓存
	PhysicalNum string //  物理cpu个数
	CoreNum     string // 核数
	Processor   string // 逻辑cpu个数
}
*/
func (this *Hardware) GetCpu() error {
	log.Debug("GetCPu start")
	rs, err := ExecCommandString("cat /proc/cpuinfo|grep -E 'physical id|model name|cpu cores|processor|cache size'|sort|uniq")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	phyn := 0
	processor := 0
	for _, x := range strings.Split(rs, "\n") {
		if strings.Contains(x, "cache size") {
			this.Cpu.CacheSize = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "cpu cores") {
			this.Cpu.CoreNum = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "model name") {
			this.Cpu.Type = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "cpu cores") {
			this.Cpu.CoreNum = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "processor") {
			processor++
		} else if strings.Contains(x, "physical id") {
			phyn++
		}
	}

	this.Cpu.Processor = fmt.Sprintf("%d", processor)
	this.Cpu.PhysicalNum = fmt.Sprintf("%d", phyn)
	log.Debug("GetCPu end")

	return nil
}

/*
type Memory struct {
	Capacity     int64  // 单卡容量
	Manufacturer string // 厂商
	Type         string // 型号
	SerialNumber string // 序列号
	Speed        string // 2400 MHz
	PartNumber   string
}
*/
func (this *Hardware) GetMemory() error {
	log.Debug("GetMemory start")
	this.Memory = []Memory{}

	rs, err := ExecCommandString("dmidecode -t memory|grep -E 'Size|Manufacturer|Serial Number|Part Number|Speed:|Type:|Locator:'|grep -Ev 'Bank|Unknow|No Module|Configured|Not Specified|Error'|sort|uniq")
	if err != nil {
		return err
	}

	tmp := strings.Split(rs, "\n")
	sn := []string{}
	la := []string{}

	var mf, pn, size, speed, tp string

	// fmt.Println("sn la", sn, la)
	for _, x := range tmp {
		if strings.Contains(x, "Manufacturer") {
			mf = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "Locator") {
			la = append(la, strings.TrimSpace(strings.Split(x, ":")[1]))
		} else if strings.Contains(x, "Part Number") {
			pn = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "Serial Number") {
			sn = append(sn, strings.TrimSpace(strings.Split(x, ":")[1]))
		} else if strings.Contains(x, "Size") {
			size = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "Speed") {
			speed = strings.TrimSpace(strings.Split(x, ":")[1])
		} else if strings.Contains(x, "Type") {
			tp = strings.TrimSpace(strings.Split(x, ":")[1])
		}
	}

	log.Debug("sn", sn, la)
	if len(sn) > 0 {
		for _, y := range sn {
			tt := Memory{
				Capacity:     size,
				Manufacturer: mf,
				SerialNumber: y,
				Type:         tp,
				Speed:        speed,
				PartNumber:   pn,
			}
			this.Memory = append(this.Memory, tt)
		}
	} else {
		log.Debug("la running", la)
		for _, y := range la {
			tt := Memory{
				Capacity:     size,
				Manufacturer: mf,
				SerialNumber: y,
				Type:         tp,
				Speed:        speed,
				PartNumber:   pn,
			}
			this.Memory = append(this.Memory, tt)
		}
	}

	log.Debug("GetMemory end")

	return nil
}
