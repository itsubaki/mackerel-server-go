package database

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itsubaki/mackerel-server-go/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GraphRepository struct {
	DB *gorm.DB
}

type GraphDef struct {
	OrgID       string `gorm:"column:org_id;       type:varchar(16); not null; primary_key"`
	Name        string `gorm:"column:name;         type:varchar(64); not null; primary_key"`
	DisplayName string `gorm:"column:display_name; type:varchar(64);"`
	Unit        string `gorm:"column:unit;         type:varchar(64);"`
	Metrics     string `gorm:"column:metrics;      type:text;"`
}

type GraphAnnotation struct {
	OrgID       string `gorm:"column:org_id;      type:varchar(16); not null;"`
	ID          string `gorm:"column:id;          type:varchar(16); not null; primary_key"`
	Title       string `gorm:"column:title;       type:varchar(64); not null;"`
	Description string `gorm:"column:description; type:varchar(64);"`
	From        int64  `gorm:"column:time_from;   type:bigint;"`
	To          int64  `gorm:"column:time_to;     type:bigint;"`
	Service     string `gorm:"column:service;     type:varchar(128); not null;"`
	Roles       string `gorm:"column:roles;       type:text;"`
}

func (a GraphAnnotation) Domain() domain.GraphAnnotation {
	return domain.GraphAnnotation{
		OrgID:       a.OrgID,
		ID:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		From:        a.From,
		To:          a.To,
		Service:     a.Service,
		Roles:       strings.Split(a.Roles, ","),
	}
}

func NewGraphRepository(handler SQLHandler) *GraphRepository {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: handler.Raw().(gorm.ConnPool),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if handler.IsDebugMode() {
		db.Logger.LogMode(4)
	}

	if err := db.AutoMigrate(&GraphDef{}); err != nil {
		panic(fmt.Errorf("auto migrate graph_def: %v", err))
	}

	if err := db.AutoMigrate(&GraphAnnotation{}); err != nil {
		panic(fmt.Errorf("auto migrate graph_annotation: %v", err))
	}

	return &GraphRepository{
		DB: db,
	}
}

func (repo *GraphRepository) SaveDef(orgID string, g []domain.GraphDef) (*domain.Success, error) {
	if err := repo.DB.Transaction(func(tx *gorm.DB) error {
		for _, r := range g {
			metrics, err := json.Marshal(r.Metrics)
			if err != nil {
				return fmt.Errorf("marshal host.Roles: %v", err)
			}

			update := GraphDef{
				DisplayName: r.DisplayName,
				Unit:        r.Unit,
				Metrics:     string(metrics),
			}

			if err := tx.Where(&GraphDef{OrgID: orgID, Name: r.Name}).Assign(&update).FirstOrCreate(&GraphDef{}).Error; err != nil {
				return fmt.Errorf("firts or create: %v", err)
			}
		}

		return nil
	}); err != nil {
		return &domain.Success{Success: false}, fmt.Errorf("transaction: %v", err)
	}

	return &domain.Success{Success: true}, nil
}

func (repo *GraphRepository) List(orgID string) (*domain.GraphAnnotations, error) {
	result := make([]GraphAnnotation, 0)
	if err := repo.DB.Where(&GraphAnnotation{OrgID: orgID}).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from graph_annotations: %v", err)
	}

	out := make([]domain.GraphAnnotation, 0)
	for _, r := range result {
		out = append(out, r.Domain())
	}

	return &domain.GraphAnnotations{GraphAnnotations: out}, nil
}

func (repo *GraphRepository) Save(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	create := GraphAnnotation{
		OrgID:       orgID,
		ID:          annotation.ID,
		Title:       annotation.Title,
		Description: annotation.Description,
		From:        annotation.From,
		To:          annotation.To,
		Service:     annotation.Service,
		Roles:       strings.Join(annotation.Roles, ","),
	}

	if err := repo.DB.Create(&create).Error; err != nil {
		return nil, fmt.Errorf("insert into graph_annotations: %v", err)
	}

	return annotation, nil
}

func (repo *GraphRepository) Update(orgID string, annotation *domain.GraphAnnotation) (*domain.GraphAnnotation, error) {
	update := GraphAnnotation{
		Title:       annotation.Title,
		Description: annotation.Description,
		From:        annotation.From,
		To:          annotation.To,
		Service:     annotation.Service,
		Roles:       strings.Join(annotation.Roles, ","),
	}

	if err := repo.DB.Where(&GraphAnnotation{OrgID: orgID, ID: annotation.ID}).Assign(&update).FirstOrCreate(&GraphAnnotation{}).Error; err != nil {
		return nil, fmt.Errorf("update graph_annotations: %v", err)
	}

	return annotation, nil
}

func (repo *GraphRepository) Delete(orgID, annotationID string) (*domain.GraphAnnotation, error) {
	result := GraphAnnotation{}
	if err := repo.DB.Where(&GraphAnnotation{OrgID: orgID, ID: annotationID}).First(&result).Error; err != nil {
		return nil, fmt.Errorf("select * from graph_annotations: %v", err)
	}

	if err := repo.DB.Delete(&GraphAnnotation{OrgID: orgID, ID: annotationID}).Error; err != nil {
		return nil, fmt.Errorf("delete from graph_annotations: %v", err)
	}

	annotation := result.Domain()
	return &annotation, nil
}
