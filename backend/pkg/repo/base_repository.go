package repo

import (
	"backend/pkg/db/model"
	"backend/pkg/utils"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repo[T model.BaseInterface] struct {
	DB *gorm.DB
}

func (r Repo[T]) Create(ctx context.Context, t *T) *gorm.DB {
	tx := r.DB.Begin()
	if tx.Error != nil {
		utils.Handle(tx.Error)
		return tx
	}

	result := tx.WithContext(ctx).Create(&t)
	if result.Error != nil {
		tx.Rollback()
		return result
	}

	if err := tx.Commit().Error; err != nil {
		utils.Handle(err)
		return tx
	}

	return result
}

func (r Repo[T]) Patch(ctx context.Context, t *T, id string, patch any) *gorm.DB {
	tx := r.DB.Begin()
	if tx.Error != nil {
		utils.Handle(tx.Error)
		return tx
	}

	result := tx.WithContext(ctx).Model(t).Clauses(clause.Returning{}).Where("id = ?", id).Updates(patch)
	if result.Error != nil {
		tx.Rollback()
		return result
	}

	if err := tx.Commit().Error; err != nil {
		utils.Handle(err)
		return tx
	}

	return result
}

func (r Repo[T]) Delete(ctx context.Context, t *T, id string) *gorm.DB {
	tx := r.DB.Begin()
	if tx.Error != nil {
		utils.Handle(tx.Error)
		return tx
	}

	patch := map[string]interface{}{
		"is_deleted": true,
		"is_active":  false,
		"deleted_at": time.Now().UTC(),
	}

	result := tx.WithContext(ctx).Model(t).Clauses(clause.Returning{}).Where("id = ?", id).Updates(patch)
	if result.Error != nil {
		tx.Rollback()
		return result
	}

	if err := tx.Commit().Error; err != nil {
		utils.Handle(err)
		return tx
	}

	return result
}

func (r Repo[T]) GetAll(ctx context.Context) ([]model.HolderEntity, error) {
	var holders []model.HolderEntity
	err := r.DB.WithContext(ctx).Find(&holders).Error
	return holders, err
}

func (r Repo[T]) DeleteAll(ctx context.Context) error {
	return r.DB.WithContext(ctx).Where("1 = 1").Delete(&model.HolderEntity{}).Error
}
