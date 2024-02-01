package addform

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	orm "github.com/lflxp/djangolang/consted"
	"github.com/lflxp/djangolang/utils/orm/sqlite"
)

// 添加表单字段
var addColumns = map[orm.FieldType]string{
	orm.Text: `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right" for="form-field-1"> $LABELS </label>
	<div class="col-sm-10">
		<input name="$NAME" placeholder="$LABELS" class="col-xs-12 col-sm-12" type="$TYPE">
	</div>
</div>`,
	orm.Textarea: `<div class="form-group">
<label class="col-sm-2 control-label no-padding-right" for="form-field-1"> $LABELS </label>
<div class="col-sm-10">
	<textarea name="$NAME" class="col-xs-12 col-sm-12" rows="10"></textarea>
</div>
</div>`,
	orm.Radio: `<div class="form-group">
	<label class="col-sm-2 control-label no-padding-right"> $LABELS </label>
	<div class="col-sm-10">
		$CONTENT
	</div>
</div>`,
	orm.Select: `<div class="form-group">
<label class="col-sm-2 control-label no-padding-right" for="$NAME"> $LABELS </label>
<div class="col-sm-10">
	<select name="$NAME">
		$CONTENT
	</select>
</div>
</div>`,
	orm.MultiSelect: `<div class="form-group">
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
</script>`,
	orm.Time: `<div class="form-group">
<label class="col-sm-2 control-label no-padding-top" for="date-timepicker1"> $LABELS </label>

<div class="col-sm-10">
	<input id="date-timepicker1" name="$NAME" type="text" class="col-xs-2 col-sm-2" />
</div>
</div><script>
jQuery(function($){
	$('#date-timepicker1').datetimepicker({
		format: 'YYYY-MM-DD hh:mm:ss',//use this option to display seconds
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
</script>`,
}

var addFuncMap = map[orm.FieldType]func(string, []string) string{
	orm.Password: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "password", -1)
	},
	orm.StringField: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "text", -1)
	},
	orm.IntField: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "number", -1)
	},
	orm.Int16Field: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "number", -1)
	},
	orm.Int64Field: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1), "$TYPE", "number", -1)
	},
	orm.Textarea: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(text, "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1)
	},
	orm.Radio: func(text string, tmp []string) string {
		tmp_radio := ""
		for n, r := range strings.Split(tmp[3], ",") {
			tmp_ro := strings.Split(r, "|")
			if n == 0 {
				tmp_radio += fmt.Sprintf("%s <input type=\"radio\" checked=\"checked\" name=\"%s\" value=\"%s\">", tmp_ro[0], tmp[2], tmp_ro[1])
			} else {
				tmp_radio += fmt.Sprintf("%s <input type=\"radio\" name=\"%s\" value=\"%s\">", tmp_ro[0], tmp[2], tmp_ro[1])
			}
		}
		return strings.Replace(strings.Replace(text, "$LABELS", tmp[0], -1), "$CONTENT", tmp_radio, -1)
	},
	orm.Select: func(text string, tmp []string) string {
		tmp_select := ""
		for _, s := range strings.Split(tmp[3], ",") {
			tmp_se := strings.Split(s, "|")
			tmp_select += fmt.Sprintf("<option value=\"%s\">%s</option>", tmp_se[1], tmp_se[0])
		}
		return strings.Replace(strings.Replace(strings.Replace(text, "$CONTENT", tmp_select, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1)
	},
	orm.MultiSelect: func(text string, tmp []string) string {
		tmp_multiselect := ""
		for _, s := range strings.Split(tmp[3], ",") {
			tmp_se := strings.Split(s, "|")
			tmp_multiselect += fmt.Sprintf("<option value=\"%s\">%s</option>", tmp_se[1], tmp_se[0])
		}
		return strings.Replace(strings.Replace(strings.Replace(strings.Replace(text, "$CONTENT", tmp_multiselect, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1), "duallist", fmt.Sprintf("%d", time.Now().Nanosecond()), -1)
	},
	orm.File: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(strings.Replace(text, "$NAME", tmp[2]+"\" id='"+tmp[2]+"'", -1), "$LABELS", tmp[0], -1), "$TYPE", "file", -1)
	},
	orm.Time: func(text string, tmp []string) string {
		return strings.Replace(strings.Replace(strings.Replace(text, "-timepicker1", fmt.Sprintf("%d", time.Now().Nanosecond()), -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1)
	},
	orm.OneToOne: func(text string, tmp []string) string {
		tmp_o2o_string := ""
		tmp_o2o := strings.Split(tmp[3], "|")
		// log.Error("tmp_o2o", tmp[3])
		sql := fmt.Sprintf("select * from %s", tmp_o2o[0])
		// log.Error(sql)
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
				// log.Error(value, t_v, string(x[value]), x)
			}
			tmp_o2o_string += strings.Replace(strings.Replace(t_s, "%id", string(x[showCol[0]]), -1), "%value", t_v, -1)
			// log.Error(tmp_o2o_string)
		}
		return strings.Replace(strings.Replace(strings.Replace(text, "$CONTENT", tmp_o2o_string, -1), "$NAME", tmp[2], -1), "$LABELS", tmp[0], -1)
	},
	orm.OneToMany: func(text string, tmp []string) string {
		tmp_o2o_string := ""
		tmp_o2o := strings.Split(tmp[3], "|")
		// log.Error("tmp_o2o", tmp[3])
		sql := fmt.Sprintf("select * from %s", tmp_o2o[0])
		// log.Error(sql)
		resultSql, err := sqlite.NewOrm().Query(sql)
		if err != nil {
			slog.Error(err.Error())
		}
		showCol := strings.Split(tmp_o2o[1], ",")
		slog.Debug("resultSql %s\n showCol %s \n tmp_o2o %s", resultSql, showCol, tmp_o2o)
		for _, x := range resultSql {
			t_s := "<option value=\"%id\">%value</option>"
			t_v := ""
			for _, value := range showCol {
				t_v += fmt.Sprintf("%s ", string(x[value]))
				// log.Error(value, t_v, string(x[value]), x)
			}
			tmp_o2o_string += strings.Replace(strings.Replace(t_s, "%id", string(x[showCol[0]]), -1), "%value", t_v, -1)
			// log.Error(tmp_o2o_string)
		}
		return strings.Replace(strings.Replace(strings.Replace(strings.Replace(text, "$CONTENT", tmp_o2o_string, -1), "$LABELS", tmp[0], -1), "$NAME", tmp[2], -1), "duallist", fmt.Sprintf("%d", time.Now().Nanosecond()), -1)
	},
}
