package dto

// BaseResponseInterface base response interface
type BaseResponseInterface interface {
	SetByteResponse(response []byte)
	GetByteResponse() []byte
}

// BaseResponse the struct, which should extend all external response structs
type BaseResponse struct {
	ByteResponse []byte
}

// SetByteResponse sets the byte response
func (r *BaseResponse) SetByteResponse(response []byte) {
	r.ByteResponse = response
}

// GetByteResponse returns the byte response
func (r *BaseResponse) GetByteResponse() []byte {
	return r.ByteResponse
}
