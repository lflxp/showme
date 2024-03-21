// https://colobu.com/2021/01/17/go-embed-tutorial/
package asset

import _ "embed"

//go:embed face.xml
var FaceXml []byte
