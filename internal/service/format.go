package service

import (
	"go.uber.org/zap"
	"golang.org/x/tools/imports"
)

type FormatService struct {
	log *zap.Logger
}

func NewFormatService(log *zap.Logger) *FormatService {
	return &FormatService{log: log}
}

func (s *FormatService) FormatCode(code string) (string, error) {
	opts := &imports.Options{
		Fragment:   true,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: false,
	}

	formatted, err := imports.Process("", []byte(code), opts)
	if err != nil {
		s.log.Error("failed to format code", zap.Error(err))
		return code, err
	}

	return string(formatted), nil
}
