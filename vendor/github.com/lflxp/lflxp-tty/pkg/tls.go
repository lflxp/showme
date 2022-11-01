package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var (
	PrivateKey string
	PublicKey  string
)

// 加密RSA
func RsaEncrypt(origData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode([]byte(PublicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密RSA
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode([]byte(PrivateKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return data, nil
}

//生成RSA私钥和公钥，保存到文件中
func GenerateRSAKey(bits int) error {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	///Users/libin/lbdata
	privateFile, err := os.Create("./server.crt")
	if err != nil {
		return err
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "CERTIFICATE", Bytes: X509PrivateKey}
	//将数据保存到文件
	pem.Encode(privateFile, &privateBlock)
	PrivateKey = string(pem.EncodeToMemory(&privateBlock))

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create("./server.key")
	if err != nil {
		return err
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "PRIVATE KEY", Bytes: X509PublicKey}
	//保存到文件
	pem.Encode(publicFile, &publicBlock)
	PublicKey = string(pem.EncodeToMemory(&publicBlock))
	return nil
}
