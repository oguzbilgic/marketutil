package marketutil

import (
	"github.com/oguzbilgic/market"
)

type LoggingMarket struct {
	m market.Market
}

func NewLoggingMarket(m market.Market) *LoggingMarket {
	return &LoggingMarket{m}
}

func (l *LoggingMarket) Ticker() (*market.Tick, error) {
	return l.m.Ticker()
}

func (l *LoggingMarket) OrderBook() ([]*market.Depth, error) {
	return l.m.OrderBook()
}
