package setting

type Settings struct {
	JWT    JWTSetting    `mapstructure:"jwt"`
	Logger LoggerSetting `mapstructure:"logger"`
	MySql  MySqlSetting  `mapstructure:"mysql"`
	Server ServerSetting `mapstructure:"server"`
	Redis  RedisSetting  `mapstructure:"redis"`
}

type JWTSetting struct {
	TOKEN_HOUR_LIFESPAN int    `mapstructure:"TOKEN_HOUR_LIFESPAN"`
	JWT_EXPIRATION      string `mapstructure:"JWT_EXPIRATION"`
	API_SECRET          string `mapstructure:"API_SECRET"`
	ACCESS_TOKEN_TTL    int    `mapstructure:"ACCESS_TOKEN_TTL"`
	REFRESH_TOKEN_TTL   int    `mapstructure:"REFRESH_TOKEN_TTL"`
}

type LoggerSetting struct {
	Log_level     string `mapstructure:"log_level"`
	File_log_name string `mapstructure:"file_log_name"`
	Max_size      int    `mapstructure:"max_size"`
	Max_backups   int    `mapstructure:"max_backups"`
	Max_age       int    `mapstructure:"max_age"`
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
