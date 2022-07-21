package fullMarketData

import "github.com/bhbosman/goCommonMarketData/fullMarketData/stream"

type FullMarketOrder struct {
	side   stream.OrderSide
	id     string
	price  float64
	volume float64
}

func newFullMarketOrder(side stream.OrderSide, id string, price float64, volume float64) *FullMarketOrder {
	return &FullMarketOrder{
		side:   side,
		id:     id,
		price:  price,
		volume: volume,
	}
}

func (self *FullMarketOrder) ReduceVolume(base float64) (leftOverVolume float64) {
	self.volume -= base
	return self.volume
}
