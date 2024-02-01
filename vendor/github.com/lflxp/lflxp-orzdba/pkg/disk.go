package pkg

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/disk"
)

var beforeDisk *MonitorDisk

func init() {
	beforeDisk, err = NewDisk()
}

type MonitorDisk struct {
	ReadCount        uint64 `json:"readCount"`
	MergedReadCount  uint64 `json:"mergedReadCount"`
	ReadBytes        uint64 `json:"readBytes"`
	ReadTime         uint64 `json:"readTime"`
	WriteCount       uint64 `json:"writeCount"`
	MergedWriteCount uint64 `json:"mergedWriteCount"`
	WriteBytes       uint64 `json:"writeBytes"`
	WriteTime        uint64 `json:"writeTime"`
	IopsInProgress   uint64 `json:"iopsInProgress"`
	IoTime           uint64 `json:"ioTime"`
	WeightedIO       uint64 `json:"weightedIO"`
}

// get all disk
func (this *MonitorDisk) Get() error {
	data, err := disk.Partitions(true)
	if err != nil {
		return err
	}

	rs := []string{}
	for _, x := range data {
		// fmt.Println(x.Device, x.Fstype, x.Mountpoint, x.Opts)
		rs = append(rs, x.Device)
	}

	dd, err := disk.IOCounters()
	if err != nil {
		return err
	}

	for _, v := range dd {
		this.ReadBytes += v.ReadBytes
		this.MergedReadCount += v.MergedReadCount
		this.ReadBytes += v.ReadBytes
		this.ReadTime += v.ReadTime
		this.WriteCount += v.WriteCount
		this.MergedWriteCount += v.MergedWriteCount
		this.WriteBytes += v.WriteBytes
		this.WriteTime += v.WriteTime
		this.IopsInProgress += v.IopsInProgress
		this.IoTime += v.IoTime
		this.WeightedIO += v.WeightedIO
	}
	return nil
}

func NewDisk() (*MonitorDisk, error) {
	data := &MonitorDisk{}
	err := data.Get()
	return data, err
}

func DiskInfo() (string, error) {
	var data_detail string
	after, err := NewDisk()
	if err != nil {
		return data_detail, err
	}

	rs_disk := float64(after.ReadCount-beforeDisk.ReadCount) / 0.9999
	// fmt.Printf("ws_disk is float64(%d-%d)/0.999\n", after.WriteCount, beforeDisk.WriteCount)
	ws_disk := float64(after.WriteCount-beforeDisk.WriteCount) / 0.9999

	// fmt.Printf("rkbs_disk is float64(%d-%d)/1.999\n", after.ReadBytes, beforeDisk.ReadBytes)
	rkbs_disk := float64(after.ReadBytes-beforeDisk.ReadBytes) / 1.9999
	// fmt.Printf("wkbs_disk is float64(%d-%d)/1.999\n", after.WriteBytes, beforeDisk.WriteBytes)
	// fmt.Println(rkbs_disk, floatToString(rkbs_disk, 1), 9-len(floatToString(rkbs_disk, 1)))
	wkbs_disk := float64(after.WriteBytes-beforeDisk.WriteBytes) / 1.9999
	// fmt.Println(wkbs_disk, floatToString(wkbs_disk, 1), 9-len(floatToString(wkbs_disk, 1)))
	queue_disk := fmt.Sprintf("%d", after.IopsInProgress)

	var await_disk float64
	var svctm_disk float64
	if (rs_disk + ws_disk) == 0.0 {
		await_disk = float64(after.ReadTime+after.WriteTime-beforeDisk.ReadTime-beforeDisk.WriteTime) / (rs_disk + ws_disk + 1)
		svctm_disk = float64(after.IoTime-beforeDisk.IoTime) / (rs_disk + ws_disk + 1)
	} else {
		await_disk = float64(after.ReadTime+after.WriteTime-beforeDisk.ReadTime-beforeDisk.WriteTime) / (rs_disk + ws_disk)
		svctm_disk = float64(after.IoTime-beforeDisk.IoTime) / (rs_disk + ws_disk)
	}

	util_disk := float64(after.IoTime-beforeDisk.IoTime) / 10
	//usr
	// fmt.Println(rs_disk, ws_disk, rkbs_disk, wkbs_disk, queue_disk, await_disk, svctm_disk, util_disk)
	// fmt.Println(strings.Repeat(" ", 6-len(floatToString(rs_disk, 1))) + floatToString(rs_disk, 1))
	if 1 != 1 {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(rs_disk, 1)))+floatToString(rs_disk, 1), "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(floatToString(rs_disk, 1), 6), "red", "", false, true)
	} else {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(rs_disk, 1)))+floatToString(rs_disk, 1), "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(rs_disk, 1), 6), "", "", false, false)
	}

	if 1 != 1 {
		// data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(ws_disk, 1)))+floatToString(ws_disk, 1), "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(floatToString(ws_disk, 1), 7), "red", "", false, true)
	} else {
		// data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(ws_disk, 1)))+floatToString(ws_disk, 1), "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(ws_disk, 1), 7), "", "", false, false)
	}

	if rkbs_disk/1024/1024 >= 1.0 {
		// data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(rkbs_disk/1024/1024, 1)))+floatToString(rkbs_disk/1024/1024, 1)+"M", "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(floatToString(rkbs_disk/1024/1024, 1), 7)+"M", "red", "", false, true)
	} else if rkbs_disk/1024 < 1.0 {
		// data_detail += Colorize(strings.Repeat(" ", 8-len(floatToString(rkbs_disk, 1)))+floatToString(rkbs_disk, 1), "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(rkbs_disk, 1), 8), "", "", false, false)
	} else if rkbs_disk/1024/1024 < 1.0 && rkbs_disk/1024 > 1.0 {
		// data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(rkbs_disk/1024, 1)))+floatToString(rkbs_disk/1024, 1)+"K", "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(rkbs_disk/1024, 1), 7)+"K", "", "", false, false)
	}

	if wkbs_disk/1024/1024 > 1.0 {
		// data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(wkbs_disk/1024/1024, 1)))+floatToString(wkbs_disk/1024/1024, 1)+"M", "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(floatToString(wkbs_disk/1024/1024, 1), 7)+"M", "red", "", false, true)
	} else if wkbs_disk/1024 < 1.0 {
		// data_detail += Colorize(strings.Repeat(" ", 8-len(floatToString(wkbs_disk, 1)))+floatToString(wkbs_disk, 1), "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(wkbs_disk, 1), 8), "", "", false, false)
	} else {
		// data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(wkbs_disk/1024, 1)))+floatToString(wkbs_disk/1024, 1)+"K", "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(wkbs_disk/1024, 1), 7)+"K", "", "", false, false)
	}

	if after.IopsInProgress > 10 {
		// data_detail += Colorize(strings.Repeat(" ", 4-len(queue_disk))+queue_disk+".0 ", "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(queue_disk, 4)+".0 ", "red", "", false, true)
	} else {
		// data_detail += Colorize(strings.Repeat(" ", 4-len(queue_disk))+queue_disk+".0 ", "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(queue_disk, 4)+".0 ", "", "", false, false)
	}

	if await_disk > 5.0 {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(await_disk, 1)))+floatToString(await_disk, 1), "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(floatToString(await_disk, 1), 6), "red", "", false, true)
	} else {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(await_disk, 1)))+floatToString(await_disk, 1), "green", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(await_disk, 1), 6), "green", "", false, false)
	}

	if svctm_disk > 5.0 {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(svctm_disk, 1)))+floatToString(svctm_disk, 1), "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(floatToString(svctm_disk, 1), 6), "red", "", false, true)
	} else {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(svctm_disk, 1)))+floatToString(svctm_disk, 1), "", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(svctm_disk, 1), 6), "", "", false, false)
	}

	if util_disk > 80.0 {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(util_disk, 1)))+floatToString(util_disk, 1), "red", "", false, true)
		data_detail += Colorize(parseRepeatSpace(floatToString(util_disk, 1), 6), "red", "", false, true)
	} else if util_disk > 100.0 {
		data_detail += Colorize(" 100.0", "green", "", false, false)
	} else {
		// data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(util_disk, 1)))+floatToString(util_disk, 1), "green", "", false, false)
		data_detail += Colorize(parseRepeatSpace(floatToString(util_disk, 1), 6), "green", "", false, false)
	}

	data_detail += Colorize("|", "dgreen", "", false, false)
	beforeDisk = after
	return data_detail, nil
}
