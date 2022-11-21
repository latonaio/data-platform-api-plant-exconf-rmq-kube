package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-plant-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-plant-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"sync"

	database "github.com/latonaio/golang-mysql-network-connector"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.PlantGeneral {
	businessPartner := *input.PlantGeneral.BusinessPartner
	plant := *input.PlantGeneral.Plant
	notKeyExistence := make([]dpfm_api_output_formatter.PlantGeneral, 0, 1)
	KeyExistence := make([]dpfm_api_output_formatter.PlantGeneral, 0, 1)

	existData := &dpfm_api_output_formatter.PlantGeneral{
		BusinessPartner: businessPartner,
		Plant:           plant,
		ExistenceConf:   false,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !e.confPlantGeneral(businessPartner, plant) {
			notKeyExistence = append(
				notKeyExistence,
				dpfm_api_output_formatter.PlantGeneral{businessPartner, plant, false},
			)
			return
		}
		KeyExistence = append(KeyExistence, dpfm_api_output_formatter.PlantGeneral{businessPartner, plant, true})
	}()

	wg.Wait()

	if len(KeyExistence) == 0 {
		return existData
	}
	if len(notKeyExistence) > 0 {
		return existData
	}

	existData.ExistenceConf = true
	return existData
}

func (e *ExistenceConf) confPlantGeneral(businessPartner int, plant string) bool {
	rows, err := e.db.Query(
		`SELECT Plant 
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_plant_general_data 
		WHERE (BusinessPartner, Plant) = (?, ?);`, businessPartner, plant,
	)
	if err != nil {
		e.l.Error(err)
		return false
	}
	if err != nil {
		e.l.Error(err)
		return false
	}

	for rows.Next() {
		var businessPartner int
		var plant string
		err := rows.Scan(&plant)
		if err != nil {
			e.l.Error(err)
			continue
		}
		if businessPartner == businessPartner {
			return true
		}
	}
	return false
}
