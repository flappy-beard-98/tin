package shares

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"time"
	"tin/adapter"
)

// Share Объект передачи информации об акции.
type Share struct {
	Figi                  string  `json:"figi,omitempty"`                     //Figi-идентификатор инструмента.
	Ticker                string  `json:"ticker,omitempty"`                   //Тикер инструмента.
	ClassCode             string  `json:"class_code,omitempty"`               //Класс-код (секция торгов).
	Isin                  string  `json:"isin,omitempty"`                     //Isin-идентификатор инструмента.
	Lot                   int32   `json:"lot,omitempty"`                      //Лотность инструмента. Возможно совершение операций только на количества ценной бумаги, кратные параметру *lot*. Подробнее: [лот](https://tinkoff.github.io/investAPI/glossary#lot)
	Currency              string  `json:"currency,omitempty"`                 //Валюта расчётов.
	Klong                 float64 `json:"klong,omitempty"`                    //Коэффициент ставки риска длинной позиции по инструменту.
	Kshort                float64 `json:"kshort,omitempty"`                   //Коэффициент ставки риска короткой позиции по инструменту.
	Dlong                 float64 `json:"dlong,omitempty"`                    //Ставка риска минимальной маржи в лонг. Подробнее: [ставка риска в лонг](https://help.tinkoff.ru/margin-trade/long/risk-rate/)
	Dshort                float64 `json:"dshort,omitempty"`                   //Ставка риска минимальной маржи в шорт. Подробнее: [ставка риска в шорт](https://help.tinkoff.ru/margin-trade/short/risk-rate/)
	DlongMin              float64 `json:"dlong_min,omitempty"`                //Ставка риска начальной маржи в лонг. Подробнее: [ставка риска в лонг](https://help.tinkoff.ru/margin-trade/long/risk-rate/)
	DshortMin             float64 `json:"dshort_min,omitempty"`               //Ставка риска начальной маржи в шорт. Подробнее: [ставка риска в шорт](https://help.tinkoff.ru/margin-trade/short/risk-rate/)
	ShortEnabledFlag      bool    `json:"short_enabled_flag,omitempty"`       //Признак доступности для операций в шорт.
	Name                  string  `json:"name,omitempty"`                     //Название инструмента.
	Exchange              string  `json:"exchange,omitempty"`                 //Торговая площадка.
	IpoDate               string  `json:"ipo_date,omitempty"`                 //Дата IPO акции в часовом поясе UTC.
	IssueSize             int64   `json:"issue_size,omitempty"`               //Размер выпуска.
	CountryOfRisk         string  `json:"country_of_risk,omitempty"`          //Код страны риска, т.е. страны, в которой компания ведёт основной бизнес.
	CountryOfRiskName     string  `json:"country_of_risk_name,omitempty"`     //Наименование страны риска, т.е. страны, в которой компания ведёт основной бизнес.
	Sector                string  `json:"sector,omitempty"`                   //Сектор экономики.
	IssueSizePlan         int64   `json:"issue_size_plan,omitempty"`          //Плановый размер выпуска.
	Nominal               float64 `json:"nominal,omitempty"`                  //Номинал.
	TradingStatus         string  `json:"trading_status,omitempty"`           //Текущий режим торгов инструмента.
	OtcFlag               bool    `json:"otc_flag,omitempty"`                 //Признак внебиржевой ценной бумаги.
	BuyAvailableFlag      bool    `json:"buy_available_flag,omitempty"`       //Признак доступности для покупки.
	SellAvailableFlag     bool    `json:"sell_available_flag,omitempty"`      //Признак доступности для продажи.
	DivYieldFlag          bool    `json:"div_yield_flag,omitempty"`           //Признак наличия дивидендной доходности.
	ShareType             string  `json:"share_type,omitempty"`               //Тип акции. Возможные значения: [ShareType](https://tinkoff.github.io/investAPI/instruments#sharetype)
	MinPriceIncrement     float64 `json:"min_price_increment,omitempty"`      //Шаг цены.
	ApiTradeAvailableFlag bool    `json:"api_trade_available_flag,omitempty"` //Параметр указывает на возможность торговать инструментом через API.
	Uid                   string  `json:"uid,omitempty"`                      //Уникальный идентификатор инструмента.
	RealExchange          string  `json:"real_exchange,omitempty"`            //Реальная площадка исполнения расчётов.
	PositionUid           string  `json:"position_uid,omitempty"`             //Уникальный идентификатор позиции инструмента.
	ForIisFlag            bool    `json:"for_iis_flag,omitempty"`             //Признак доступности для ИИС.
	ForQualInvestorFlag   bool    `json:"for_qual_investor_flag,omitempty"`   //Флаг отображающий доступность торговли инструментом только для квалифицированных инвесторов.
	WeekendFlag           bool    `json:"weekend_flag,omitempty"`             //Флаг отображающий доступность торговли инструментом по выходным
	BlockedTcaFlag        bool    `json:"blocked_tca_flag,omitempty"`         //Флаг заблокированного ТКС
	LiquidityFlag         bool    `json:"liquidity_flag,omitempty"`           //Флаг достаточной ликвидности
	First_1MinCandleDate  string  `json:"first_1min_candle_date,omitempty"`   //Дата первой минутной свечи.
	First_1DayCandleDate  string  `json:"first_1day_candle_date,omitempty"`   //Дата первой дневной свечи.
}

func New(share *investapi.Share) *Share {
	return &Share{
		Figi:                  share.Figi,
		Ticker:                share.Ticker,
		ClassCode:             share.ClassCode,
		Isin:                  share.Isin,
		Lot:                   share.Lot,
		Currency:              share.Currency,
		Klong:                 share.Klong.ToFloat(),
		Kshort:                share.Kshort.ToFloat(),
		Dlong:                 share.Dlong.ToFloat(),
		Dshort:                share.Dshort.ToFloat(),
		DlongMin:              share.DlongMin.ToFloat(),
		DshortMin:             share.DshortMin.ToFloat(),
		ShortEnabledFlag:      share.ShortEnabledFlag,
		Name:                  share.Name,
		Exchange:              share.Exchange,
		IpoDate:               share.IpoDate.AsTime().Format("2006-01-02"),
		IssueSize:             share.IssueSize,
		CountryOfRisk:         share.CountryOfRisk,
		CountryOfRiskName:     share.CountryOfRiskName,
		Sector:                share.Sector,
		IssueSizePlan:         share.IssueSizePlan,
		Nominal:               share.Nominal.ToFloat(),
		TradingStatus:         share.TradingStatus.String(),
		OtcFlag:               share.OtcFlag,
		BuyAvailableFlag:      share.BuyAvailableFlag,
		SellAvailableFlag:     share.SellAvailableFlag,
		DivYieldFlag:          share.DivYieldFlag,
		ShareType:             share.ShareType.String(),
		MinPriceIncrement:     share.MinPriceIncrement.ToFloat(),
		ApiTradeAvailableFlag: share.ApiTradeAvailableFlag,
		Uid:                   share.Uid,
		RealExchange:          share.RealExchange.String(),
		PositionUid:           share.PositionUid,
		ForIisFlag:            share.ForIisFlag,
		ForQualInvestorFlag:   share.ForQualInvestorFlag,
		WeekendFlag:           share.WeekendFlag,
		BlockedTcaFlag:        share.BlockedTcaFlag,
		LiquidityFlag:         share.LiquidityFlag,
		First_1MinCandleDate:  share.First_1MinCandleDate.AsTime().Format("2006-01-02"),
		First_1DayCandleDate:  share.First_1DayCandleDate.AsTime().Format("2006-01-02"),
	}
}

func (o *Share) GetIpoDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.IpoDate)
	return t
}

func (o *Share) GetFirst_1MinCandleDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.First_1MinCandleDate)
	return t
}

func (o *Share) GetFirst_1DayCandleDate() time.Time {
	t, _ := time.Parse("2006-01-02", o.First_1DayCandleDate)
	return t
}

type Shares struct {
	shares []Share
}

func NewShares() *Shares {
	return &Shares{
		shares: make([]Share, 0),
	}
}

//go:embed schema.sql
var schema string

func (o *Shares) Scheme(ctx context.Context, db *adapter.Db) error {
	_, err := db.ExecContext(ctx, schema)
	return err
}

func (o *Shares) Get() []Share {
	return o.shares
}

func (o *Shares) GetFigis() []string {
	figis := make([]string, 0)
	for _, v := range o.shares {
		figis = append(figis, v.Figi)
	}
	return figis
}

func (o *Shares) GetUids() []string {
	uids := make([]string, 0)
	for _, v := range o.shares {
		uids = append(uids, v.Uid)
	}
	return uids
}

//go:embed insert.sql
var insert string

func (o *Shares) Insert(ctx context.Context, db *adapter.Db) error {
	for _, s := range o.shares {
		_, err := db.DB.NamedExecContext(ctx, insert, s)
		if err != nil {
			return err
		}
	}
	return nil
}

//go:embed read.sql
var read string

func (o *Shares) ReadByCurrency(ctx context.Context, db *adapter.Db, currency string) error {
	byCurrency := fmt.Sprintf("%s\nWHERE Currency = ?", read)
	err := db.SelectContext(ctx, &o.shares, byCurrency, currency)
	return err
}

func (o *Shares) Import(client *adapter.TApi) error {
	service := client.NewInstrumentsServiceClient()

	response, err := service.Shares(investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE)

	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("empty shares response")
	}

	o.shares = make([]Share, 0)

	for _, v := range response.Instruments {
		o.shares = append(o.shares, *New(v))
	}
	return nil
}
