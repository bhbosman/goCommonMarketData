package fullMarketDataManagerService

type PublishFullMarketData struct {
	PublishInstrument string
	PublishTo         string
}

func NewPublishFullMarketData(publishInstrument string, publishTo string) *PublishFullMarketData {
	return &PublishFullMarketData{PublishInstrument: publishInstrument, PublishTo: publishTo}
}

//
//type FullMarketDepthBookInstruments struct {
//	Instruments []string
//}
