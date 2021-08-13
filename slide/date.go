package slide

import "time"

type DateOp struct {
	nowUTC time.Time
	nowJST time.Time
}

func newDateOp() *DateOp {
	now := time.Now()
	nowUTC := now.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := nowUTC.In(jst)

	return &DateOp{
		nowUTC: nowUTC,
		nowJST: nowJST,
	}
}

func (d *DateOp) getDateJST() string {
	return d.nowJST.Format("20060102150405")
}
