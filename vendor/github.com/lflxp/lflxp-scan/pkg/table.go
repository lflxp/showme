package pkg

import (
	"errors"
	"fmt"
	"io"
)

type Col struct {
	Data      string
	Color     string
	BgColor   string
	TextAlign Align
}

func NewCol() *Col {
	return &Col{}
}

func (this *Col) SetColor(Color string) *Col {
	this.Color = Color
	return this
}

func (this *Col) SetBgColor(Color string) *Col {
	this.BgColor = Color
	return this
}

func (this *Col) SetTextAlign(align Align) *Col {
	this.TextAlign = align
	return this
}

type Table struct {
	width        int
	Header       []*Col
	Rows         [][]*Col // row 行 col 列
	ColumnWidths []int
	ShowHeader   bool
}

func NewTable(width int) *Table {
	return &Table{
		width: width,
		Rows:  [][]*Col{},
	}
}

func (this *Table) CalColumnWidths() error {
	if this.width == 0 {
		return errors.New("width is nil or 0")
	}
	colCount := len(this.Rows)
	colWidth := this.width / colCount
	for i := 0; i < colCount; i++ {
		this.ColumnWidths = append(this.ColumnWidths, colWidth)
	}

	return nil
}

// 顺序添加列名
func (this *Table) AddCol(title string) *Col {
	t := &Col{Data: title}
	this.Header = append(this.Header, t)
	this.Rows = append(this.Rows, []*Col{})
	return t
}

// 顺序增加第row行的数据
func (this *Table) AddRow(row int, Data *Col) {
	tmp := this.Rows[row]
	tmp = append(tmp, Data)
	this.Rows[row] = tmp
}

// 顺序增加第row行的数据 指定index
func (this *Table) AddRowByIndex(row int, Data *Col) {
	tmp := this.Rows[row]
	tmp[0] = Data
	this.Rows[row] = tmp
}

func (this Table) FprintHeader(w io.Writer) {
	width := this.width / len(this.Header)
	for i := 0; i < len(this.Header); i++ {
		tmp := this.Header[i]
		var s string
		switch tmp.TextAlign {
		case TextLeft:
			s = AlignLeft(tmp.Data, width)
		case TextRight:
			s = AlignRight(tmp.Data, width)
		case TextCenter:
			s = AlignCenter(tmp.Data, width)
		}
		fmt.Fprintf(w, Colorize(s, tmp.Color, tmp.BgColor, true, true))
	}
	fmt.Fprintln(w, "")
}

// 打印内容
// 正序其实打印数据只有一个
func (this *Table) Fprint(w io.Writer) {
	// 字段个数
	colNum := len(this.Rows)
	count := len(this.Rows[0])
	// 行列转换
	// 倒序打印
	for cc := 0; cc < count; cc++ {
		for col := 0; col < colNum; col++ {
			tmp := this.Rows[col][cc]
			var s string
			switch tmp.TextAlign {
			case TextLeft:
				s = AlignLeft(tmp.Data, this.ColumnWidths[col])
			case TextRight:
				s = AlignRight(tmp.Data, this.ColumnWidths[col])
			case TextCenter:
				s = AlignCenter(tmp.Data, this.ColumnWidths[col])
			}
			fmt.Fprintf(w, Colorize(s, tmp.Color, tmp.BgColor, false, true))
		}
		fmt.Fprintln(w, "")
	}
	// for rlength, r := range this.Rows {
	// 	for _, col := range r {
	// 		// fmt.Fprintln(w, fmt.Sprintf("%s %d %d", col.Data, this.ColumnWidths[n], n))
	// 		var s string
	// 		switch col.TextAlign {
	// 		case TextLeft:
	// 			s = AlignLeft(col.Data, this.ColumnWidths[rlength])
	// 		case TextRight:
	// 			s = AlignRight(col.Data, this.ColumnWidths[rlength])
	// 		case TextCenter:
	// 			s = AlignCenter(col.Data, this.ColumnWidths[rlength])
	// 		}
	// 		fmt.Fprintf(w, s)
	// 	}
	// 	fmt.Fprintln(w, "\n")
	// }
}

// 打印内容
// 倒序必须是全量数据
func (this *Table) FprintOrderDesc(w io.Writer) {
	// 字段个数
	colNum := len(this.Rows)
	count := len(this.Rows[0])
	// 行列转换
	// 倒序打印
	for cc := count - 1; cc >= 0; cc-- {
		for col := 0; col < colNum; col++ {
			tmp := this.Rows[col][cc]
			var s string
			switch tmp.TextAlign {
			case TextLeft:
				s = AlignLeft(tmp.Data, this.ColumnWidths[col])
			case TextRight:
				s = AlignRight(tmp.Data, this.ColumnWidths[col])
			case TextCenter:
				s = AlignCenter(tmp.Data, this.ColumnWidths[col])
			}
			fmt.Fprintf(w, Colorize(s, tmp.Color, tmp.BgColor, false, true))
		}
		fmt.Fprintln(w, "")
	}
	// for rlength, r := range this.Rows {
	// 	for _, col := range r {
	// 		// fmt.Fprintln(w, fmt.Sprintf("%s %d %d", col.Data, this.ColumnWidths[n], n))
	// 		var s string
	// 		switch col.TextAlign {
	// 		case TextLeft:
	// 			s = AlignLeft(col.Data, this.ColumnWidths[rlength])
	// 		case TextRight:
	// 			s = AlignRight(col.Data, this.ColumnWidths[rlength])
	// 		case TextCenter:
	// 			s = AlignCenter(col.Data, this.ColumnWidths[rlength])
	// 		}
	// 		fmt.Fprintf(w, s)
	// 	}
	// 	fmt.Fprintln(w, "\n")
	// }
}

func TableTest(w io.Writer, width int) {
	table := NewTable(width)
	table.ShowHeader = true

	table.AddCol("name").SetColor("red").SetTextAlign(TextCenter).SetBgColor("dgreen")
	table.AddCol("rank").SetColor("blue").SetTextAlign(TextLeft).SetBgColor("dgreen")
	table.AddCol("time").SetColor("green").SetTextAlign(TextRight).SetBgColor("dgreen")

	// fmt.Println(table)
	for i := 0; i <= 100; i++ {
		tmp := NewCol()
		tmp.Data = fmt.Sprintf("name %d", i)
		tmp.TextAlign = TextCenter
		tmp.Color = "blue"
		table.AddRow(0, tmp)
		rank := NewCol()
		rank.Data = fmt.Sprintf("rank %d", i)
		rank.TextAlign = TextCenter
		rank.Color = "white"
		table.AddRow(1, rank)
		time := NewCol()
		time.Data = fmt.Sprintf("time %d", i)
		time.TextAlign = TextCenter
		time.Color = "yellow"
		table.AddRow(2, time)
	}
	table.CalColumnWidths()
	table.Fprint(w)
}
