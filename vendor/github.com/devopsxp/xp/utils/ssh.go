// https://www.cnblogs.com/chenqionghe/p/8267326.html

package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	privateKey = fmt.Sprintf("%s/.ssh/id_rsa", home)
}

var privateKey string

type Cli struct {
	IP         string       //IP地址
	Username   string       //用户名
	Password   string       //密码
	Port       int          //端口号
	client     *ssh.Client  //ssh客户端
	sftpClient *sftp.Client // sftp客户端
	LastResult string       //最近一次Run的结果
}

// 创建命令行对象
// @param ip IP地址
// @param username 用户名
// @param password 密码
// @param port 端口号,默认22
func New(ip string, username string, password string, port ...int) *Cli {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	return cli
}

// 执行shell
// @param shell shell脚本命令
func (c *Cli) Run(shell string) (string, error) {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return "", err
		}
	}
	defer c.client.Close()

	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}

// 连接
func (c *Cli) connect() error {
	var config ssh.ClientConfig
	if c.Password != "" {
		config = ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	} else {
		key, err := ioutil.ReadFile(privateKey)
		if err != nil {
			slog.Error("Unable to read private key: %v\n", err)
			return err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			slog.Error("Unable to parse private key: %v\n", err)
			return err
		}

		config = ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	}

	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}

// 上传文件String
func (c *Cli) SftpUploadTemplateString(data, remotepath string) error {
	if c.sftpClient == nil {
		if err := c.sftpConnect(); err != nil {
			return err
		}
	}
	defer c.sftpClient.Close()

	dstFile, err := c.sftpClient.Create(remotepath)
	if err != nil {
		slog.Error(fmt.Sprintf("创建文件 %s 失败: %s", remotepath, err.Error()))
		return err
	}
	defer dstFile.Close()

	// 按byte传输文件
	dstFile.Write([]byte(data))

	slog.Debug(fmt.Sprintf("模板文件 %s 上传成功", remotepath))
	return nil
}

// 上传文件
// 注意： remotepath是文件路径不是文件夹路径
func (c *Cli) SftpUploadToRemote(localpath, remotepath string) error {
	if c.sftpClient == nil {
		if err := c.sftpConnect(); err != nil {
			return err
		}
	}
	defer c.sftpClient.Close()

	// remoteFileName := path.Base(localpath)
	// slog.Warn(localpath, remoteFileName)
	srcFile, err := os.Open(localpath)
	if err != nil {
		slog.Error("打开文件失败", err)
		return err
	}
	defer srcFile.Close()

	dstFile, err := c.sftpClient.Create(remotepath)
	if err != nil {
		slog.Error("创建文件 %s 失败: %s", remotepath, err.Error())
		return err
	}
	defer dstFile.Close()

	// 按byte传输文件
	buffer := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				slog.Debug("已读取到文件末尾")
				break
			} else {
				slog.Debug(fmt.Sprintf("读取文件出错 %v", err))
				return err
			}
		}
		//注意，由于文件大小不定，不可直接使用buffer，否则会在文件末尾重复写入，以填充1024的整数倍
		dstFile.Write(buffer[:n])
	}
	slog.Debug(fmt.Sprintf("文件 %s 上传成功", remotepath))
	return nil
}

// 下载文件
func (c *Cli) SftpDownloadToLocal(localpath, remotepath string) error {
	if c.sftpClient == nil {
		if err := c.sftpConnect(); err != nil {
			return err
		}
	}
	defer c.sftpClient.Close()

	srcFile, err := c.sftpClient.Open(remotepath)
	if err != nil {
		slog.Error("文件读取失败", err)
		return err
	}
	defer srcFile.Close()

	// localFilename := path.Base(remotepath)
	dstFile, err := os.Create(localpath)
	if err != nil {
		slog.Error("文件创建失败", err.Error())
		return err
	}
	defer dstFile.Close()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		slog.Error("额外那件写入失败", err.Error())
		return err
	}

	slog.Warn("文件 %s 下载成功", localpath)
	return nil
}

// 利用sftp传输文件连接
// https://www.cnblogs.com/-xuzhankun/p/11056576.html
func (c *Cli) sftpConnect() error {

	var clientConfig *ssh.ClientConfig
	if c.Password != "" {
		clientConfig = &ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	} else {
		key, err := ioutil.ReadFile(privateKey)
		if err != nil {
			slog.Error("Unable to read private key: %v\n", err)
			return err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			slog.Error("Unable to parse private key: %v\n", err)
			return err
		}

		clientConfig = &ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	}

	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, clientConfig) //连接ssh
	if err != nil {
		slog.Error("连接ssh失败", err)
		return err
	}

	if sftpClient, err := sftp.NewClient(sshClient); err != nil { //创建客户端
		slog.Error("创建客户端失败", err)
		return err
	} else {
		c.sftpClient = sftpClient
	}

	return nil
}

// 执行带交互的命令
func (c *Cli) RunTerminal(shell string, stdout, stderr io.Writer) error {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}

	err = session.Run(shell)
	return err
}
