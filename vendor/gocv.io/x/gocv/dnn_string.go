package gocv

func (c NetBackendType) String() string {
	switch c {
	case NetBackendDefault:
		return ""
	case NetBackendHalide:
		return "halide"
	case NetBackendOpenVINO:
		return "openvino"
	case NetBackendOpenCV:
		return "opencv"
	case NetBackendVKCOM:
		return "vulkan"
	case NetBackendCUDA:
		return "cuda"
	}
	return ""
}

func (c NetTargetType) String() string {
	switch c {
	case NetTargetCPU:
		return "cpu"
	case NetTargetFP32:
		return "fp32"
	case NetTargetFP16:
		return "fp16"
	case NetTargetVPU:
		return "vpu"
	case NetTargetVulkan:
		return "vulkan"
	case NetTargetFPGA:
		return "fpga"
	case NetTargetCUDA:
		return "cuda"
	case NetTargetCUDAFP16:
		return "cudafp16"
	}
	return ""
}
