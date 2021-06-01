package configs

import (
	"ZinX-shm/utils/logger"
	"ZinX-shm/ziface"
)

/*插件*/

//接口
type Options interface {
	Apply(config *Config) error
}

type Config struct{
	//日志管理
	Logger *logger.Logger
	//封包功能
	DataPack ziface.IDataPack
	//心跳监测功能
	Session ziface.ISession
}

//将Server端指定的Config，赋予全局Config，调用时通过全局的Config调用
func (conf *Config) Apply(config *Config) error{
	if config != conf {
		*config = *conf
	}
	return nil
}