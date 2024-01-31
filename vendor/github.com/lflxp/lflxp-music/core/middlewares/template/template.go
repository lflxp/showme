package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/lflxp/tools/orm/sqlite"
)

var FuncMap template.FuncMap = map[string]interface{}{
	"beegoli":      BeegoLi,
	"admincolumns": AdminColumns,
	"formcolumns":  FormColumns,
	"str2html":     str2html,
}

func AddFuncMap(funcName string, funcs interface{}) {
	if _, ok := FuncMap[funcName]; !ok {
		FuncMap[funcName] = funcs
	}
}

// https://docs.djangoproject.com/en/1.11/ref/contrib/admin/
// https://github.com/beego/beego/blob/fea7c914ccc789cfe0cd90afd32f1e9a0f3a90bb/server/web/templatefunc.go
// func init() {
// 	FuncMap = template.FuncMap{}
// 	FuncMap["beegoli"] = BeegoLi
// 	FuncMap["admincolumns"] = AdminColumns
// 	FuncMap["formcolumns"] = FormColumns
// 	FuncMap["str2html"] = str2html
// }

func str2html(raw string) template.HTML {
	return template.HTML(raw)
}

func BeegoLi(info []map[string]string) string {
	if len(info) == 0 {
		return "data is none"
	}

	// 按Service分组
	tmp := map[string][]map[string]string{}

	for _, in := range info {
		slog.Debug(fmt.Sprintf("Service %s", in["Services"]))
		if d, ok := tmp[in["Services"]]; ok {
			d = append(d, in)
			tmp[in["Services"]] = d
			// slog.Debugf("Service ADD %v", tmp)
		} else {
			t := []map[string]string{in}
			tmp[in["Services"]] = t
			// slog.Debugf("New Service %v", tmp)
		}
	}

	// 去重
	var sortsKey []string
	for key := range tmp {
		sortsKey = append(sortsKey, key)
	}

	// 排序
	sort.Strings(sortsKey)

	result := ""
	for _, key := range sortsKey {
		value := tmp[key]
		// result := fmt.Sprintf("<li class=\"list-group-item list-group-item-info\">%s</li>", strings.ToUpper(beego.AppConfig.String("appname")))
		result += fmt.Sprintf("<li class=\"list-group-item list-group-item-info\">Project %s</li>", strings.ToUpper(key))
		for _, data := range value {
			// result += fmt.Sprintf("<li class=\"list-group-item list-group-item-info\">%s</li>", data["Struct"])
			if data["Struct"] != "User" && data["Struct"] != "Group" && data["Struct"] != "Userauth" && data["Struct"] != "History" {
				result += fmt.Sprintf("<li class=\"list-group-item\"><a class=\"badge\" href=\"/admin/table?name=%s\"><span class=\"glyphicon glyphicon-edit\"></span>Change</a><a class=\"badge\" href=\"/admin/add?name=%s\"><span class=\"glyphicon glyphicon-plus\"></span>Add</a><a href=\"/admin/table?name=%s\" target=\"_self\">%s</a></li>", data["Struct"], data["Struct"], data["Struct"], data["Struct"])
			}
			// for _, x := range strings.Split(data["Name"], " ") {
			// 	result += fmt.Sprintf("<li class=\"list-group-item\"><a class=\"badge\" href=\"#\">Change</a><a class=\"badge\" href=\"#\">Add</a>%s</li>", x)
			// }
		}
	}

	result += "<div class=\"row\">&nbsp;</div>"
	return result
}

func AdminColumns(data map[string]string) string {
	//get columns
	// col := strings.TrimSpace(data["Names"])
	col := strings.TrimSpace(data["List"]) + " 操作"
	result, err := DirectJson(strings.Split(col, " ")...)
	if err != nil {
		return err.Error()
	}
	return result
}

type Columns struct {
	Field    string //state
	Title    string
	Checkbox bool
	Rowspan  int
	Colspan  int
	Align    string
	Valign   string
	Sortable bool
	editable bool
}

func NewColumns(filed string, checkbox bool) *Columns {
	return &Columns{
		Field:    filed,
		Title:    strings.ToUpper(filed),
		Checkbox: checkbox,
		Align:    "center",
		Valign:   "middle",
	}
}

func MutilColumms(data ...string) []Columns {
	tmp := []Columns{Columns{Checkbox: true}}
	for _, x := range data {
		tmp = append(tmp, *NewColumns(x, false))
	}
	return tmp
}

func DirectJson(data ...string) (string, error) {
	tmp := MutilColumms(data...)
	jsons, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}
	return strings.ToLower(string(jsons)), nil
}

/*
[map[Tag:name:"ip" search:"false" Name: Id Vpn Name Ip Id_search:"true" Name_search:"false" Ip_name:"ip" Ip_search:"false" Struct:Vpn Type:string Id_name:"id" Vpn_name:"vpn" Vpn_search:"true" Name_name:"name"] map[Name_search:"true" Struct:Machine Type:time.Time Name: Id Sn Mac Ip Name Create Update Sn_name:"sn" Sn_search:"true" Mac_search:"true" Tag:xorm:"updated" Mac_xorm:"mac" Mac_name:"mac" Update_xorm:"updated" Ip_xorm:"ip" Name_name:"name" Create_xorm:"created" Id_name:"id" Id_search:"true" Sn_xorm:"sn" Ip_name:"ip" Ip_search:"true" Name_xorm:"name"] map[Type:string Name: Id Cdn_name Type Type_name:"type" Struct:Cdn Tag:name:"type" search:"false" Cdn_name_name:"cdn_name" Cdn_name_search:"true" Type_search:"false" Id_name:"id" Id_search:"true"]]
*/
func FormColumns(data map[string]string) string {
	result := ""
	text := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right" for="form-field-1"> $LABELS </label>
	<div class="col-sm-10">
		<input name="$NAME" placeholder="$LABELS" class="col-xs-12 col-sm-12" type="$TYPE">
	</div>
</div>`
	textarea := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right" for="form-field-1"> $LABELS </label>
	<div class="col-sm-10">
		<textarea name="$NAME" class="col-xs-12 col-sm-12" rows="10"></textarea>
	</div>
</div>`
	radio := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right"> $LABELS </label>
	<div class="col-sm-10">
		$CONTENT
	</div>
</div>`
	selected := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right" for="$NAME"> $LABELS </label>
	<div class="col-sm-10">
		<select name="$NAME">
			$CONTENT
		</select>
	</div>
</div>`
	multiselect := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-top" for="duallist"> $LABELS </label>
	<div class="col-sm-10">
		<select multiple="multiple" size="10" name="$NAME" id="duallist" class="col-xs-10 col-sm-10" >
			$CONTENT
		</select>

		<div class="hr hr-16 hr-dotted"></div>
	</div>
</div><script>
jQuery(function($){
		var demo1 = $('#duallist').bootstrapDualListbox({infoTextFiltered: '<span class="label label-purple label-lg">Filtered</span>'});
		var container1 = demo1.bootstrapDualListbox('getContainer');
		container1.find('.btn').addClass('btn-white btn-info btn-bold');

		//in ajax mode, remove remaining elements before leaving page
		$(document).one('ajaxloadstart.page', function(e) {
			$('[class*=select2]').remove();
			$('#duallist').bootstrapDualListbox('destroy');
			$('.rating').raty('destroy');
			$('.multiselect').multiselect('destroy');
		});
});
</script>`
	timed := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-top" for="date-timepicker1"> $LABELS </label>

	<div class="col-sm-10">
		<input id="date-timepicker1" name="$NAME" type="text" class="col-xs-2 col-sm-2" />
	</div>
</div><script>
jQuery(function($){
		$('#date-timepicker1').datetimepicker({
			//format: 'YYYY/MM/DD h:mm:ss',//use this option to display seconds
			icons: {
			time: 'fa fa-clock-o',
			date: 'fa fa-calendar',
			up: 'fa fa-chevron-up',
			down: 'fa fa-chevron-down',
			previous: 'fa fa-chevron-left',
			next: 'fa fa-chevron-right',
			today: 'fa fa-arrows ',
			clear: 'fa fa-trash',
			close: 'fa fa-times'
			}
	   });
});
</script>`
	if data["Col"] == "" {
		return "Col is none"
	}

	for _, info := range strings.Split(strings.TrimSpace(data["Col"]), " ") {
		tmp := strings.Split(info, ":")
		switch tmp[1] {
		case "string":
			result += strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "text", -1)
		case "int", "int16", "int64":
			result += strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "number", -1)
		case "textarea":
			result += strings.Replace(strings.Replace(textarea, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1)
		case "radio":
			tmp_radio := ""
			for n, r := range strings.Split(tmp[3], ",") {
				tmp_ro := strings.Split(r, "|")
				if n == 0 {
					tmp_radio += fmt.Sprintf("%s <input type=\"radio\" checked=\"checked\" name=\"%s\" value=\"%s\">", tmp_ro[0], tmp[2], tmp_ro[1])
				} else {
					tmp_radio += fmt.Sprintf("%s <input type=\"radio\" name=\"%s\" value=\"%s\">", tmp_ro[0], tmp[2], tmp_ro[1])
				}
			}
			result += strings.Replace(strings.Replace(radio, "$LABELS", tmp[0], -1), "$CONTENT", tmp_radio, -1)
		case "select":
			tmp_select := ""
			for _, s := range strings.Split(tmp[3], ",") {
				tmp_se := strings.Split(s, "|")
				tmp_select += fmt.Sprintf("<option value=\"%s\">%s</option>", tmp_se[1], tmp_se[0])
			}
			result += strings.Replace(strings.Replace(strings.Replace(selected, "$CONTENT", tmp_select, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1)
		case "multiselect":
			tmp_multiselect := ""
			for _, s := range strings.Split(tmp[3], ",") {
				tmp_se := strings.Split(s, "|")
				tmp_multiselect += fmt.Sprintf("<option value=\"%s\">%s</option>", tmp_se[1], tmp_se[0])
			}
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(multiselect, "$CONTENT", tmp_multiselect, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1), "duallist", fmt.Sprintf("%d", time.Now().Nanosecond()), -1)
		case "file":
			result += strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "file", -1)
		case "time", "time.Time":
			result += strings.Replace(strings.Replace(strings.Replace(timed, "-timepicker1", fmt.Sprintf("%d", time.Now().Nanosecond()), -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1)
		case "o2o":
			tmp_o2o_string := ""
			tmp_o2o := strings.Split(tmp[3], "|")
			// slog.Error("tmp_o2o", tmp[3])
			sql := fmt.Sprintf("select * from %s", tmp_o2o[0])
			// slog.Error(sql)
			resultSql, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				slog.Error(err.Error())
			}
			showCol := strings.Split(fmt.Sprintf("id,%s", tmp_o2o[1]), ",")
			for _, x := range resultSql {
				t_s := "<option value=\"%id\">%value</option>"
				t_v := ""
				for _, value := range showCol {
					t_v += fmt.Sprintf("%s ", string(x[value]))
					// slog.Error(value, t_v, string(x[value]), x)
				}
				tmp_o2o_string += strings.Replace(strings.Replace(t_s, "%id", string(x[showCol[0]]), -1), "%value", t_v, -1)
				// slog.Error(tmp_o2o_string)
			}
			result += strings.Replace(strings.Replace(strings.Replace(selected, "$CONTENT", tmp_o2o_string, -1), "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1)
		case "o2m":
			tmp_o2o_string := ""
			tmp_o2o := strings.Split(tmp[3], "|")
			// slog.Error("tmp_o2o", tmp[3])
			sql := fmt.Sprintf("select * from %s", tmp_o2o[0])
			// slog.Error(sql)
			resultSql, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				slog.Error(err.Error())
			}
			showCol := strings.Split(fmt.Sprintf("id,%s", tmp_o2o[1]), ",")
			for _, x := range resultSql {
				t_s := "<option value=\"%id\">%value</option>"
				t_v := ""
				for _, value := range showCol {
					t_v += fmt.Sprintf("%s ", string(x[value]))
					// slog.Error(value, t_v, string(x[value]), x)
				}
				tmp_o2o_string += strings.Replace(strings.Replace(t_s, "%id", string(x[showCol[0]]), -1), "%value", t_v, -1)
				// slog.Error(tmp_o2o_string)
			}
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(multiselect, "$CONTENT", tmp_o2o_string, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1), "duallist", fmt.Sprintf("%d", time.Now().Nanosecond()), -1)
		}
	}

	return result
}

// 根据查询到的结果映射为对应的form表单及数据
func EditFormColumns(data map[string]string, dbinfo map[string][]byte) string {
	result := fmt.Sprintf("<input name=\"id\" value=\"%s\" style=\"display:none\" />", dbinfo["id"])
	text := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right" for="form-field-1"> $LABELS </label>
	<div class="col-sm-10">
		<input name="$NAME" placeholder="$LABELS" class="col-xs-12 col-sm-12" type="$TYPE" value="$VALUE">
	</div>
</div>`
	textarea := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right" for="form-field-1"> $LABELS </label>
	<div class="col-sm-10">
		<textarea name="$NAME" class="col-xs-12 col-sm-12" rows="10">$VALUE</textarea>
	</div>
</div>`
	radio := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right"> $LABELS </label>
	<div class="col-sm-10">
		$CONTENT
	</div>
</div>`
	selected := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right" for="$NAME"> $LABELS </label>
	<div class="col-sm-10">
		<select name="$NAME">
			$CONTENT
		</select>
	</div>
</div>`
	multiselect := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-top" for="duallist"> $LABELS </label>
	<div class="col-sm-10">
		<select multiple="multiple" size="10" name="$NAME" id="duallist" class="col-xs-10 col-sm-10" >
			$CONTENT
		</select>

		<div class="hr hr-16 hr-dotted"></div>
	</div>
</div><script>
jQuery(function($){
		var demo1 = $('#duallist').bootstrapDualListbox({infoTextFiltered: '<span class="label label-purple label-lg">Filtered</span>'});
		var container1 = demo1.bootstrapDualListbox('getContainer');
		container1.find('.btn').addClass('btn-white btn-info btn-bold');

		//in ajax mode, remove remaining elements before leaving page
		$(document).one('ajaxloadstart.page', function(e) {
			$('[class*=select2]').remove();
			$('#duallist').bootstrapDualListbox('destroy');
			$('.rating').raty('destroy');
			$('.multiselect').multiselect('destroy');
		});
});
</script>`
	timed := `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-top" for="date-timepicker1"> $LABELS </label>

	<div class="col-sm-10">
		<input id="date-timepicker1" name="$NAME" type="text" class="col-xs-2 col-sm-2" value="$VALUE" />
	</div>
</div><script>
jQuery(function($){
		$('#date-timepicker1').datetimepicker({
			//format: 'YYYY/MM/DD h:mm:ss',//use this option to display seconds
			icons: {
			time: 'fa fa-clock-o',
			date: 'fa fa-calendar',
			up: 'fa fa-chevron-up',
			down: 'fa fa-chevron-down',
			previous: 'fa fa-chevron-left',
			next: 'fa fa-chevron-right',
			today: 'fa fa-arrows ',
			clear: 'fa fa-trash',
			close: 'fa fa-times'
			}
	   });
});
</script>`
	for _, info := range strings.Split(strings.TrimSpace(data["Col"]), " ") {
		// slog.Error(info)
		tmp := strings.Split(info, ":")
		switch tmp[1] {
		case "password":
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "password", -1), "$VALUE", string(dbinfo[tmp[2]]), -1)
		case "string":
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "text", -1), "$VALUE", string(dbinfo[tmp[2]]), -1)
		case "int", "int16", "int64":
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(text, "$VALUE", string(dbinfo[tmp[2]]), -1), "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "number", -1)
		case "textarea":
			result += strings.Replace(strings.Replace(strings.Replace(textarea, "$VALUE", string(dbinfo[tmp[2]]), -1), "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1)
		case "radio":
			tmp_radio := ""
			for _, r := range strings.Split(tmp[3], ",") {
				tmp_ro := strings.Split(r, "|")
				if tmp_ro[1] == string(dbinfo[tmp[2]]) {
					tmp_radio += fmt.Sprintf("%s <input type=\"radio\" checked=\"checked\" name=\"%s\" value=\"%s\">", tmp_ro[0], tmp[2], tmp_ro[1])
				} else {
					tmp_radio += fmt.Sprintf("%s <input type=\"radio\" name=\"%s\" value=\"%s\">", tmp_ro[0], tmp[2], tmp_ro[1])
				}
			}
			result += strings.Replace(strings.Replace(radio, "$LABELS", tmp[0], -1), "$CONTENT", tmp_radio, -1)
		case "select":
			tmp_select := ""
			for _, s := range strings.Split(tmp[3], ",") {
				tmp_se := strings.Split(s, "|")
				if tmp_se[1] == string(dbinfo[tmp[2]]) {
					tmp_select += fmt.Sprintf("<option value=\"%s\" selected=\"selected\">%s</option>", tmp_se[1], tmp_se[0])
				} else {
					tmp_select += fmt.Sprintf("<option value=\"%s\">%s</option>", tmp_se[1], tmp_se[0])
				}
			}
			result += strings.Replace(strings.Replace(strings.Replace(selected, "$CONTENT", tmp_select, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1)
		case "multiselect":
			tmp_multiselect := ""
			for _, s := range strings.Split(tmp[3], ",") {
				tmp_se := strings.Split(s, "|")
				if strings.Contains(string(dbinfo[tmp[2]]), tmp_se[1]) {
					tmp_multiselect += fmt.Sprintf("<option value=\"%s\" selected=\"selected\">%s</option>", tmp_se[1], tmp_se[0])
				} else {
					tmp_multiselect += fmt.Sprintf("<option value=\"%s\">%s</option>", tmp_se[1], tmp_se[0])
				}
			}
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(multiselect, "$CONTENT", tmp_multiselect, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1), "duallist", fmt.Sprintf("%d", time.Now().Nanosecond()), -1)
		case "file":
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(text, "$VALUE", string(dbinfo[tmp[2]]), -1), "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "file", -1)
		case "time", "time.Time":
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(timed, "$VALUE", string(dbinfo[tmp[2]]), -1), "-timepicker1", fmt.Sprintf("%d", time.Now().Nanosecond()), -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1)
		case "o2o":
			tmp_o2o_string := ""
			tmp_o2o := strings.Split(tmp[3], "|")
			// slog.Error("tmp_o2o", tmp[3])
			sql := fmt.Sprintf("select * from %s", tmp_o2o[0])
			// slog.Error(sql)
			resultSql, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				slog.Error(err.Error())
			}
			showCol := strings.Split(fmt.Sprintf("id,%s", tmp_o2o[1]), ",")
			for _, x := range resultSql {
				var t_s string
				if strings.Contains(string(dbinfo[tmp[2]]), string(x[showCol[0]])) {
					t_s = "<option value=\"%id\" selected=\"selected\">%value</option>"
				} else {
					t_s = "<option value=\"%id\">%value</option>"
				}
				t_v := ""
				for _, value := range showCol {
					t_v += fmt.Sprintf("%s ", string(x[value]))
					// slog.Error(value, t_v, string(x[value]), x)
				}
				tmp_o2o_string += strings.Replace(strings.Replace(t_s, "%id", string(x[showCol[0]]), -1), "%value", t_v, -1)
				// slog.Error(tmp_o2o_string)
			}
			result += strings.Replace(strings.Replace(strings.Replace(selected, "$CONTENT", tmp_o2o_string, -1), "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1)
		case "o2m":
			tmp_o2o_string := ""
			tmp_o2o := strings.Split(tmp[3], "|")
			// slog.Error("tmp_o2o", tmp[3])
			sql := fmt.Sprintf("select * from %s", tmp_o2o[0])
			// slog.Error(sql)
			resultSql, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				slog.Error(err.Error())
			}
			showCol := strings.Split(fmt.Sprintf("id,%s", tmp_o2o[1]), ",")
			for _, x := range resultSql {
				var t_s string
				if strings.Contains(string(dbinfo[tmp[2]]), string(x[showCol[0]])) {
					t_s = "<option value=\"%id\" selected=\"selected\">%value</option>"
				} else {
					t_s = "<option value=\"%id\">%value</option>"
				}
				t_v := ""
				for _, value := range showCol {
					t_v += fmt.Sprintf("%s ", string(x[value]))
					// slog.Error(value, t_v, string(x[value]), x)
				}
				tmp_o2o_string += strings.Replace(strings.Replace(t_s, "%id", string(x[showCol[0]]), -1), "%value", t_v, -1)
				// slog.Error(tmp_o2o_string)
			}
			result += strings.Replace(strings.Replace(strings.Replace(strings.Replace(multiselect, "$CONTENT", tmp_o2o_string, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1), "duallist", fmt.Sprintf("%d", time.Now().Nanosecond()), -1)
		}
	}

	return result
}
