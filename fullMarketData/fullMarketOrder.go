package fullMarketData

import "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"

type FullMarketOrder struct {
	Side      stream.OrderSide
	Id        string
	Price     float64
	Volume    float64
	ExtraData []byte
}

func newFullMarketOrder(
	side stream.OrderSide,
	id string,
	price float64,
	volume float64,
	extraData []byte,
) *FullMarketOrder {
	return &FullMarketOrder{
		Side:      side,
		Id:        id,
		Price:     price,
		Volume:    volume,
		ExtraData: extraData,
	}
}

func (self *FullMarketOrder) ReduceVolume(base float64) (leftOverVolume float64) {
	self.Volume -= base
	return self.Volume
}
