
//----------------------------------------
// The code is automatically generated by the GenlibLcl tool.
// Copyright © ying32. All Rights Reserved.
// 
// Licensed under Apache License 2.0
//
//----------------------------------------


package vcl


import (
    . "github.com/ying32/govcl/vcl/api"
    . "github.com/ying32/govcl/vcl/types"
    "unsafe"
)

type TDrawGrid struct {
    IWinControl
    instance uintptr
    // 特殊情况下使用，主要应对Go的GC问题，与LCL没有太多关系。
    ptr unsafe.Pointer
}

// 创建一个新的对象。
// 
// Create a new object.
func NewDrawGrid(owner IComponent) *TDrawGrid {
    d := new(TDrawGrid)
    d.instance = DrawGrid_Create(CheckPtr(owner))
    d.ptr = unsafe.Pointer(d.instance)
    return d
}

// 动态转换一个已存在的对象实例。
// 
// Dynamically convert an existing object instance.
func AsDrawGrid(obj interface{}) *TDrawGrid {
    instance, ptr := getInstance(obj)
    if instance == 0 { return nil }
    return &TDrawGrid{instance: instance, ptr: ptr}
}

// -------------------------- Deprecated begin --------------------------
// 新建一个对象来自已经存在的对象实例指针。
// 
// Create a new object from an existing object instance pointer.
// Deprecated: use AsDrawGrid.
func DrawGridFromInst(inst uintptr) *TDrawGrid {
    return AsDrawGrid(inst)
}

// 新建一个对象来自已经存在的对象实例。
// 
// Create a new object from an existing object instance.
// Deprecated: use AsDrawGrid.
func DrawGridFromObj(obj IObject) *TDrawGrid {
    return AsDrawGrid(obj)
}

// 新建一个对象来自不安全的地址。注意：使用此函数可能造成一些不明情况，慎用。
// 
// Create a new object from an unsecured address. Note: Using this function may cause some unclear situations and be used with caution..
// Deprecated: use AsDrawGrid.
func DrawGridFromUnsafePointer(ptr unsafe.Pointer) *TDrawGrid {
    return AsDrawGrid(ptr)
}

// -------------------------- Deprecated end --------------------------
// 释放对象。
// 
// Free object.
func (d *TDrawGrid) Free() {
    if d.instance != 0 {
        DrawGrid_Free(d.instance)
        d.instance, d.ptr = 0, nullptr
    }
}

// 返回对象实例指针。
// 
// Return object instance pointer.
func (d *TDrawGrid) Instance() uintptr {
    return d.instance
}

// 获取一个不安全的地址。
// 
// Get an unsafe address.
func (d *TDrawGrid) UnsafeAddr() unsafe.Pointer {
    return d.ptr
}

// 检测地址是否为空。
// 
// Check if the address is empty.
func (d *TDrawGrid) IsValid() bool {
    return d.instance != 0
}

// 检测当前对象是否继承自目标对象。
// 
// Checks whether the current object is inherited from the target object.
func (d *TDrawGrid) Is() TIs {
    return TIs(d.instance)
}

// 动态转换当前对象为目标对象。
// 
// Dynamically convert the current object to the target object.
//func (d *TDrawGrid) As() TAs {
//    return TAs(d.instance)
//}

// 获取类信息指针。
// 
// Get class information pointer.
func TDrawGridClass() TClass {
    return DrawGrid_StaticClassType()
}

func (d *TDrawGrid) CellRect(ACol int32, ARow int32) TRect {
    return DrawGrid_CellRect(d.instance, ACol , ARow)
}

func (d *TDrawGrid) MouseToCell(X int32, Y int32, ACol *int32, ARow *int32) {
    DrawGrid_MouseToCell(d.instance, X , Y , ACol , ARow)
}

func (d *TDrawGrid) MouseCoord(X int32, Y int32) TGridCoord {
    return DrawGrid_MouseCoord(d.instance, X , Y)
}

// 是否可以获得焦点。
func (d *TDrawGrid) CanFocus() bool {
    return DrawGrid_CanFocus(d.instance)
}

// 返回是否包含指定控件。
//
// it's contain a specified control.
func (d *TDrawGrid) ContainsControl(Control IControl) bool {
    return DrawGrid_ContainsControl(d.instance, CheckPtr(Control))
}

// 返回指定坐标及相关属性位置控件。
//
// Returns the specified coordinate and the relevant attribute position control..
func (d *TDrawGrid) ControlAtPos(Pos TPoint, AllowDisabled bool, AllowWinControls bool, AllLevels bool) *TControl {
    return AsControl(DrawGrid_ControlAtPos(d.instance, Pos , AllowDisabled , AllowWinControls , AllLevels))
}

// 禁用控件的对齐。
//
// Disable control alignment.
func (d *TDrawGrid) DisableAlign() {
    DrawGrid_DisableAlign(d.instance)
}

// 启用控件对齐。
//
// Enabled control alignment.
func (d *TDrawGrid) EnableAlign() {
    DrawGrid_EnableAlign(d.instance)
}

// 查找子控件。
//
// Find sub controls.
func (d *TDrawGrid) FindChildControl(ControlName string) *TControl {
    return AsControl(DrawGrid_FindChildControl(d.instance, ControlName))
}

func (d *TDrawGrid) FlipChildren(AllLevels bool) {
    DrawGrid_FlipChildren(d.instance, AllLevels)
}

// 返回是否获取焦点。
//
// Return to get focus.
func (d *TDrawGrid) Focused() bool {
    return DrawGrid_Focused(d.instance)
}

// 句柄是否已经分配。
//
// Is the handle already allocated.
func (d *TDrawGrid) HandleAllocated() bool {
    return DrawGrid_HandleAllocated(d.instance)
}

// 插入一个控件。
//
// Insert a control.
func (d *TDrawGrid) InsertControl(AControl IControl) {
    DrawGrid_InsertControl(d.instance, CheckPtr(AControl))
}

// 要求重绘。
//
// Redraw.
func (d *TDrawGrid) Invalidate() {
    DrawGrid_Invalidate(d.instance)
}

// 绘画至指定DC。
//
// Painting to the specified DC.
func (d *TDrawGrid) PaintTo(DC HDC, X int32, Y int32) {
    DrawGrid_PaintTo(d.instance, DC , X , Y)
}

// 移除一个控件。
//
// Remove a control.
func (d *TDrawGrid) RemoveControl(AControl IControl) {
    DrawGrid_RemoveControl(d.instance, CheckPtr(AControl))
}

// 重新对齐。
//
// Realign.
func (d *TDrawGrid) Realign() {
    DrawGrid_Realign(d.instance)
}

// 重绘。
//
// Repaint.
func (d *TDrawGrid) Repaint() {
    DrawGrid_Repaint(d.instance)
}

// 按比例缩放。
//
// Scale by.
func (d *TDrawGrid) ScaleBy(M int32, D int32) {
    DrawGrid_ScaleBy(d.instance, M , D)
}

// 滚动至指定位置。
//
// Scroll by.
func (d *TDrawGrid) ScrollBy(DeltaX int32, DeltaY int32) {
    DrawGrid_ScrollBy(d.instance, DeltaX , DeltaY)
}

// 设置组件边界。
//
// Set component boundaries.
func (d *TDrawGrid) SetBounds(ALeft int32, ATop int32, AWidth int32, AHeight int32) {
    DrawGrid_SetBounds(d.instance, ALeft , ATop , AWidth , AHeight)
}

// 设置控件焦点。
//
// Set control focus.
func (d *TDrawGrid) SetFocus() {
    DrawGrid_SetFocus(d.instance)
}

// 控件更新。
//
// Update.
func (d *TDrawGrid) Update() {
    DrawGrid_Update(d.instance)
}

// 将控件置于最前。
//
// Bring the control to the front.
func (d *TDrawGrid) BringToFront() {
    DrawGrid_BringToFront(d.instance)
}

// 将客户端坐标转为绝对的屏幕坐标。
//
// Convert client coordinates to absolute screen coordinates.
func (d *TDrawGrid) ClientToScreen(Point TPoint) TPoint {
    return DrawGrid_ClientToScreen(d.instance, Point)
}

// 将客户端坐标转为父容器坐标。
//
// Convert client coordinates to parent container coordinates.
func (d *TDrawGrid) ClientToParent(Point TPoint, AParent IWinControl) TPoint {
    return DrawGrid_ClientToParent(d.instance, Point , CheckPtr(AParent))
}

// 是否在拖拽中。
//
// Is it in the middle of dragging.
func (d *TDrawGrid) Dragging() bool {
    return DrawGrid_Dragging(d.instance)
}

// 是否有父容器。
//
// Is there a parent container.
func (d *TDrawGrid) HasParent() bool {
    return DrawGrid_HasParent(d.instance)
}

// 隐藏控件。
//
// Hidden control.
func (d *TDrawGrid) Hide() {
    DrawGrid_Hide(d.instance)
}

// 发送一个消息。
//
// Send a message.
func (d *TDrawGrid) Perform(Msg uint32, WParam uintptr, LParam int) int {
    return DrawGrid_Perform(d.instance, Msg , WParam , LParam)
}

// 刷新控件。
//
// Refresh control.
func (d *TDrawGrid) Refresh() {
    DrawGrid_Refresh(d.instance)
}

// 将屏幕坐标转为客户端坐标。
//
// Convert screen coordinates to client coordinates.
func (d *TDrawGrid) ScreenToClient(Point TPoint) TPoint {
    return DrawGrid_ScreenToClient(d.instance, Point)
}

// 将父容器坐标转为客户端坐标。
//
// Convert parent container coordinates to client coordinates.
func (d *TDrawGrid) ParentToClient(Point TPoint, AParent IWinControl) TPoint {
    return DrawGrid_ParentToClient(d.instance, Point , CheckPtr(AParent))
}

// 控件至于最后面。
//
// The control is placed at the end.
func (d *TDrawGrid) SendToBack() {
    DrawGrid_SendToBack(d.instance)
}

// 显示控件。
//
// Show control.
func (d *TDrawGrid) Show() {
    DrawGrid_Show(d.instance)
}

// 获取控件的字符，如果有。
//
// Get the characters of the control, if any.
func (d *TDrawGrid) GetTextBuf(Buffer *string, BufSize int32) int32 {
    return DrawGrid_GetTextBuf(d.instance, Buffer , BufSize)
}

// 获取控件的字符长，如果有。
//
// Get the character length of the control, if any.
func (d *TDrawGrid) GetTextLen() int32 {
    return DrawGrid_GetTextLen(d.instance)
}

// 设置控件字符，如果有。
//
// Set control characters, if any.
func (d *TDrawGrid) SetTextBuf(Buffer string) {
    DrawGrid_SetTextBuf(d.instance, Buffer)
}

// 查找指定名称的组件。
//
// Find the component with the specified name.
func (d *TDrawGrid) FindComponent(AName string) *TComponent {
    return AsComponent(DrawGrid_FindComponent(d.instance, AName))
}

// 获取类名路径。
//
// Get the class name path.
func (d *TDrawGrid) GetNamePath() string {
    return DrawGrid_GetNamePath(d.instance)
}

// 复制一个对象，如果对象实现了此方法的话。
//
// Copy an object, if the object implements this method.
func (d *TDrawGrid) Assign(Source IObject) {
    DrawGrid_Assign(d.instance, CheckPtr(Source))
}

// 获取类的类型信息。
//
// Get class type information.
func (d *TDrawGrid) ClassType() TClass {
    return DrawGrid_ClassType(d.instance)
}

// 获取当前对象类名称。
//
// Get the current object class name.
func (d *TDrawGrid) ClassName() string {
    return DrawGrid_ClassName(d.instance)
}

// 获取当前对象实例大小。
//
// Get the current object instance size.
func (d *TDrawGrid) InstanceSize() int32 {
    return DrawGrid_InstanceSize(d.instance)
}

// 判断当前类是否继承自指定类。
//
// Determine whether the current class inherits from the specified class.
func (d *TDrawGrid) InheritsFrom(AClass TClass) bool {
    return DrawGrid_InheritsFrom(d.instance, AClass)
}

// 与一个对象进行比较。
//
// Compare with an object.
func (d *TDrawGrid) Equals(Obj IObject) bool {
    return DrawGrid_Equals(d.instance, CheckPtr(Obj))
}

// 获取类的哈希值。
//
// Get the hash value of the class.
func (d *TDrawGrid) GetHashCode() int32 {
    return DrawGrid_GetHashCode(d.instance)
}

// 文本类信息。
//
// Text information.
func (d *TDrawGrid) ToString() string {
    return DrawGrid_ToString(d.instance)
}

func (d *TDrawGrid) AnchorToNeighbour(ASide TAnchorKind, ASpace int32, ASibling IControl) {
    DrawGrid_AnchorToNeighbour(d.instance, ASide , ASpace , CheckPtr(ASibling))
}

func (d *TDrawGrid) AnchorParallel(ASide TAnchorKind, ASpace int32, ASibling IControl) {
    DrawGrid_AnchorParallel(d.instance, ASide , ASpace , CheckPtr(ASibling))
}

// 置于指定控件的横向中心。
func (d *TDrawGrid) AnchorHorizontalCenterTo(ASibling IControl) {
    DrawGrid_AnchorHorizontalCenterTo(d.instance, CheckPtr(ASibling))
}

// 置于指定控件的纵向中心。
func (d *TDrawGrid) AnchorVerticalCenterTo(ASibling IControl) {
    DrawGrid_AnchorVerticalCenterTo(d.instance, CheckPtr(ASibling))
}

func (d *TDrawGrid) AnchorSame(ASide TAnchorKind, ASibling IControl) {
    DrawGrid_AnchorSame(d.instance, ASide , CheckPtr(ASibling))
}

func (d *TDrawGrid) AnchorAsAlign(ATheAlign TAlign, ASpace int32) {
    DrawGrid_AnchorAsAlign(d.instance, ATheAlign , ASpace)
}

func (d *TDrawGrid) AnchorClient(ASpace int32) {
    DrawGrid_AnchorClient(d.instance, ASpace)
}

func (d *TDrawGrid) ScaleDesignToForm(ASize int32) int32 {
    return DrawGrid_ScaleDesignToForm(d.instance, ASize)
}

func (d *TDrawGrid) ScaleFormToDesign(ASize int32) int32 {
    return DrawGrid_ScaleFormToDesign(d.instance, ASize)
}

func (d *TDrawGrid) Scale96ToForm(ASize int32) int32 {
    return DrawGrid_Scale96ToForm(d.instance, ASize)
}

func (d *TDrawGrid) ScaleFormTo96(ASize int32) int32 {
    return DrawGrid_ScaleFormTo96(d.instance, ASize)
}

func (d *TDrawGrid) Scale96ToFont(ASize int32) int32 {
    return DrawGrid_Scale96ToFont(d.instance, ASize)
}

func (d *TDrawGrid) ScaleFontTo96(ASize int32) int32 {
    return DrawGrid_ScaleFontTo96(d.instance, ASize)
}

func (d *TDrawGrid) ScaleScreenToFont(ASize int32) int32 {
    return DrawGrid_ScaleScreenToFont(d.instance, ASize)
}

func (d *TDrawGrid) ScaleFontToScreen(ASize int32) int32 {
    return DrawGrid_ScaleFontToScreen(d.instance, ASize)
}

func (d *TDrawGrid) Scale96ToScreen(ASize int32) int32 {
    return DrawGrid_Scale96ToScreen(d.instance, ASize)
}

func (d *TDrawGrid) ScaleScreenTo96(ASize int32) int32 {
    return DrawGrid_ScaleScreenTo96(d.instance, ASize)
}

func (d *TDrawGrid) AutoAdjustLayout(AMode TLayoutAdjustmentPolicy, AFromPPI int32, AToPPI int32, AOldFormWidth int32, ANewFormWidth int32) {
    DrawGrid_AutoAdjustLayout(d.instance, AMode , AFromPPI , AToPPI , AOldFormWidth , ANewFormWidth)
}

func (d *TDrawGrid) FixDesignFontsPPI(ADesignTimePPI int32) {
    DrawGrid_FixDesignFontsPPI(d.instance, ADesignTimePPI)
}

func (d *TDrawGrid) ScaleFontsPPI(AToPPI int32, AProportion float64) {
    DrawGrid_ScaleFontsPPI(d.instance, AToPPI , AProportion)
}

func (d *TDrawGrid) SetOnColRowMoved(fn TGridOperationEvent) {
    DrawGrid_SetOnColRowMoved(d.instance, fn)
}

func (d *TDrawGrid) SetOnPrepareCanvas(fn TOnPrepareCanvasEvent) {
    DrawGrid_SetOnPrepareCanvas(d.instance, fn)
}

// 获取控件自动调整。
//
// Get Control automatically adjusts.
func (d *TDrawGrid) Align() TAlign {
    return DrawGrid_GetAlign(d.instance)
}

// 设置控件自动调整。
//
// Set Control automatically adjusts.
func (d *TDrawGrid) SetAlign(value TAlign) {
    DrawGrid_SetAlign(d.instance, value)
}

// 获取四个角位置的锚点。
func (d *TDrawGrid) Anchors() TAnchors {
    return DrawGrid_GetAnchors(d.instance)
}

// 设置四个角位置的锚点。
func (d *TDrawGrid) SetAnchors(value TAnchors) {
    DrawGrid_SetAnchors(d.instance, value)
}

func (d *TDrawGrid) BiDiMode() TBiDiMode {
    return DrawGrid_GetBiDiMode(d.instance)
}

func (d *TDrawGrid) SetBiDiMode(value TBiDiMode) {
    DrawGrid_SetBiDiMode(d.instance, value)
}

// 获取窗口边框样式。比如：无边框，单一边框等。
func (d *TDrawGrid) BorderStyle() TBorderStyle {
    return DrawGrid_GetBorderStyle(d.instance)
}

// 设置窗口边框样式。比如：无边框，单一边框等。
func (d *TDrawGrid) SetBorderStyle(value TBorderStyle) {
    DrawGrid_SetBorderStyle(d.instance, value)
}

// 获取颜色。
//
// Get color.
func (d *TDrawGrid) Color() TColor {
    return DrawGrid_GetColor(d.instance)
}

// 设置颜色。
//
// Set color.
func (d *TDrawGrid) SetColor(value TColor) {
    DrawGrid_SetColor(d.instance, value)
}

func (d *TDrawGrid) ColCount() int32 {
    return DrawGrid_GetColCount(d.instance)
}

func (d *TDrawGrid) SetColCount(value int32) {
    DrawGrid_SetColCount(d.instance, value)
}

// 获取约束控件大小。
func (d *TDrawGrid) Constraints() *TSizeConstraints {
    return AsSizeConstraints(DrawGrid_GetConstraints(d.instance))
}

// 设置约束控件大小。
func (d *TDrawGrid) SetConstraints(value *TSizeConstraints) {
    DrawGrid_SetConstraints(d.instance, CheckPtr(value))
}

func (d *TDrawGrid) DefaultColWidth() int32 {
    return DrawGrid_GetDefaultColWidth(d.instance)
}

func (d *TDrawGrid) SetDefaultColWidth(value int32) {
    DrawGrid_SetDefaultColWidth(d.instance, value)
}

func (d *TDrawGrid) DefaultRowHeight() int32 {
    return DrawGrid_GetDefaultRowHeight(d.instance)
}

func (d *TDrawGrid) SetDefaultRowHeight(value int32) {
    DrawGrid_SetDefaultRowHeight(d.instance, value)
}

func (d *TDrawGrid) DefaultDrawing() bool {
    return DrawGrid_GetDefaultDrawing(d.instance)
}

func (d *TDrawGrid) SetDefaultDrawing(value bool) {
    DrawGrid_SetDefaultDrawing(d.instance, value)
}

// 获取设置控件双缓冲。
//
// Get Set control double buffering.
func (d *TDrawGrid) DoubleBuffered() bool {
    return DrawGrid_GetDoubleBuffered(d.instance)
}

// 设置设置控件双缓冲。
//
// Set Set control double buffering.
func (d *TDrawGrid) SetDoubleBuffered(value bool) {
    DrawGrid_SetDoubleBuffered(d.instance, value)
}

// 获取设置控件拖拽时的光标。
//
// Get Set the cursor when the control is dragged.
func (d *TDrawGrid) DragCursor() TCursor {
    return DrawGrid_GetDragCursor(d.instance)
}

// 设置设置控件拖拽时的光标。
//
// Set Set the cursor when the control is dragged.
func (d *TDrawGrid) SetDragCursor(value TCursor) {
    DrawGrid_SetDragCursor(d.instance, value)
}

// 获取拖拽方式。
//
// Get Drag and drop.
func (d *TDrawGrid) DragKind() TDragKind {
    return DrawGrid_GetDragKind(d.instance)
}

// 设置拖拽方式。
//
// Set Drag and drop.
func (d *TDrawGrid) SetDragKind(value TDragKind) {
    DrawGrid_SetDragKind(d.instance, value)
}

// 获取拖拽模式。
//
// Get Drag mode.
func (d *TDrawGrid) DragMode() TDragMode {
    return DrawGrid_GetDragMode(d.instance)
}

// 设置拖拽模式。
//
// Set Drag mode.
func (d *TDrawGrid) SetDragMode(value TDragMode) {
    DrawGrid_SetDragMode(d.instance, value)
}

// 获取控件启用。
//
// Get the control enabled.
func (d *TDrawGrid) Enabled() bool {
    return DrawGrid_GetEnabled(d.instance)
}

// 设置控件启用。
//
// Set the control enabled.
func (d *TDrawGrid) SetEnabled(value bool) {
    DrawGrid_SetEnabled(d.instance, value)
}

func (d *TDrawGrid) FixedColor() TColor {
    return DrawGrid_GetFixedColor(d.instance)
}

func (d *TDrawGrid) SetFixedColor(value TColor) {
    DrawGrid_SetFixedColor(d.instance, value)
}

func (d *TDrawGrid) FixedCols() int32 {
    return DrawGrid_GetFixedCols(d.instance)
}

func (d *TDrawGrid) SetFixedCols(value int32) {
    DrawGrid_SetFixedCols(d.instance, value)
}

func (d *TDrawGrid) RowCount() int32 {
    return DrawGrid_GetRowCount(d.instance)
}

func (d *TDrawGrid) SetRowCount(value int32) {
    DrawGrid_SetRowCount(d.instance, value)
}

func (d *TDrawGrid) FixedRows() int32 {
    return DrawGrid_GetFixedRows(d.instance)
}

func (d *TDrawGrid) SetFixedRows(value int32) {
    DrawGrid_SetFixedRows(d.instance, value)
}

// 获取字体。
//
// Get Font.
func (d *TDrawGrid) Font() *TFont {
    return AsFont(DrawGrid_GetFont(d.instance))
}

// 设置字体。
//
// Set Font.
func (d *TDrawGrid) SetFont(value *TFont) {
    DrawGrid_SetFont(d.instance, CheckPtr(value))
}

func (d *TDrawGrid) GridLineWidth() int32 {
    return DrawGrid_GetGridLineWidth(d.instance)
}

func (d *TDrawGrid) SetGridLineWidth(value int32) {
    DrawGrid_SetGridLineWidth(d.instance, value)
}

func (d *TDrawGrid) Options() TGridOptions {
    return DrawGrid_GetOptions(d.instance)
}

func (d *TDrawGrid) SetOptions(value TGridOptions) {
    DrawGrid_SetOptions(d.instance, value)
}

// 获取使用父容器颜色。
//
// Get parent color.
func (d *TDrawGrid) ParentColor() bool {
    return DrawGrid_GetParentColor(d.instance)
}

// 设置使用父容器颜色。
//
// Set parent color.
func (d *TDrawGrid) SetParentColor(value bool) {
    DrawGrid_SetParentColor(d.instance, value)
}

// 获取使用父容器双缓冲。
//
// Get Parent container double buffering.
func (d *TDrawGrid) ParentDoubleBuffered() bool {
    return DrawGrid_GetParentDoubleBuffered(d.instance)
}

// 设置使用父容器双缓冲。
//
// Set Parent container double buffering.
func (d *TDrawGrid) SetParentDoubleBuffered(value bool) {
    DrawGrid_SetParentDoubleBuffered(d.instance, value)
}

// 获取使用父容器字体。
//
// Get Parent container font.
func (d *TDrawGrid) ParentFont() bool {
    return DrawGrid_GetParentFont(d.instance)
}

// 设置使用父容器字体。
//
// Set Parent container font.
func (d *TDrawGrid) SetParentFont(value bool) {
    DrawGrid_SetParentFont(d.instance, value)
}

// 获取以父容器的ShowHint属性为准。
func (d *TDrawGrid) ParentShowHint() bool {
    return DrawGrid_GetParentShowHint(d.instance)
}

// 设置以父容器的ShowHint属性为准。
func (d *TDrawGrid) SetParentShowHint(value bool) {
    DrawGrid_SetParentShowHint(d.instance, value)
}

// 获取右键菜单。
//
// Get Right click menu.
func (d *TDrawGrid) PopupMenu() *TPopupMenu {
    return AsPopupMenu(DrawGrid_GetPopupMenu(d.instance))
}

// 设置右键菜单。
//
// Set Right click menu.
func (d *TDrawGrid) SetPopupMenu(value IComponent) {
    DrawGrid_SetPopupMenu(d.instance, CheckPtr(value))
}

func (d *TDrawGrid) ScrollBars() TScrollStyle {
    return DrawGrid_GetScrollBars(d.instance)
}

func (d *TDrawGrid) SetScrollBars(value TScrollStyle) {
    DrawGrid_SetScrollBars(d.instance, value)
}

// 获取显示鼠标悬停提示。
//
// Get Show mouseover tips.
func (d *TDrawGrid) ShowHint() bool {
    return DrawGrid_GetShowHint(d.instance)
}

// 设置显示鼠标悬停提示。
//
// Set Show mouseover tips.
func (d *TDrawGrid) SetShowHint(value bool) {
    DrawGrid_SetShowHint(d.instance, value)
}

// 获取Tab切换顺序序号。
//
// Get Tab switching sequence number.
func (d *TDrawGrid) TabOrder() TTabOrder {
    return DrawGrid_GetTabOrder(d.instance)
}

// 设置Tab切换顺序序号。
//
// Set Tab switching sequence number.
func (d *TDrawGrid) SetTabOrder(value TTabOrder) {
    DrawGrid_SetTabOrder(d.instance, value)
}

// 获取控件可视。
//
// Get the control visible.
func (d *TDrawGrid) Visible() bool {
    return DrawGrid_GetVisible(d.instance)
}

// 设置控件可视。
//
// Set the control visible.
func (d *TDrawGrid) SetVisible(value bool) {
    DrawGrid_SetVisible(d.instance, value)
}

func (d *TDrawGrid) VisibleColCount() int32 {
    return DrawGrid_GetVisibleColCount(d.instance)
}

func (d *TDrawGrid) VisibleRowCount() int32 {
    return DrawGrid_GetVisibleRowCount(d.instance)
}

// 设置控件单击事件。
//
// Set control click event.
func (d *TDrawGrid) SetOnClick(fn TNotifyEvent) {
    DrawGrid_SetOnClick(d.instance, fn)
}

// 设置上下文弹出事件，一般是右键时弹出。
//
// Set Context popup event, usually pop up when right click.
func (d *TDrawGrid) SetOnContextPopup(fn TContextPopupEvent) {
    DrawGrid_SetOnContextPopup(d.instance, fn)
}

// 设置双击事件。
func (d *TDrawGrid) SetOnDblClick(fn TNotifyEvent) {
    DrawGrid_SetOnDblClick(d.instance, fn)
}

// 设置拖拽下落事件。
//
// Set Drag and drop event.
func (d *TDrawGrid) SetOnDragDrop(fn TDragDropEvent) {
    DrawGrid_SetOnDragDrop(d.instance, fn)
}

// 设置拖拽完成事件。
//
// Set Drag and drop completion event.
func (d *TDrawGrid) SetOnDragOver(fn TDragOverEvent) {
    DrawGrid_SetOnDragOver(d.instance, fn)
}

func (d *TDrawGrid) SetOnDrawCell(fn TDrawCellEvent) {
    DrawGrid_SetOnDrawCell(d.instance, fn)
}

// 设置停靠结束事件。
//
// Set Dock end event.
func (d *TDrawGrid) SetOnEndDock(fn TEndDragEvent) {
    DrawGrid_SetOnEndDock(d.instance, fn)
}

// 设置拖拽结束。
//
// Set End of drag.
func (d *TDrawGrid) SetOnEndDrag(fn TEndDragEvent) {
    DrawGrid_SetOnEndDrag(d.instance, fn)
}

// 设置焦点进入。
//
// Set Focus entry.
func (d *TDrawGrid) SetOnEnter(fn TNotifyEvent) {
    DrawGrid_SetOnEnter(d.instance, fn)
}

// 设置焦点退出。
//
// Set Focus exit.
func (d *TDrawGrid) SetOnExit(fn TNotifyEvent) {
    DrawGrid_SetOnExit(d.instance, fn)
}

func (d *TDrawGrid) SetOnGetEditMask(fn TGetEditEvent) {
    DrawGrid_SetOnGetEditMask(d.instance, fn)
}

func (d *TDrawGrid) SetOnGetEditText(fn TGetEditEvent) {
    DrawGrid_SetOnGetEditText(d.instance, fn)
}

// 设置键盘按键按下事件。
//
// Set Keyboard button press event.
func (d *TDrawGrid) SetOnKeyDown(fn TKeyEvent) {
    DrawGrid_SetOnKeyDown(d.instance, fn)
}

// 设置键键下事件。
func (d *TDrawGrid) SetOnKeyPress(fn TKeyPressEvent) {
    DrawGrid_SetOnKeyPress(d.instance, fn)
}

// 设置键盘按键抬起事件。
//
// Set Keyboard button lift event.
func (d *TDrawGrid) SetOnKeyUp(fn TKeyEvent) {
    DrawGrid_SetOnKeyUp(d.instance, fn)
}

// 设置鼠标按下事件。
//
// Set Mouse down event.
func (d *TDrawGrid) SetOnMouseDown(fn TMouseEvent) {
    DrawGrid_SetOnMouseDown(d.instance, fn)
}

// 设置鼠标进入事件。
//
// Set Mouse entry event.
func (d *TDrawGrid) SetOnMouseEnter(fn TNotifyEvent) {
    DrawGrid_SetOnMouseEnter(d.instance, fn)
}

// 设置鼠标离开事件。
//
// Set Mouse leave event.
func (d *TDrawGrid) SetOnMouseLeave(fn TNotifyEvent) {
    DrawGrid_SetOnMouseLeave(d.instance, fn)
}

// 设置鼠标移动事件。
func (d *TDrawGrid) SetOnMouseMove(fn TMouseMoveEvent) {
    DrawGrid_SetOnMouseMove(d.instance, fn)
}

// 设置鼠标抬起事件。
//
// Set Mouse lift event.
func (d *TDrawGrid) SetOnMouseUp(fn TMouseEvent) {
    DrawGrid_SetOnMouseUp(d.instance, fn)
}

// 设置鼠标滚轮按下事件。
func (d *TDrawGrid) SetOnMouseWheelDown(fn TMouseWheelUpDownEvent) {
    DrawGrid_SetOnMouseWheelDown(d.instance, fn)
}

// 设置鼠标滚轮抬起事件。
func (d *TDrawGrid) SetOnMouseWheelUp(fn TMouseWheelUpDownEvent) {
    DrawGrid_SetOnMouseWheelUp(d.instance, fn)
}

func (d *TDrawGrid) SetOnSelectCell(fn TSelectCellEvent) {
    DrawGrid_SetOnSelectCell(d.instance, fn)
}

func (d *TDrawGrid) SetOnSetEditText(fn TSetEditEvent) {
    DrawGrid_SetOnSetEditText(d.instance, fn)
}

// 设置启动停靠。
func (d *TDrawGrid) SetOnStartDock(fn TStartDockEvent) {
    DrawGrid_SetOnStartDock(d.instance, fn)
}

func (d *TDrawGrid) SetOnTopLeftChanged(fn TNotifyEvent) {
    DrawGrid_SetOnTopLeftChanged(d.instance, fn)
}

// 获取画布。
func (d *TDrawGrid) Canvas() *TCanvas {
    return AsCanvas(DrawGrid_GetCanvas(d.instance))
}

func (d *TDrawGrid) Col() int32 {
    return DrawGrid_GetCol(d.instance)
}

func (d *TDrawGrid) SetCol(value int32) {
    DrawGrid_SetCol(d.instance, value)
}

func (d *TDrawGrid) EditorMode() bool {
    return DrawGrid_GetEditorMode(d.instance)
}

func (d *TDrawGrid) SetEditorMode(value bool) {
    DrawGrid_SetEditorMode(d.instance, value)
}

func (d *TDrawGrid) GridHeight() int32 {
    return DrawGrid_GetGridHeight(d.instance)
}

func (d *TDrawGrid) GridWidth() int32 {
    return DrawGrid_GetGridWidth(d.instance)
}

func (d *TDrawGrid) LeftCol() int32 {
    return DrawGrid_GetLeftCol(d.instance)
}

func (d *TDrawGrid) SetLeftCol(value int32) {
    DrawGrid_SetLeftCol(d.instance, value)
}

func (d *TDrawGrid) Selection() TGridRect {
    return DrawGrid_GetSelection(d.instance)
}

func (d *TDrawGrid) SetSelection(value TGridRect) {
    DrawGrid_SetSelection(d.instance, value)
}

func (d *TDrawGrid) Row() int32 {
    return DrawGrid_GetRow(d.instance)
}

func (d *TDrawGrid) SetRow(value int32) {
    DrawGrid_SetRow(d.instance, value)
}

func (d *TDrawGrid) TopRow() int32 {
    return DrawGrid_GetTopRow(d.instance)
}

func (d *TDrawGrid) SetTopRow(value int32) {
    DrawGrid_SetTopRow(d.instance, value)
}

// 获取Tab可停留。
//
// Get Tab can stay.
func (d *TDrawGrid) TabStop() bool {
    return DrawGrid_GetTabStop(d.instance)
}

// 设置Tab可停留。
//
// Set Tab can stay.
func (d *TDrawGrid) SetTabStop(value bool) {
    DrawGrid_SetTabStop(d.instance, value)
}

// 获取依靠客户端总数。
func (d *TDrawGrid) DockClientCount() int32 {
    return DrawGrid_GetDockClientCount(d.instance)
}

// 获取停靠站点。
//
// Get Docking site.
func (d *TDrawGrid) DockSite() bool {
    return DrawGrid_GetDockSite(d.instance)
}

// 设置停靠站点。
//
// Set Docking site.
func (d *TDrawGrid) SetDockSite(value bool) {
    DrawGrid_SetDockSite(d.instance, value)
}

// 获取鼠标是否在客户端，仅VCL有效。
//
// Get Whether the mouse is on the client, only VCL is valid.
func (d *TDrawGrid) MouseInClient() bool {
    return DrawGrid_GetMouseInClient(d.instance)
}

// 获取当前停靠的可视总数。
//
// Get The total number of visible calls currently docked.
func (d *TDrawGrid) VisibleDockClientCount() int32 {
    return DrawGrid_GetVisibleDockClientCount(d.instance)
}

// 获取画刷对象。
//
// Get Brush.
func (d *TDrawGrid) Brush() *TBrush {
    return AsBrush(DrawGrid_GetBrush(d.instance))
}

// 获取子控件数。
//
// Get Number of child controls.
func (d *TDrawGrid) ControlCount() int32 {
    return DrawGrid_GetControlCount(d.instance)
}

// 获取控件句柄。
//
// Get Control handle.
func (d *TDrawGrid) Handle() HWND {
    return DrawGrid_GetHandle(d.instance)
}

// 获取父容器句柄。
//
// Get Parent container handle.
func (d *TDrawGrid) ParentWindow() HWND {
    return DrawGrid_GetParentWindow(d.instance)
}

// 设置父容器句柄。
//
// Set Parent container handle.
func (d *TDrawGrid) SetParentWindow(value HWND) {
    DrawGrid_SetParentWindow(d.instance, value)
}

func (d *TDrawGrid) Showing() bool {
    return DrawGrid_GetShowing(d.instance)
}

// 获取使用停靠管理。
func (d *TDrawGrid) UseDockManager() bool {
    return DrawGrid_GetUseDockManager(d.instance)
}

// 设置使用停靠管理。
func (d *TDrawGrid) SetUseDockManager(value bool) {
    DrawGrid_SetUseDockManager(d.instance, value)
}

func (d *TDrawGrid) Action() *TAction {
    return AsAction(DrawGrid_GetAction(d.instance))
}

func (d *TDrawGrid) SetAction(value IComponent) {
    DrawGrid_SetAction(d.instance, CheckPtr(value))
}

func (d *TDrawGrid) BoundsRect() TRect {
    return DrawGrid_GetBoundsRect(d.instance)
}

func (d *TDrawGrid) SetBoundsRect(value TRect) {
    DrawGrid_SetBoundsRect(d.instance, value)
}

// 获取客户区高度。
//
// Get client height.
func (d *TDrawGrid) ClientHeight() int32 {
    return DrawGrid_GetClientHeight(d.instance)
}

// 设置客户区高度。
//
// Set client height.
func (d *TDrawGrid) SetClientHeight(value int32) {
    DrawGrid_SetClientHeight(d.instance, value)
}

func (d *TDrawGrid) ClientOrigin() TPoint {
    return DrawGrid_GetClientOrigin(d.instance)
}

// 获取客户区矩形。
//
// Get client rectangle.
func (d *TDrawGrid) ClientRect() TRect {
    return DrawGrid_GetClientRect(d.instance)
}

// 获取客户区宽度。
//
// Get client width.
func (d *TDrawGrid) ClientWidth() int32 {
    return DrawGrid_GetClientWidth(d.instance)
}

// 设置客户区宽度。
//
// Set client width.
func (d *TDrawGrid) SetClientWidth(value int32) {
    DrawGrid_SetClientWidth(d.instance, value)
}

// 获取控件状态。
//
// Get control state.
func (d *TDrawGrid) ControlState() TControlState {
    return DrawGrid_GetControlState(d.instance)
}

// 设置控件状态。
//
// Set control state.
func (d *TDrawGrid) SetControlState(value TControlState) {
    DrawGrid_SetControlState(d.instance, value)
}

// 获取控件样式。
//
// Get control style.
func (d *TDrawGrid) ControlStyle() TControlStyle {
    return DrawGrid_GetControlStyle(d.instance)
}

// 设置控件样式。
//
// Set control style.
func (d *TDrawGrid) SetControlStyle(value TControlStyle) {
    DrawGrid_SetControlStyle(d.instance, value)
}

func (d *TDrawGrid) Floating() bool {
    return DrawGrid_GetFloating(d.instance)
}

// 获取控件父容器。
//
// Get control parent container.
func (d *TDrawGrid) Parent() *TWinControl {
    return AsWinControl(DrawGrid_GetParent(d.instance))
}

// 设置控件父容器。
//
// Set control parent container.
func (d *TDrawGrid) SetParent(value IWinControl) {
    DrawGrid_SetParent(d.instance, CheckPtr(value))
}

// 获取左边位置。
//
// Get Left position.
func (d *TDrawGrid) Left() int32 {
    return DrawGrid_GetLeft(d.instance)
}

// 设置左边位置。
//
// Set Left position.
func (d *TDrawGrid) SetLeft(value int32) {
    DrawGrid_SetLeft(d.instance, value)
}

// 获取顶边位置。
//
// Get Top position.
func (d *TDrawGrid) Top() int32 {
    return DrawGrid_GetTop(d.instance)
}

// 设置顶边位置。
//
// Set Top position.
func (d *TDrawGrid) SetTop(value int32) {
    DrawGrid_SetTop(d.instance, value)
}

// 获取宽度。
//
// Get width.
func (d *TDrawGrid) Width() int32 {
    return DrawGrid_GetWidth(d.instance)
}

// 设置宽度。
//
// Set width.
func (d *TDrawGrid) SetWidth(value int32) {
    DrawGrid_SetWidth(d.instance, value)
}

// 获取高度。
//
// Get height.
func (d *TDrawGrid) Height() int32 {
    return DrawGrid_GetHeight(d.instance)
}

// 设置高度。
//
// Set height.
func (d *TDrawGrid) SetHeight(value int32) {
    DrawGrid_SetHeight(d.instance, value)
}

// 获取控件光标。
//
// Get control cursor.
func (d *TDrawGrid) Cursor() TCursor {
    return DrawGrid_GetCursor(d.instance)
}

// 设置控件光标。
//
// Set control cursor.
func (d *TDrawGrid) SetCursor(value TCursor) {
    DrawGrid_SetCursor(d.instance, value)
}

// 获取组件鼠标悬停提示。
//
// Get component mouse hints.
func (d *TDrawGrid) Hint() string {
    return DrawGrid_GetHint(d.instance)
}

// 设置组件鼠标悬停提示。
//
// Set component mouse hints.
func (d *TDrawGrid) SetHint(value string) {
    DrawGrid_SetHint(d.instance, value)
}

// 获取组件总数。
//
// Get the total number of components.
func (d *TDrawGrid) ComponentCount() int32 {
    return DrawGrid_GetComponentCount(d.instance)
}

// 获取组件索引。
//
// Get component index.
func (d *TDrawGrid) ComponentIndex() int32 {
    return DrawGrid_GetComponentIndex(d.instance)
}

// 设置组件索引。
//
// Set component index.
func (d *TDrawGrid) SetComponentIndex(value int32) {
    DrawGrid_SetComponentIndex(d.instance, value)
}

// 获取组件所有者。
//
// Get component owner.
func (d *TDrawGrid) Owner() *TComponent {
    return AsComponent(DrawGrid_GetOwner(d.instance))
}

// 获取组件名称。
//
// Get the component name.
func (d *TDrawGrid) Name() string {
    return DrawGrid_GetName(d.instance)
}

// 设置组件名称。
//
// Set the component name.
func (d *TDrawGrid) SetName(value string) {
    DrawGrid_SetName(d.instance, value)
}

// 获取对象标记。
//
// Get the control tag.
func (d *TDrawGrid) Tag() int {
    return DrawGrid_GetTag(d.instance)
}

// 设置对象标记。
//
// Set the control tag.
func (d *TDrawGrid) SetTag(value int) {
    DrawGrid_SetTag(d.instance, value)
}

// 获取左边锚点。
func (d *TDrawGrid) AnchorSideLeft() *TAnchorSide {
    return AsAnchorSide(DrawGrid_GetAnchorSideLeft(d.instance))
}

// 设置左边锚点。
func (d *TDrawGrid) SetAnchorSideLeft(value *TAnchorSide) {
    DrawGrid_SetAnchorSideLeft(d.instance, CheckPtr(value))
}

// 获取顶边锚点。
func (d *TDrawGrid) AnchorSideTop() *TAnchorSide {
    return AsAnchorSide(DrawGrid_GetAnchorSideTop(d.instance))
}

// 设置顶边锚点。
func (d *TDrawGrid) SetAnchorSideTop(value *TAnchorSide) {
    DrawGrid_SetAnchorSideTop(d.instance, CheckPtr(value))
}

// 获取右边锚点。
func (d *TDrawGrid) AnchorSideRight() *TAnchorSide {
    return AsAnchorSide(DrawGrid_GetAnchorSideRight(d.instance))
}

// 设置右边锚点。
func (d *TDrawGrid) SetAnchorSideRight(value *TAnchorSide) {
    DrawGrid_SetAnchorSideRight(d.instance, CheckPtr(value))
}

// 获取底边锚点。
func (d *TDrawGrid) AnchorSideBottom() *TAnchorSide {
    return AsAnchorSide(DrawGrid_GetAnchorSideBottom(d.instance))
}

// 设置底边锚点。
func (d *TDrawGrid) SetAnchorSideBottom(value *TAnchorSide) {
    DrawGrid_SetAnchorSideBottom(d.instance, CheckPtr(value))
}

func (d *TDrawGrid) ChildSizing() *TControlChildSizing {
    return AsControlChildSizing(DrawGrid_GetChildSizing(d.instance))
}

func (d *TDrawGrid) SetChildSizing(value *TControlChildSizing) {
    DrawGrid_SetChildSizing(d.instance, CheckPtr(value))
}

// 获取边框间距。
func (d *TDrawGrid) BorderSpacing() *TControlBorderSpacing {
    return AsControlBorderSpacing(DrawGrid_GetBorderSpacing(d.instance))
}

// 设置边框间距。
func (d *TDrawGrid) SetBorderSpacing(value *TControlBorderSpacing) {
    DrawGrid_SetBorderSpacing(d.instance, CheckPtr(value))
}

func (d *TDrawGrid) ColWidths(Index int32) int32 {
    return DrawGrid_GetColWidths(d.instance, Index)
}

func (d *TDrawGrid) SetColWidths(Index int32, value int32) {
    DrawGrid_SetColWidths(d.instance, Index, value)
}

func (d *TDrawGrid) RowHeights(Index int32) int32 {
    return DrawGrid_GetRowHeights(d.instance, Index)
}

func (d *TDrawGrid) SetRowHeights(Index int32, value int32) {
    DrawGrid_SetRowHeights(d.instance, Index, value)
}

// 获取指定索引停靠客户端。
func (d *TDrawGrid) DockClients(Index int32) *TControl {
    return AsControl(DrawGrid_GetDockClients(d.instance, Index))
}

// 获取指定索引子控件。
func (d *TDrawGrid) Controls(Index int32) *TControl {
    return AsControl(DrawGrid_GetControls(d.instance, Index))
}

// 获取指定索引组件。
//
// Get the specified index component.
func (d *TDrawGrid) Components(AIndex int32) *TComponent {
    return AsComponent(DrawGrid_GetComponents(d.instance, AIndex))
}

// 获取锚侧面。
func (d *TDrawGrid) AnchorSide(AKind TAnchorKind) *TAnchorSide {
    return AsAnchorSide(DrawGrid_GetAnchorSide(d.instance, AKind))
}
