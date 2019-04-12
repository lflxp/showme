package utils

import (
	"fmt"
	"time"
)

func GetNowTime() string {
	return fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05"))
}

// 文字字体 参数介绍：text->文本内容 status->文字颜色 background->背景颜色 underline->是否下划线 highshow->是否高亮
// http://www.cnblogs.com/frydsh/p/4139922.html
func Colorize(text string, status string, background string, underline bool, highshow bool) string {
	out_one := "\033["
	out_two := ""
	out_three := ""
	out_four := ""
	//可动态配置字体颜色 背景色 高亮
	// 显示：0(默认)、1(粗体/高亮)、22(非粗体)、4(单条下划线)、24(无下划线)、5(闪烁)、25(无闪烁)、7(反显、翻转前景色和背景色)、27(无反显)
	// 颜色：0(黑)、1(红)、2(绿)、 3(黄)、4(蓝)、5(洋红)、6(青)、7(白)
	//  前景色为30+颜色值，如31表示前景色为红色；背景色为40+颜色值，如41表示背景色为红色。
	if underline == true && highshow == true {
		out_four = ";1;4m" //高亮
	} else if underline != true && highshow == true {
		out_four = ";1m"
	} else if underline == true && highshow != true {
		out_four = ";4m"
	} else {
		out_four = ";22m"
	}

	switch status {
	case "black":
		out_two = "30"
	case "red":
		out_two = "31"
	case "green":
		out_two = "32"
	case "yellow":
		out_two = "33"
	case "blue":
		out_two = "34"
	case "purple":
		out_two = "35"
	case "dgreen":
		out_two = "36"
	case "white":
		out_two = "37"
	default:
		out_two = ""
	}

	switch background {
	case "black":
		out_three = "40;"
	case "red":
		out_three = "41;"
	case "green":
		out_three = "42;"
	case "yellow":
		out_three = "43;"
	case "blue":
		out_three = "44;"
	case "purple":
		out_three = "45;"
	case "dgreen":
		out_three = "46;"
	case "white":
		out_three = "47;"
	default:
		out_three = ""
	}
	return out_one + out_three + out_two + out_four + text + "\033[0m"
}
