package repositories

import (
    "context"
    "errors"

    "github.com/cureerel/cserver/internal/domain/entity"
    "github.com/cureerel/cserver/internal/domain/repository"
    "github.com/cureerel/cserver/internal/infrastructure/postgres/models"
    "gorm.io/gorm"
)

type membershipRepository struct {
    db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) repository.MembershipRepository {
    return &membershipRepository{db: db}
}

func (r *membershipRepository) Create(ctx context.Context, membership *entity.Membership) error {
    m := models.MembershipFromDomain(membership)
    if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
        return err
    }
    membership.ID = m.ID
    return nil
}

func (r *membershipRepository) GetByUserID(ctx context.Context, userID uint) (*entity.Membership, error) {
    var m models.Membership
    if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return m.ToDomain(), nil
}

func (r *membershipRepository) Update(ctx context.Context, membership *entity.Membership) error {
    m := models.MembershipFromDomain(membership)
    return r.db.WithContext(ctx).Save(m).Error
}

func (r *membershipRepository) Cancel(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).
        Model(&models.Membership{}).
        Where("id = ?", id).
        Update("status", string(entity.MembershipCancelled)).Error
}