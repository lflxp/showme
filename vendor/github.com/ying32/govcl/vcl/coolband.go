
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

type TCoolBand struct {
    IObject
    instance uintptr
    // 特殊情况下使用，主要应对Go的GC问题，与LCL没有太多关系。
    ptr unsafe.Pointer
}

// 创建一个新的对象。
// 
// Create a new object.
func NewCoolBand(AOwner *TCollection) *TCoolBand {
    c := new(TCoolBand)
    c.instance = CoolBand_Create(CheckPtr(AOwner))
    c.ptr = unsafe.Pointer(c.instance)
    setFinalizer(c, (*TCoolBand).Free)
    return c
}

// 动态转换一个已存在的对象实例。
// 
// Dynamically convert an existing object instance.
func AsCoolBand(obj interface{}) *TCoolBand {
    instance, ptr := getInstance(obj)
    if instance == 0 { return nil }
    return &TCoolBand{instance: instance, ptr: ptr}
}

// -------------------------- Deprecated begin --------------------------
// 新建一个对象来自已经存在的对象实例指针。
// 
// Create a new object from an existing object instance pointer.
// Deprecated: use AsCoolBand.
func CoolBandFromInst(inst uintptr) *TCoolBand {
    return AsCoolBand(inst)
}

// 新建一个对象来自已经存在的对象实例。
// 
// Create a new object from an existing object instance.
// Deprecated: use AsCoolBand.
func CoolBandFromObj(obj IObject) *TCoolBand {
    return AsCoolBand(obj)
}

// 新建一个对象来自不安全的地址。注意：使用此函数可能造成一些不明情况，慎用。
// 
// Create a new object from an unsecured address. Note: Using this function may cause some unclear situations and be used with caution..
// Deprecated: use AsCoolBand.
func CoolBandFromUnsafePointer(ptr unsafe.Pointer) *TCoolBand {
    return AsCoolBand(ptr)
}

// -------------------------- Deprecated end --------------------------
// 释放对象。
// 
// Free object.
func (c *TCoolBand) Free() {
    if c.instance != 0 {
        CoolBand_Free(c.instance)
        c.instance, c.ptr = 0, nullptr
    }
}

// 返回对象实例指针。
// 
// Return object instance pointer.
func (c *TCoolBand) Instance() uintptr {
    return c.instance
}

// 获取一个不安全的地址。
// 
// Get an unsafe address.
func (c *TCoolBand) UnsafeAddr() unsafe.Pointer {
    return c.ptr
}

// 检测地址是否为空。
// 
// Check if the address is empty.
func (c *TCoolBand) IsValid() bool {
    return c.instance != 0
}

// 检测当前对象是否继承自目标对象。
// 
// Checks whether the current object is inherited from the target object.
func (c *TCoolBand) Is() TIs {
    return TIs(c.instance)
}

// 动态转换当前对象为目标对象。
// 
// Dynamically convert the current object to the target object.
//func (c *TCoolBand) As() TAs {
//    return TAs(c.instance)
//}

// 获取类信息指针。
// 
// Get class information pointer.
func TCoolBandClass() TClass {
    return CoolBand_StaticClassType()
}

// 复制一个对象，如果对象实现了此方法的话。
//
// Copy an object, if the object implements this method.
func (c *TCoolBand) Assign(Source IObject) {
    CoolBand_Assign(c.instance, CheckPtr(Source))
}

// 获取类名路径。
//
// Get the class name path.
func (c *TCoolBand) GetNamePath() string {
    return CoolBand_GetNamePath(c.instance)
}

// 获取类的类型信息。
//
// Get class type information.
func (c *TCoolBand) ClassType() TClass {
    return CoolBand_ClassType(c.instance)
}

// 获取当前对象类名称。
//
// Get the current object class name.
func (c *TCoolBand) ClassName() string {
    return CoolBand_ClassName(c.instance)
}

// 获取当前对象实例大小。
//
// Get the current object instance size.
func (c *TCoolBand) InstanceSize() int32 {
    return CoolBand_InstanceSize(c.instance)
}

// 判断当前类是否继承自指定类。
//
// Determine whether the current class inherits from the specified class.
func (c *TCoolBand) InheritsFrom(AClass TClass) bool {
    return CoolBand_InheritsFrom(c.instance, AClass)
}

// 与一个对象进行比较。
//
// Compare with an object.
func (c *TCoolBand) Equals(Obj IObject) bool {
    return CoolBand_Equals(c.instance, CheckPtr(Obj))
}

// 获取类的哈希值。
//
// Get the hash value of the class.
func (c *TCoolBand) GetHashCode() int32 {
    return CoolBand_GetHashCode(c.instance)
}

// 文本类信息。
//
// Text information.
func (c *TCoolBand) ToString() string {
    return CoolBand_ToString(c.instance)
}

// 获取高度。
//
// Get height.
func (c *TCoolBand) Height() int32 {
    return CoolBand_GetHeight(c.instance)
}

func (c *TCoolBand) Bitmap() *TBitmap {
    return AsBitmap(CoolBand_GetBitmap(c.instance))
}

func (c *TCoolBand) SetBitmap(value *TBitmap) {
    CoolBand_SetBitmap(c.instance, CheckPtr(value))
}

// 获取窗口边框样式。比如：无边框，单一边框等。
func (c *TCoolBand) BorderStyle() TBorderStyle {
    return CoolBand_GetBorderStyle(c.instance)
}

// 设置窗口边框样式。比如：无边框，单一边框等。
func (c *TCoolBand) SetBorderStyle(value TBorderStyle) {
    CoolBand_SetBorderStyle(c.instance, value)
}

func (c *TCoolBand) Break() bool {
    return CoolBand_GetBreak(c.instance)
}

func (c *TCoolBand) SetBreak(value bool) {
    CoolBand_SetBreak(c.instance, value)
}

// 获取颜色。
//
// Get color.
func (c *TCoolBand) Color() TColor {
    return CoolBand_GetColor(c.instance)
}

// 设置颜色。
//
// Set color.
func (c *TCoolBand) SetColor(value TColor) {
    CoolBand_SetColor(c.instance, value)
}

func (c *TCoolBand) Control() *TWinControl {
    return AsWinControl(CoolBand_GetControl(c.instance))
}

func (c *TCoolBand) SetControl(value IWinControl) {
    CoolBand_SetControl(c.instance, CheckPtr(value))
}

func (c *TCoolBand) FixedBackground() bool {
    return CoolBand_GetFixedBackground(c.instance)
}

func (c *TCoolBand) SetFixedBackground(value bool) {
    CoolBand_SetFixedBackground(c.instance, value)
}

func (c *TCoolBand) FixedSize() bool {
    return CoolBand_GetFixedSize(c.instance)
}

func (c *TCoolBand) SetFixedSize(value bool) {
    CoolBand_SetFixedSize(c.instance, value)
}

func (c *TCoolBand) HorizontalOnly() bool {
    return CoolBand_GetHorizontalOnly(c.instance)
}

func (c *TCoolBand) SetHorizontalOnly(value bool) {
    CoolBand_SetHorizontalOnly(c.instance, value)
}

// 获取图像在images中的索引。
func (c *TCoolBand) ImageIndex() int32 {
    return CoolBand_GetImageIndex(c.instance)
}

// 设置图像在images中的索引。
func (c *TCoolBand) SetImageIndex(value int32) {
    CoolBand_SetImageIndex(c.instance, value)
}

func (c *TCoolBand) MinHeight() int32 {
    return CoolBand_GetMinHeight(c.instance)
}

func (c *TCoolBand) SetMinHeight(value int32) {
    CoolBand_SetMinHeight(c.instance, value)
}

func (c *TCoolBand) MinWidth() int32 {
    return CoolBand_GetMinWidth(c.instance)
}

func (c *TCoolBand) SetMinWidth(value int32) {
    CoolBand_SetMinWidth(c.instance, value)
}

// 获取使用父容器颜色。
//
// Get parent color.
func (c *TCoolBand) ParentColor() bool {
    return CoolBand_GetParentColor(c.instance)
}

// 设置使用父容器颜色。
//
// Set parent color.
func (c *TCoolBand) SetParentColor(value bool) {
    CoolBand_SetParentColor(c.instance, value)
}

func (c *TCoolBand) ParentBitmap() bool {
    return CoolBand_GetParentBitmap(c.instance)
}

func (c *TCoolBand) SetParentBitmap(value bool) {
    CoolBand_SetParentBitmap(c.instance, value)
}

// 获取文本。
func (c *TCoolBand) Text() string {
    return CoolBand_GetText(c.instance)
}

// 设置文本。
func (c *TCoolBand) SetText(value string) {
    CoolBand_SetText(c.instance, value)
}

// 获取控件可视。
//
// Get the control visible.
func (c *TCoolBand) Visible() bool {
    return CoolBand_GetVisible(c.instance)
}

// 设置控件可视。
//
// Set the control visible.
func (c *TCoolBand) SetVisible(value bool) {
    CoolBand_SetVisible(c.instance, value)
}

// 获取宽度。
//
// Get width.
func (c *TCoolBand) Width() int32 {
    return CoolBand_GetWidth(c.instance)
}

// 设置宽度。
//
// Set width.
func (c *TCoolBand) SetWidth(value int32) {
    CoolBand_SetWidth(c.instance, value)
}

func (c *TCoolBand) Collection() *TCollection {
    return AsCollection(CoolBand_GetCollection(c.instance))
}

func (c *TCoolBand) SetCollection(value *TCollection) {
    CoolBand_SetCollection(c.instance, CheckPtr(value))
}

func (c *TCoolBand) Index() int32 {
    return CoolBand_GetIndex(c.instance)
}

func (c *TCoolBand) SetIndex(value int32) {
    CoolBand_SetIndex(c.instance, value)
}

func (c *TCoolBand) DisplayName() string {
    return CoolBand_GetDisplayName(c.instance)
}

func (c *TCoolBand) SetDisplayName(value string) {
    CoolBand_SetDisplayName(c.instance, value)
}
