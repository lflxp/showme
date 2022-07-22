package sync

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

/*
	断点续传：
		文件传递：文件复制
		srcFile 复制到 destFile
	思路：
		边复制，边记录复制的总量（记录总复制字节数保存至tempFile）
		复制一半中断后，若再次复制，会先读取tempFile文件里记录的已完成字节总数
		然后将这个字节数作为文件读写offset（偏移量），
		来结合whence（偏移位置）来决定后面读写文件的起点位置。
		seek(offset,whence),设置指针光标的位置
		第一个参数：偏移量
		第二个参数：如何设置
			0：seekStart表示相对于文件开始，
			1：seekCurrent表示相对于当前偏移量，
			2：seek end表示相对于结束。
*/

// 本地单文件复制
func LocalSync(src, dest string, debug bool) error {
	tmpFile := fmt.Sprintf("%s.temp", src)

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	destFile, _ := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	tempFile, _ := os.OpenFile(tmpFile, os.O_CREATE|os.O_RDWR, os.ModePerm)

	defer srcFile.Close()
	defer destFile.Close()

	//1.读取临时文件中的数据，根据seek
	tempFile.Seek(0, io.SeekStart)
	bs := make([]byte, 100, 100)
	n1, err := tempFile.Read(bs)
	if err != nil {
		log.Printf("tempFile.Read %s\n", err.Error())
	}
	countStr := string(bs[:n1])
	count, _ := strconv.ParseInt(countStr, 10, 64)
	fmt.Printf("开始复制: %s start %d end %d \n", src, n1, count)

	// 2. 设置读，写的偏移量
	srcFile.Seek(count, 0)  // 设置file1下一次读或者写的起点
	destFile.Seek(count, 0) // 设置file2下一次读或者写的起点
	data := make([]byte, 1000, 1000)

	n2 := -1            // 读取的数据量， 默认值
	n3 := -1            //写出的数据量
	total := int(count) //读取的总量

	srcFileSize, _ := srcFile.Stat()

	for {
		//3.读取数据
		// 基于上面的起点，读取file1文件len（data）个字节，
		// n2：实际读取的字节数（小于等于len（data），将读取的字节存入data。
		n2, err = srcFile.Read(data)
		if err == io.EOF {
			fmt.Printf("文件: %s Size: %d 复制完毕。。\n", src, total)
			tempFile.Close()
			os.Remove(tmpFile) //复制完，先不删除，验证最后存储的字节总数是不是跟复制的文件大小一致。
			break
		}
		//将数据写入到目标文件
		// 基于上面的起点，向file2文件写入len（data[:n2]）个字节，也就是写入data中前n2个元素；
		// n3：实际写入的字节数（小于等于n2）。
		n3, _ = destFile.Write(data[:n2])
		total += n3

		process := (float64(total) / float64(srcFileSize.Size())) * 100
		// if int64(total)/srcFileSize.Size()*100%50 == 0 {
		if debug && int(process)%10 == 0 {
			fmt.Printf("复制文件: %s 进度: %.2f %d Size: %d \n", src, process, int64(process), total)
		}
		// }

		//将复制总量，存储到临时文件中
		tempFile.Seek(0, io.SeekStart) // 设置file3的下次读写起点为源点：0点，即覆盖重写。
		tempFile.WriteString(strconv.Itoa(total))

		//假装断电
		//if total>1800000{
		//  panic("假装断电了。。。，假装的。。。")
		//}
	}
	return nil
}
