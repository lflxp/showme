
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

type TClipboard struct {
    IObject
    instance uintptr
    // 特殊情况下使用，主要应对Go的GC问题，与LCL没有太多关系。
    ptr unsafe.Pointer
}

// 创建一个新的对象。
// 
// Create a new object.
func NewClipboard() *TClipboard {
    c := new(TClipboard)
    c.instance = Clipboard_Create()
    c.ptr = unsafe.Pointer(c.instance)
    setFinalizer(c, (*TClipboard).Free)
    return c
}

// 动态转换一个已存在的对象实例。
// 
// Dynamically convert an existing object instance.
func AsClipboard(obj interface{}) *TClipboard {
    instance, ptr := getInstance(obj)
    if instance == 0 { return nil }
    return &TClipboard{instance: instance, ptr: ptr}
}

// -------------------------- Deprecated begin --------------------------
// 新建一个对象来自已经存在的对象实例指针。
// 
// Create a new object from an existing object instance pointer.
// Deprecated: use AsClipboard.
func ClipboardFromInst(inst uintptr) *TClipboard {
    return AsClipboard(inst)
}

// 新建一个对象来自已经存在的对象实例。
// 
// Create a new object from an existing object instance.
// Deprecated: use AsClipboard.
func ClipboardFromObj(obj IObject) *TClipboard {
    return AsClipboard(obj)
}

// 新建一个对象来自不安全的地址。注意：使用此函数可能造成一些不明情况，慎用。
// 
// Create a new object from an unsecured address. Note: Using this function may cause some unclear situations and be used with caution..
// Deprecated: use AsClipboard.
func ClipboardFromUnsafePointer(ptr unsafe.Pointer) *TClipboard {
    return AsClipboard(ptr)
}

// -------------------------- Deprecated end --------------------------
// 释放对象。
// 
// Free object.
func (c *TClipboard) Free() {
    if c.instance != 0 {
        Clipboard_Free(c.instance)
        c.instance, c.ptr = 0, nullptr
    }
}

// 返回对象实例指针。
// 
// Return object instance pointer.
func (c *TClipboard) Instance() uintptr {
    return c.instance
}

// 获取一个不安全的地址。
// 
// Get an unsafe address.
func (c *TClipboard) UnsafeAddr() unsafe.Pointer {
    return c.ptr
}

// 检测地址是否为空。
// 
// Check if the address is empty.
func (c *TClipboard) IsValid() bool {
    return c.instance != 0
}

// 检测当前对象是否继承自目标对象。
// 
// Checks whether the current object is inherited from the target object.
func (c *TClipboard) Is() TIs {
    return TIs(c.instance)
}

// 动态转换当前对象为目标对象。
// 
// Dynamically convert the current object to the target object.
//func (c *TClipboard) As() TAs {
//    return TAs(c.instance)
//}

// 获取类信息指针。
// 
// Get class information pointer.
func TClipboardClass() TClass {
    return Clipboard_StaticClassType()
}

func (c *TClipboard) FindPictureFormatID() TClipboardFormat {
    return Clipboard_FindPictureFormatID(c.instance)
}

func (c *TClipboard) FindFormatID(FormatName string) TClipboardFormat {
    return Clipboard_FindFormatID(c.instance, FormatName)
}

func (c *TClipboard) SupportedFormats(List IStrings) {
    Clipboard_SupportedFormats(c.instance, CheckPtr(List))
}

func (c *TClipboard) HasFormatName(FormatName string) bool {
    return Clipboard_HasFormatName(c.instance, FormatName)
}

func (c *TClipboard) HasPictureFormat() bool {
    return Clipboard_HasPictureFormat(c.instance)
}

func (c *TClipboard) SetAsHtml(Html string, PlainText string) {
    Clipboard_SetAsHtml(c.instance, Html , PlainText)
}

func (c *TClipboard) GetFormat(FormatID TClipboardFormat, Stream IStream) bool {
    return Clipboard_GetFormat(c.instance, FormatID , CheckPtr(Stream))
}

// 复制一个对象，如果对象实现了此方法的话。
//
// Copy an object, if the object implements this method.
func (c *TClipboard) Assign(Source IObject) {
    Clipboard_Assign(c.instance, CheckPtr(Source))
}

// 清除。
func (c *TClipboard) Clear() {
    Clipboard_Clear(c.instance)
}

// 关闭。
func (c *TClipboard) Close() {
    Clipboard_Close(c.instance)
}

func (c *TClipboard) Open() {
    Clipboard_Open(c.instance)
}

// 设置控件字符，如果有。
//
// Set control characters, if any.
func (c *TClipboard) SetTextBuf(Buffer string) {
    Clipboard_SetTextBuf(c.instance, Buffer)
}

// 获取类名路径。
//
// Get the class name path.
func (c *TClipboard) GetNamePath() string {
    return Clipboard_GetNamePath(c.instance)
}

// 获取类的类型信息。
//
// Get class type information.
func (c *TClipboard) ClassType() TClass {
    return Clipboard_ClassType(c.instance)
}

// 获取当前对象类名称。
//
// Get the current object class name.
func (c *TClipboard) ClassName() string {
    return Clipboard_ClassName(c.instance)
}

// 获取当前对象实例大小。
//
// Get the current object instance size.
func (c *TClipboard) InstanceSize() int32 {
    return Clipboard_InstanceSize(c.instance)
}

// 判断当前类是否继承自指定类。
//
// Determine whether the current class inherits from the specified class.
func (c *TClipboard) InheritsFrom(AClass TClass) bool {
    return Clipboard_InheritsFrom(c.instance, AClass)
}

// 与一个对象进行比较。
//
// Compare with an object.
func (c *TClipboard) Equals(Obj IObject) bool {
    return Clipboard_Equals(c.instance, CheckPtr(Obj))
}

// 获取类的哈希值。
//
// Get the hash value of the class.
func (c *TClipboard) GetHashCode() int32 {
    return Clipboard_GetHashCode(c.instance)
}

// 文本类信息。
//
// Text information.
func (c *TClipboard) ToString() string {
    return Clipboard_ToString(c.instance)
}

func (c *TClipboard) FormatCount() int32 {
    return Clipboard_GetFormatCount(c.instance)
}

func (c *TClipboard) Formats(Index int32) TClipboardFormat {
    return Clipboard_GetFormats(c.instance, Index)
}
