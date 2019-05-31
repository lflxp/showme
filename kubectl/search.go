package kubectl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
)

func searchBuffer(g *gocui.Gui, v *gocui.View) error {
	origin.BeforeSearch = v.Name()
	maxX, maxY := g.Size()
	if v, err := g.SetView("searchBuffer", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = fmt.Sprintf("Search %s", origin.BeforeSearch)
		v.Highlight = true
		v.Editable = true
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		if _, err := g.SetCurrentView("searchBuffer"); err != nil {
			return err
		}
	}

	return nil
}

func delsearchBuffer(g *gocui.Gui, v *gocui.View) error {
	// getline
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	// getResult
	if v, err := g.View(origin.BeforeSearch); err != nil {
		return err
	} else {
		// v.Autoscroll = true
		// v.Highlight = true
		// v.Wrap = true
		// deleteOrigin
		if err := g.DeleteView("searchBuffer"); err != nil {
			return err
		}
		// color for search word and replace before result
		tmprs := v.BufferLines()
		v.Clear()
		front := []string{}
		backend := []string{}
		title := ""
		countTitle := 0
		for _, x := range tmprs {
			x = strings.Replace(x, ">", "*", -1)
			match, _ := regexp.MatchString(l, x)
			if match {
				// fmt.Fprintln(v, utils.Colorize(strings.Replace(x, " ", ">", 3), "yellow", "", false, true))
				front = append(front, utils.Colorize(strings.Replace(x, "*", ">", 1), "yellow", "", false, true))
			} else {
				match1, _ := regexp.MatchString("NAME", x)
				if match1 && countTitle == 0 {
					// fmt.Fprintln(v, utils.Colorize(x, "dgreen", "black", false, false))
					title = utils.Colorize(x, "dgreen", "black", true, false)
					countTitle++
				} else {
					// fmt.Fprintln(v, strings.Replace(x, ">", " ", -1))
					backend = append(backend, x)
				}
			}
		}
		// print
		fmt.Fprintln(v, title)
		for _, o := range front {
			fmt.Fprintln(v, o)
		}
		for _, t := range backend {
			fmt.Fprintln(v, t)
		}
		// setTop
		if _, err := g.SetCurrentView(origin.BeforeSearch); err != nil {
			return err
		}
	}

	return nil
}
