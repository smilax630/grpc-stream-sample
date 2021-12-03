package entity //

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	options struct {
		level string
	}
	// Option represents NewLogger option.
	Option func(*options)
)

// WithLevel returns level option.
// Available level is https://github.com/uber-go/zap/blob/master/zapcore/level.go/L126.
func WithLevel(level string) Option {
	return func(opts *options) {
		opts.level = level
	}
}

// Logger ...
var Logger = newLogger()

// newLogger returns a zap logger.
func newLogger(opts ...Option) *zap.Logger {
	dopts := &options{level: "info"}
	for _, opt := range opts {
		opt(dopts)
	}
	level := new(zapcore.Level)
	if err := level.Set(dopts.level); err != nil {
		log.Fatalf("failed to new looger: %v", err)
	}

	build, err := newConfig(*level).Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("failed to new logger: %v", err)
	}

	return build
}

func newConfig(level zapcore.Level) zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    newEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "eventTime",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func encodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString("DEBUG")
	case zapcore.InfoLevel:
		enc.AppendString("INFO")
	case zapcore.WarnLevel:
		enc.AppendString("WARNING")
	case zapcore.ErrorLevel:
		enc.AppendString("ERROR")
	case zapcore.DPanicLevel:
		enc.AppendString("CRITICAL")
	case zapcore.PanicLevel:
		enc.AppendString("ALERT")
	case zapcore.FatalLevel:
		enc.AppendString("EMERGENCY")
	}
}

type walletByUser struct {
	Time          string `json:"time"`
	FanUUID       string `json:"fan_uuid"`
	Balance       int64  `json:"balance"`        // 合計
	PaidBalance   int64  `json:"paid_balance"`   // 有償
	EarnedBalance int64  `json:"earned_balance"` // 無償
	PointBalance  int64  `json:"point_balance"`  // point
}

type consumeSimplyLogBody struct {
	Time                string     `json:"time"`
	InfluecnerUUID      string     `json:"influecner_uuid"`
	FanUUID             string     `json:"fan_uuid"`
	ItemCode            string     `json:"item_code"`
	ItemCodeCoin        int64      `json:"item_code_coin"`
	SumConsumeCoin      int64      `json:"sum_consume_coin"`
	FanMeetingID        uint       `json:"fan_meeting_id"`
	FanMeetingStartTime *time.Time `json:"fan_meeting_start_time"`
	NumExtension        uint32     `json:"num_extension"`
	Wallet              *Wallet    `json:"wallet"`
}

// facade logging
type facadeSimplyLogBody struct {
	Time           string  `json:"time"`
	FanUUID        string  `json:"fan_uuid"`
	ProductCode    string  `json:"product_code"`
	ProductCodeYen int64   `json:"product_code_yen"`
	TransactionID  string  `json:"transaction_id"`
	Wallet         *Wallet `json:"wallet"`
}

// point logging
type pointSimplyLogBody struct {
	Time          string  `json:"time"`
	FanUUID       string  `json:"fan_uuid"`
	IncentiveCode string  `json:"incentive_code"`
	IncentiveName string  `json:"incentive_name"`
	Point         int     `json:"point"`
	ExpireDays    int     `json:"expire_days"`
	Wallet        *Wallet `json:"wallet"`
}

// Wallet is user wallet struct
type Wallet struct {
	WalletID      string `json:"wallet_id"`
	CurrencyID    string `json:"currency_id"`
	Balance       int64  `json:"balance"`        // 合計
	PaidBalance   int64  `json:"paid_balance"`   // 有償
	EarnedBalance int64  `json:"earned_balance"` // 無償
	PointBalance  int64  `json:"point_balance"`  // point
}

//SetWalletByUserLog ...
func SetWalletByUserLog(fanUUID string, balance, paid, earned, point int64) {
	data := walletByUser{
		FanUUID:       fanUUID,
		Balance:       balance,
		PaidBalance:   paid,
		EarnedBalance: earned,
		PointBalance:  point,
	}
	j, _ := json.Marshal(data)
	sdLogging.simplyWalletLogger.Log(logging.Entry{Payload: json.RawMessage(j)})
}

//microFormat is  time format
const microFormat = "2006-01-02T15:04:05.000+0900"

//stackdriverLogging is
type stackdriverLogging struct {
	client              *logging.Client
	simplyConsumeLogger *logging.Logger
	simplyFacadeLogger  *logging.Logger
	simplyPointLogger   *logging.Logger
	simplyWalletLogger  *logging.Logger
}

var sdLogging *stackdriverLogging

//NewStackdriverLogging is
func NewStackdriverLogging(simplyConsumeLogID, simplyFacadeLogID, simplyPointLogID, simplyWalletLogID string) {
	client, err := logging.NewClient(context.Background(), "cyberpal4545")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	sdLogging = &stackdriverLogging{
		client:              client,
		simplyConsumeLogger: client.Logger(simplyConsumeLogID),
		simplyFacadeLogger:  client.Logger(simplyFacadeLogID),
		simplyPointLogger:   client.Logger(simplyPointLogID),
		simplyWalletLogger:  client.Logger(simplyWalletLogID),
	}
}

//CloseStackdriverLogging ...
func CloseStackdriverLogging() {
	sdLogging.client.Close()
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)
