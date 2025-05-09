package logs

type GormZapLogger struct {
	logger *Logger
}

func NewGormZapLogger() *GormZapLogger {
	return &GormZapLogger{
		logger: NewLogger(),
	}
}

func (l *GormZapLogger) Printf(prefix string, v ...interface{}) {
	if len(v) == 4 {
		l.logger.Infof("[gorm-sql]:time=%v | row=%v | pos=%v |sql=%v", v[1], v[2], v[0], v[3])
		return
	}
	if len(v) == 5 {
		l.logger.Infof("[gorm-sql]:time=%v | row=%v | pos=%v | sql=%v | err=%v", v[2], v[3], v[0], v[4], v[1])
		return
	}
	l.logger.Infof("[gorm-sql]:%v", v)
}
