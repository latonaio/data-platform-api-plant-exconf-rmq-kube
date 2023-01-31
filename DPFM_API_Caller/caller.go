package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-plant-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-plant-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"encoding/json"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
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

func (e *ExistenceConf) Conf(msg rabbitmq.RabbitmqMessage) interface{} {
	var ret interface{}
	ret = map[string]interface{}{
		"ExistenceConf": false,
	}
	input := make(map[string]interface{})
	err := json.Unmarshal(msg.Raw(), &input)
	if err != nil {
		return ret
	}

	_, ok := input["PlantGeneral"]
	if ok {
		input := &dpfm_api_input_reader.SDC{}
		err = json.Unmarshal(msg.Raw(), input)
		ret = e.confPlantGeneral(input)
		goto endProcess
	}

	err = xerrors.Errorf("can not get exconf check target")
endProcess:
	if err != nil {
		e.l.Error(err)
	}
	return ret
}

func (e *ExistenceConf) confPlantGeneral(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.PlantGeneral {
	exconf := dpfm_api_output_formatter.PlantGeneral{
		ExistenceConf: false,
	}
	if input.PlantGeneral.Plant == nil {
		return &exconf
	}
	if input.PlantGeneral.BusinessPartner == nil {
		return &exconf
	}
	exconf = dpfm_api_output_formatter.PlantGeneral{
		Plant:           *input.PlantGeneral.Plant,
		BusinessPartner: *input.PlantGeneral.BusinessPartner,
		ExistenceConf:   false,
	}

	rows, err := e.db.Query(
		`SELECT Plant, BusinessPartner
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_plant_general_data 
		WHERE (Plant, BusinessPartner) = (?, ?);`, exconf.Plant, exconf.BusinessPartner,
	)
	if err != nil {
		e.l.Error(err)
		return &exconf
	}
	defer rows.Close()

	exconf.ExistenceConf = rows.Next()
	return &exconf
}
