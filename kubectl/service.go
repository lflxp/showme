package kubectl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/lflxp/showme/utils/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KeyService(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, Services); err != nil {
		return err
	}
	if err := g.SetKeybinding("Serviceed", gocui.KeyEnter, gocui.ModNone, getServices); err != nil {
		return err
	}
	if err := g.SetKeybinding("msgservice", gocui.KeyEnter, gocui.ModNone, delservicemessage); err != nil {
		return err
	}
	if err := g.SetKeybinding("Serviceed", gocui.KeyDelete, gocui.ModNone, deleteServiceView); err != nil {
		return err
	}
	if err := g.SetKeybinding("deleteServiceView", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		return err
	}
	return nil
}

func deleteServiceView(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	rs := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	maxX, maxY := g.Size()
	if v, err := g.SetView("deleteServiceView", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		origin.CurrentPod = strings.TrimSpace(rs[2])
		origin.DefaultNS = strings.TrimSpace(rs[1])
		v.Title = fmt.Sprintf("确认删除[Service] %s:%s?(y/N)", strings.TrimSpace(rs[2]), strings.TrimSpace(rs[1]))
		v.Highlight = true
		v.Editable = true
		// v.Frame = false
		// v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorRed
		// fmt.Fprintln(v, strings.Trim(l, " "))
		// fmt.Fprintln(v, l)
		// selectId = strings.Trim(l, " ")
		// fmt.Fprintln(v, fmt.Sprintf("Your Selectd Range: %s", l))
		if _, err := g.SetCurrentView("deleteServiceView"); err != nil {
			return err
		}
	}
	return nil
}

func delservicemessage(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msgservice"); err != nil {
		return err
	}
	if _, err := setCurrentViewOnTop(g, "Serviceed"); err != nil {
		return err
	}
	return nil
}

func getServices(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	ttt := strings.Split(strings.Replace(l, ">", "*", 1), "*")
	if len(ttt) > 1 {
		maxX, maxY := g.Size()
		if v, err := g.SetView("msgservice", maxX*8/100, maxY*8/100, maxX*92/100, maxY*92/100); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			namespace := strings.TrimSpace(ttt[1])
			name := strings.TrimSpace(ttt[2])
			v.Title = fmt.Sprintf("Current: %s %s", namespace, name)
			v.Highlight = true
			v.SelFgColor = gocui.ColorMagenta
			v.SelBgColor = gocui.ColorCyan
			v.Editable = true
			v.Wrap = true
			// fmt.Fprintln(v, strings.Trim(l, " "))
			// fmt.Fprintln(v, l)
			// selectId = strings.Trim(l, " ")

			pod, err := origin.ClientSet.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				fmt.Fprintln(v, err.Error())
			} else {
				// json格式美化
				b, err := json.MarshalIndent(pod, "", "\t")
				if err != nil {
					fmt.Fprintln(v, err.Error())
				} else {
					fmt.Fprintln(v, utils.Colorize(string(b), "red", "", false, true))
				}
			}

			if _, err := g.SetCurrentView("msgservice"); err != nil {
				return err
			}

		}
	}

	return nil
}

func Services(g *gocui.Gui, v *gocui.View) error {
	fmt.Println("Servicesssss")
	if err = delOtherViewNoBack(g); err != nil {
		return err
	}
	if v, err := g.SetView("Serviceed", 0, 0, origin.maxX-1, origin.maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		tmp_services, err := origin.ClientSet.CoreV1().Services("").List(metav1.ListOptions{})
		if err != nil {
			return err
		}

		v.Title = fmt.Sprintf("Service/%d", len(tmp_services.Items))
		v.Highlight = true
		v.Editable = true
		// v.Wrap = true
		// v.MoveCursor(startx, endy, false)
		if _, err := setCurrentViewOnTop(g, "Serviceed"); err != nil {
			return err
		}

		tableNow := table.NewTable(origin.maxX - 1)

		// NAMESPACE     NAME                   TYPE        CLUSTER-IP     EXTERNAL-IP    PORT(S)                  AGE   SELECTOR
		// default       deploy-heketi          ClusterIP   10.68.96.124   10.128.25.68   8080/TCP                 16d   name=deploy-heketi
		// default       kubernetes             ClusterIP   10.68.0.1      <none>         443/TCP                  16d   <none>
		// kube-system   kube-dns               ClusterIP   10.68.0.2      <none>         53/UDP,53/TCP,9153/TCP   16d   k8s-app=kube-dns
		// kube-system   kubernetes-dashboard   NodePort    10.68.19.163   <none>         443:28537/TCP            16d   k8s-app=kubernetes-dashboard
		// kube-system   metrics-server         ClusterIP   10.68.43.4     <none>         443/TCP                  16d   k8s-app=metrics-server
		// kube-system   tiller-deploy          ClusterIP   10.68.5.143    <none>         44134/TCP                16d   app=helm,name=tiller
		tableNow.AddCol("Namespace").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("NAME").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("TYPE").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("CLUSTER-IP").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("EXTERNAL-IP").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("PORT(S)").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("AGE").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.AddCol("SELECTOR").SetColor("dgreen").SetTextAlign(table.TextLeft).SetBgColor("black")
		tableNow.CalColumnWidths()

		for num, value := range tmp_services.Items {
			if num == 0 {
				tableNow.FprintHeader(v)
			}

			name := table.NewCol()
			name.Data = fmt.Sprintf("*%s", value.Namespace)
			name.TextAlign = table.TextLeft
			name.Color = "yellow"
			tableNow.AddRow(0, name)

			ns := table.NewCol()
			ns.Data = fmt.Sprintf("*%s", value.Name)
			ns.TextAlign = table.TextLeft
			ns.Color = "yellow"
			tableNow.AddRow(1, ns)

			ttype := table.NewCol()
			ttype.Data = fmt.Sprintf("*%s", value.Spec.Type)
			ttype.TextAlign = table.TextLeft
			ttype.Color = "yellow"
			tableNow.AddRow(2, ttype)

			cip := table.NewCol()
			cip.Data = fmt.Sprintf("%s", value.Spec.ClusterIP)
			cip.TextAlign = table.TextLeft
			cip.Color = "yellow"
			tableNow.AddRow(3, cip)

			eip := table.NewCol()
			eip.Data = fmt.Sprintf("%s", strings.Join(value.Spec.ExternalIPs, ","))
			eip.TextAlign = table.TextLeft
			eip.Color = "yellow"
			tableNow.AddRow(4, eip)

			pp := []string{}
			for _, x := range value.Spec.Ports {
				pp = append(pp, fmt.Sprintf("%d/%s", x.Port, x.Protocol))
			}
			ports := table.NewCol()
			ports.Data = fmt.Sprintf("%s", strings.Join(pp, ","))
			ports.TextAlign = table.TextLeft
			ports.Color = "yellow"
			tableNow.AddRow(5, ports)

			age := table.NewCol()
			age.Data = fmt.Sprintf("%s", strings.Replace(fmt.Sprintf("%v", value.CreationTimestamp.Sub(time.Now())), "-", "", -1))
			age.TextAlign = table.TextLeft
			age.Color = "yellow"
			tableNow.AddRow(6, age)

			kk := []string{}
			for k, v := range value.Spec.Selector {
				kk = append(kk, fmt.Sprintf("%s:%s", k, v))
			}
			selector := table.NewCol()
			selector.Data = fmt.Sprintf("%s", strings.Join(kk, ","))
			selector.TextAlign = table.TextLeft
			selector.Color = "yellow"
			tableNow.AddRow(7, selector)

			// fmt.Fprintln(w, s)
		}
		tableNow.Fprint(v)
	}
	return nil
}
