package logger

import (
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type gcpEncoder struct {
	zapcore.Encoder
}

type entryCaller struct {
	*zapcore.EntryCaller
}

func newGCPEncoder(cfg *zapcore.EncoderConfig) zapcore.Encoder {
	return gcpEncoder{
		Encoder: zapcore.NewJSONEncoder(*cfg),
	}
}

func (enc gcpEncoder) Clone() zapcore.Encoder {
	return gcpEncoder{
		Encoder: enc.Encoder.Clone(),
	}
}

//nolint:gocritic // Implement interface of Zap, so no need to check hugeParam
func (enc gcpEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	fields = append(fields, zap.Object("logging.googleapis.com/sourceLocation", entryCaller{EntryCaller: &ent.Caller}))
	return enc.Encoder.EncodeEntry(ent, fields)
}

func (ent entryCaller) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("function", ent.EntryCaller.Function)
	enc.AddString("file", ent.EntryCaller.TrimmedPath())
	enc.AddString("line", strconv.Itoa(ent.EntryCaller.Line))
	return nil
}
