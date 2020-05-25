package controller

import (
	"github.com/itsubaki/mackerel-api/pkg/interface/database"
)

type DowntimeController struct {
}

func NewDowntimeController(handler database.SQLHandler) *DowntimeController {
	return &DowntimeController{}
}

func (s *DowntimeController) List(c Context) {

}

func (s *DowntimeController) Save(c Context) {

}

func (s *DowntimeController) Update(c Context) {

}

func (s DowntimeController) Downtime(c Context) {

}

func (s DowntimeController) Delete(c Context) {

}
