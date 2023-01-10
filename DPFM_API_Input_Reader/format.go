package dpfm_api_input_reader

import (
	"data-platform-api-plant-exconf-rmq-kube/DPFM_API_Caller/requests"
)

func (sdc *SDC) ConvertToPlant() *requests.PlantGeneral {
	data := sdc.PlantGeneral
	return &requests.PlantGeneral{
		BusinessPartner: data.BusinessPartner,
		Plant:           data.Plant,
	}
}
