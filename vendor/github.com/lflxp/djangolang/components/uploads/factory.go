package uploads

import (
	"fmt"
	"log/slog"
	"strings"
)

// string转html
type Uploads struct {
}

func (u *Uploads) Transfer() interface{} {
	return func(data map[string]string) string {
		if err := u.check(data); err != nil {
			return err.Error()
		}

		var result string

		// Id:int64:id: Country:string:Country: Zoom:string:Zoom: Company:string:Company: Items:string:Items: Production:string:Production:
		// Count:string:Count: Serial:string:Serial: Extend:string:Extend: 上传文件:file:file:
		if value, ok := data["Col"]; !ok {
			slog.Error("%v not contains Col key", data)
			return fmt.Sprintf("%v not contains Col key", data)
		} else {
			for index, v := range strings.Split(strings.TrimSpace(value), " ") {
				tmp := strings.Split(v, ":")
				if len(tmp) < 3 {
					slog.Error("upload %v format error, length is %d", v, len(tmp))
				}

				if tmp[1] == "file" {
					result += fmt.Sprintf(`formData.append('%s', $('#myform')[%d].files[0];\n`, tmp[2], index)
				} else {
					result += fmt.Sprintf(`formData.append('%s', $('[name="%s"]').val());\n`, tmp[2], tmp[2])
				}
			}
		}
		// log.Debugf("====upload transfer is %s", result)
		return result
	}
}

func (u *Uploads) check(data map[string]string) error {
	// log.Debug("upload raw", data)
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}
	return nil
}
