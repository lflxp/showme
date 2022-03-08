// Copyright 2020 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package extended contains an extended set of terminal descriptions.
// Applications desiring to have a better chance of Just Working by
// default should include this package.  This will significantly increase
// the size of the program.
package extended

import (
	// The following imports just register themselves --
	// these are the terminal types we aggregate in this package.
	_ "gitee.com/lflxp/tcell/v2/terminfo/a/aixterm"
	_ "gitee.com/lflxp/tcell/v2/terminfo/a/alacritty"
	_ "gitee.com/lflxp/tcell/v2/terminfo/a/ansi"
	_ "gitee.com/lflxp/tcell/v2/terminfo/b/beterm"
	_ "gitee.com/lflxp/tcell/v2/terminfo/c/cygwin"
	_ "gitee.com/lflxp/tcell/v2/terminfo/d/dtterm"
	_ "gitee.com/lflxp/tcell/v2/terminfo/e/emacs"
	_ "gitee.com/lflxp/tcell/v2/terminfo/f/foot"
	_ "gitee.com/lflxp/tcell/v2/terminfo/g/gnome"
	_ "gitee.com/lflxp/tcell/v2/terminfo/h/hpterm"
	_ "gitee.com/lflxp/tcell/v2/terminfo/k/konsole"
	_ "gitee.com/lflxp/tcell/v2/terminfo/k/kterm"
	_ "gitee.com/lflxp/tcell/v2/terminfo/l/linux"
	_ "gitee.com/lflxp/tcell/v2/terminfo/p/pcansi"
	_ "gitee.com/lflxp/tcell/v2/terminfo/r/rxvt"
	_ "gitee.com/lflxp/tcell/v2/terminfo/s/screen"
	_ "gitee.com/lflxp/tcell/v2/terminfo/s/simpleterm"
	_ "gitee.com/lflxp/tcell/v2/terminfo/s/sun"
	_ "gitee.com/lflxp/tcell/v2/terminfo/t/termite"
	_ "gitee.com/lflxp/tcell/v2/terminfo/t/tmux"
	_ "gitee.com/lflxp/tcell/v2/terminfo/v/vt100"
	_ "gitee.com/lflxp/tcell/v2/terminfo/v/vt102"
	_ "gitee.com/lflxp/tcell/v2/terminfo/v/vt220"
	_ "gitee.com/lflxp/tcell/v2/terminfo/v/vt320"
	_ "gitee.com/lflxp/tcell/v2/terminfo/v/vt400"
	_ "gitee.com/lflxp/tcell/v2/terminfo/v/vt420"
	_ "gitee.com/lflxp/tcell/v2/terminfo/v/vt52"
	_ "gitee.com/lflxp/tcell/v2/terminfo/w/wy50"
	_ "gitee.com/lflxp/tcell/v2/terminfo/w/wy60"
	_ "gitee.com/lflxp/tcell/v2/terminfo/w/wy99_ansi"
	_ "gitee.com/lflxp/tcell/v2/terminfo/x/xfce"
	_ "gitee.com/lflxp/tcell/v2/terminfo/x/xterm"
	_ "gitee.com/lflxp/tcell/v2/terminfo/x/xterm_kitty"
	_ "gitee.com/lflxp/tcell/v2/terminfo/x/xterm_termite"
)
