package clickhousestorage

import (
	"AlexSarva/tender/utils/dbutils"
	"context"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouse struct {
	Database driver.Conn
	ctx      context.Context
}

func MyClickHouseDB(path string) *ClickHouse {
	parsedCfg, parseCfgErr := dbutils.ParseConfigDB(path)
	if parseCfgErr != nil {
		log.Println("Проблема с конфиогом для ClickHouse")
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: parsedCfg.DatabaseName,
			Username: parsedCfg.User,
			Password: parsedCfg.Password,
		},
		//Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    50,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour * 24,
	})
	if err != nil {
		log.Println("НЕТ подключения к ClickHouse: ", err)
	}
	return &ClickHouse{
		Database: conn,
		ctx:      context.Background(),
	}
}

func (c *ClickHouse) Ping() bool {
	pingErr := c.Database.Ping(c.ctx)
	if pingErr != nil {
		log.Println("НЕ пингуется", pingErr)
		return false
	}
	return true
}
