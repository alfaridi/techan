package techan

import (
	"testing"
	"time"

	"github.com/sdcoffey/big"
	"github.com/stretchr/testify/assert"
)

func TestTimeSeries_AddCandle(t *testing.T) {
	t.Run("Throws if nil candle passed", func(t *testing.T) {
		ts := NewTimeSeries()
		assert.Panics(t, func() {
			ts.AddCandle(nil)
		})
	})

	t.Run("Adds candle if last is nil", func(t *testing.T) {
		ts := NewTimeSeries()

		now := time.Now()
		candle := NewCandle(NewTimePeriod(now, now.Add(time.Minute)))
		candle.ClosePrice = big.NewDecimal(1)

		ts.AddCandle(candle)

		assert.Len(t, ts.Candles, 1)
	})

	t.Run("Does not add candle if before last candle", func(t *testing.T) {
		ts := NewTimeSeries()

		now := time.Now()
		candle := NewCandle(NewTimePeriod(now, now.Add(time.Minute)))
		candle.ClosePrice = big.NewDecimal(1)

		ts.AddCandle(candle)
		then := now.Add(-time.Minute * 10)

		nextCandle := NewCandle(NewTimePeriod(then, now.Add(time.Minute)))
		candle.ClosePrice = big.NewDecimal(2)

		ts.AddCandle(nextCandle)

		assert.Len(t, ts.Candles, 1)
		assert.EqualValues(t, now.UnixNano(), ts.Candles[0].Period.Start.UnixNano())
	})
}

func TestTimeSeries_LastCandle(t *testing.T) {
	ts := NewTimeSeries()

	now := time.Now()
	candle := NewCandle(NewTimePeriod(now, now.Add(time.Minute)))
	candle.ClosePrice = big.NewDecimal(1)

	ts.AddCandle(candle)

	assert.EqualValues(t, now.UnixNano(), ts.LastCandle().Period.Start.UnixNano())
	assert.EqualValues(t, 1, ts.LastCandle().ClosePrice.Float())

	next := time.Now().Add(time.Minute)
	newCandle := NewCandle(NewTimePeriod(next, now.Add(time.Minute)))
	newCandle.ClosePrice = big.NewDecimal(2)

	ts.AddCandle(newCandle)

	assert.Len(t, ts.Candles, 2)

	assert.EqualValues(t, next.UnixNano(), ts.LastCandle().Period.Start.UnixNano())
	assert.EqualValues(t, 2, ts.LastCandle().ClosePrice.Float())
}

func TestTimeSeries_LastIndex(t *testing.T) {
	ts := NewTimeSeries()

	now := time.Now()
	candle := NewCandle(NewTimePeriod(now, now.Add(time.Minute)))
	ts.AddCandle(candle)

	assert.EqualValues(t, 0, ts.LastIndex())

	now2 := time.Now().Add(time.Minute)
	candle = NewCandle(NewTimePeriod(now2, now2.Add(time.Minute)))
	ts.AddCandle(candle)

	assert.EqualValues(t, 1, ts.LastIndex())
}
