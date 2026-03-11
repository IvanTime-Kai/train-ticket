package setting

type Settings struct {
	Logger LoggerSetting `mapstructure:"logger"`
	MySql  MySqlSetting  `mapstructure:"mysql"`
	Server ServerSetting `mapstructure:"server"`
	Redis  RedisSetting  `mapstructure:"redis"`
	GRPC   GRPCSetting   `mapstructure:"grpc"`
}

type LoggerSetting struct {
	LogLevel     string `mapstructure:"log_level"`
	FileLogName string `mapstructure:"file_log_name"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
}

type MySqlSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DbName          string `mapstructure:"dbName"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifeTime"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type GRPCSetting struct {
	TrainService string `mapstructure:"train_service"`
}
