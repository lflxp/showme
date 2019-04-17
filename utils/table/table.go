package table

import (
	"errors"
	"fmt"
	"io"

	"github.com/lflxp/showme/utils"
)

type Col struct {
	data      string
	color     string
	TextAlign Align
}

func NewCol() *Col {
	return &Col{}
}

func (this *Col) SetColor(color string) *Col {
	this.color = color
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

func (this *Table) CalColumnWidths(cs []int) error {
	if this.width == 0 {
		return errors.New("width is nil or 0")
	}
	if len(this.ColumnWidths) == 0 {
		colCount := len(this.Rows[0])
		colWidth := this.width / colCount
		for i := 0; i < colCount; i++ {
			this.ColumnWidths = append(this.ColumnWidths, colWidth)
		}
	} else {
		this.ColumnWidths = cs
	}

	return nil
}

// 顺序添加列名
func (this *Table) AddCol(title string) *Col {
	t := &Col{data: title}
	this.Rows = append(this.Rows, []*Col{t})
	return t
}

// 顺序增加第row行的数据
func (this *Table) AddRow(row int, data *Col) {
	tmp := this.Rows[row]
	tmp = append(tmp, data)
	this.Rows[row] = tmp
}

// 打印内容
func (this *Table) Fprint(w io.Writer) {
	for _, r := range this.Rows {
		for n, col := range r {
			fmt.Fprintf(w, col.data)
			if this.ShowHeader == true && n == 0 {
				var s string
				switch col.TextAlign {
				case TextLeft:
					s = AlignLeft(utils.Colorize(col.data, col.color, "dgreen", false, true), this.ColumnWidths[n])
				case TextRight:
					s = AlignRight(utils.Colorize(col.data, col.color, "dgreen", false, true), this.ColumnWidths[n])
				case TextCenter:
					s = AlignCenter(utils.Colorize(col.data, col.color, "dgreen", false, true), this.ColumnWidths[n])
				}
				fmt.Fprintln(w, s)
			} else if n > 0 {
				var s string
				switch col.TextAlign {
				case TextLeft:
					s = AlignLeft(utils.Colorize(col.data, col.color, "dgreen", false, true), this.ColumnWidths[n])
				case TextRight:
					s = AlignRight(utils.Colorize(col.data, col.color, "dgreen", false, true), this.ColumnWidths[n])
				case TextCenter:
					s = AlignCenter(utils.Colorize(col.data, col.color, "dgreen", false, true), this.ColumnWidths[n])
				}
				fmt.Fprintf(w, s)
			}
		}
	}
}

func TableTest(w io.Writer, width int) {
	table := NewTable(width)
	table.ShowHeader = true

	table.AddCol("name").SetColor("red").SetTextAlign(TextCenter)
	table.AddCol("rank").SetColor("blue").SetTextAlign(TextLeft)
	table.AddCol("time").SetColor("green").SetTextAlign(TextRight)

	fmt.Println(table)
	for i := 0; i <= 100; i++ {
		tmp := NewCol()
		tmp.data = fmt.Sprintf("name %d", i)
		tmp.TextAlign = TextCenter
		tmp.color = "blue"
		table.AddRow(0, tmp)
		rank := NewCol()
		rank.data = fmt.Sprintf("rank %d", i)
		rank.TextAlign = TextCenter
		rank.color = "white"
		table.AddRow(1, rank)
		time := NewCol()
		time.data = fmt.Sprintf("time %d", i)
		time.TextAlign = TextCenter
		time.color = "white"
		table.AddRow(2, time)
	}
	table.CalColumnWidths(nil)
	table.Fprint(w)
}
