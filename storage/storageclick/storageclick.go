package clickhousestorage

import (
	"AlexSarva/tender/models"
	"AlexSarva/tender/utils/dbutils"
	"context"
	"fmt"
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

func (c *ClickHouse) GetOrgInfo(inn string) (models.Organization, error) {
	log.Printf("%T", inn)
	log.Printf("%v", inn)
	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 1000,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	var orgInfo []models.Organization
	err := c.Database.Select(ctx, &orgInfo, `
select ogrn, inn, kpp,
       replaceAll(case when short_name = '' then full_name else short_name end,'&quot;','"') short_name,
       replaceAll(full_name,'&quot;','"') full_name, reg_date, end_date, okved_id, capital,
       cast(region as Int8) region_id,
       replaceAll(case when area = '' then '' else area end ||
       case when city = '' then '' when area = '' then city else ', '||city end ||
       case when settlement = '' then '' when city = '' then settlement else ', '||settlement end ||
       case when street = '' then '' when settlement ='' then street else ', '||street end ||
       case when house = '' then '' else ', '||house end ||
       case when corpus in ('','-') then '' else ', '||corpus end ||
       case when apartment in ('','-') then '' else ', пом. '||apartment end,'&quot;','"')  address
       from reestr_company.org_full
        where inn = $1
`, inn)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v\n", orgInfo)
	return orgInfo[0], nil
}
