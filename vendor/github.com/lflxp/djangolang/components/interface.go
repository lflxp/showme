package components

// ============================处理流程接口梳理==============================
type FormFactory interface {
	Transfer() interface{} // 转换逻辑函数, 返回处理函数和错误信息
}
