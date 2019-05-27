package k8s

import (
	"encoding/json"
	"fmt"
	"strings"

	// "gitlab.yc/ares/k8sApi/models"
	. "github.com/lflxp/showme/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSecretList() (*apiv1.SecretList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	sl, err := cli.CoreV1().Secrets("").List(metav1.ListOptions{})
	return sl, err
}

func GetSecretByName(namespace, name string) (*apiv1.Secret, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	limit, err := cli.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	return limit, err
}

func GetSecretByLables(namespace, labels string) (*apiv1.SecretList, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	options := metav1.ListOptions{
		LabelSelector: labels,
	}

	limit, err := cli.CoreV1().Secrets(namespace).List(options)
	return limit, err
}

func DeleteSecret(namespace, name string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	err = cli.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}

type DockerConfigJson struct {
	Namespace   string            `namespace`
	Name        string            `name`
	Password    string            `password`
	Server      string            `server`
	Annotations map[string]string `annotations`
	Labels      map[string]string `labels`
}

// {"auths":{"(server)":{"username":"***","password":"***","email":"***","auth":"base64(***)"}}}
func (this DockerConfigJson) GetAuth() []byte {
	// var result string
	tmp := map[string]interface{}{}
	if this.Namespace != "" && this.Name != "" && this.Password != "" {
		t := map[string]interface{}{
			this.Server: map[string]string{
				"username": this.Name,
				"password": this.Password,
				"auth":     EncodeBase64(fmt.Sprintf("%s:%s", strings.TrimSpace(this.Name), strings.TrimSpace(this.Password))),
			},
		}
		tmp["auths"] = t
		// result = fmt.Sprintf("{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"auth\":\"%s\"}}}", this.Server, this.Name, this.Password, EncodeBase64(fmt.Sprintf("%s:%s", strings.TrimSpace(this.Name), strings.TrimSpace(this.Password))))
		// beego.Critical(result, this)
	}
	rs, err := json.Marshal(tmp)
	if err != nil {
		// beego.Critical(err.Error())
		return []byte(err.Error())
	}
	// beego.Critical(string(rs))
	// result = EncodeBase64(string(rs))
	// beego.Critical(result)
	// return result
	return rs
}

func CreateSecretWithDockerConfigJson(info DockerConfigJson) (*apiv1.Secret, error) {
	cli, err := GetClientSet()
	if err != nil {
		return nil, err
	}

	secrets := apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s.%s", Jiami(info.Name), strings.Replace(info.Server, ".", "-", -1)),
			Namespace:   info.Namespace,
			Labels:      info.Labels,
			Annotations: info.Annotations,
		},
		Data: map[string][]byte{".dockerconfigjson": info.GetAuth()},
		Type: apiv1.SecretTypeDockerConfigJson,
	}
	// beego.Critical(secrets)
	cf, err := cli.CoreV1().Secrets(info.Namespace).Create(&secrets)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			cf, err = cli.CoreV1().Secrets(info.Namespace).Update(&secrets)
		}
	}
	return cf, err
}

func DeleteSecretWithDockerConfigJson(namespace, username, server string) error {
	cli, err := GetClientSet()
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s.%s", Jiami(username), strings.Replace(server, ".", "-", -1))
	err = cli.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}

// func UpdateSecret(data *models.LimitRanged) (*apiv1.LimitRange, error) {
// 	cli, err := GetClientSet()
// 	if err != nil {
// 		return nil, err
// 	}

// 	lr, err := cli.CoreV1().LimitRanges(data.Metadata.Namespace).Update(data.Resolved())
// 	return lr, err
// }
