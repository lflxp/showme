//----------------------------------------
//
// Copyright © ying32. All Rights Reserved.
//
// Licensed under Apache License 2.0
//
//----------------------------------------

// +build !windows

package vcl

import (
	"fmt"
)

func showError(err interface{}) {
	fmt.Println(err)
}

func tryLoadAppIcon() {
	// no code
}
