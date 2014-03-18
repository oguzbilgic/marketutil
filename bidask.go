package marketutil

import (
	"github.com/oguzbilgic/fpd"
	"github.com/oguzbilgic/market"
	"time"
)

type BidAsk struct {
	ask     *fpd.Decimal
	askTime time.Time

	bid     *fpd.Decimal
	bidTime time.Time

	lastTrade     *fpd.Decimal
	lastTradeTime time.Time
}

func NewBidAsk(sMarket market.StreamingMarket) *BidAsk {
	ba := &BidAsk{}

	tickChan := sMarket.NewTickChan()
	go func() {
		for {
			tick := <-tickChan

			if tick.Time.After(ba.askTime) {
				ba.ask = tick.Ask
				ba.askTime = tick.Time
			}

			if tick.Time.After(ba.bidTime) {
				ba.bid = tick.Bid
				ba.bidTime = tick.Time
			}
		}
	}()

	tradeChan := sMarket.NewTradeChan()
	go func() {
		for {
			trade := <-tradeChan

			if trade.Currency != market.USD {
				continue
			}

			ba.lastTrade = trade.Price

			if ba.bid == nil || ba.ask == nil {
				continue
			}

			if trade.Time.After(ba.bidTime) {
				if trade.Price.Cmp(ba.bid) != 1 {
					ba.bid = trade.Price
					ba.bidTime = trade.Time
				}
			}

			if trade.Time.After(ba.askTime) {
				if trade.Price.Cmp(ba.ask) != -1 {
					ba.ask = trade.Price
					ba.askTime = trade.Time
				}
			}
		}
	}()

	return ba
}

func (ba *BidAsk) Ask() *fpd.Decimal {
	return ba.ask
}

func (ba *BidAsk) Bid() *fpd.Decimal {
	return ba.bid
}

func (ba *BidAsk) LastTrade() *fpd.Decimal {
	return ba.lastTrade
}

func (ba *BidAsk) Freshness() time.Duration {
	return (time.Since(ba.bidTime) + time.Since(ba.askTime)) / 2
}
