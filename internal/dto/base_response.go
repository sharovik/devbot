package dto

type BaseResponseInterface interface {
	SetByteResponse(response []byte)
	GetByteResponse() []byte
}

type BaseResponse struct {
	ByteResponse []byte
}

func (r *BaseResponse) SetByteResponse(response []byte) {
	r.ByteResponse = response
}

func (r *BaseResponse) GetByteResponse() []byte {
	return r.ByteResponse
}
