package service

import (
	"context"
	"errors"
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
)

type MembershipService struct {
	membershipRepo repository.MembershipRepository
}

func NewMembershipService(membershipRepo repository.MembershipRepository) *MembershipService {
	return &MembershipService{membershipRepo: membershipRepo}
}

var planDurations = map[entity.MembershipPlan]time.Duration{
	entity.PlanFree:  30 * 24 * time.Hour,
	entity.PlanBasic: 30 * 24 * time.Hour,
	entity.PlanPro:   90 * 24 * time.Hour,
}

func (s *MembershipService) Activate(ctx context.Context, userID uint, plan entity.MembershipPlan) (*entity.Membership, error) {
	duration, ok := planDurations[plan]
	if !ok {
		return nil, errors.New("invalid membership plan")
	}

	now := time.Now()

	// Check if user already has a membership record
	existing, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Renew or upgrade: extend/reset expiry and update plan
		existing.Plan = plan
		existing.Status = entity.MembershipActive
		existing.StartsAt = now
		existing.ExpiresAt = now.Add(duration)
		existing.UpdatedAt = now
		if err := s.membershipRepo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// No existing record — create fresh
	membership := &entity.Membership{
		UserID:    userID,
		Plan:      plan,
		Status:    entity.MembershipActive,
		StartsAt:  now,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		return nil, err
	}
	return membership, nil
}

func (s *MembershipService) GetByUserID(ctx context.Context, userID uint) (*entity.Membership, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, errors.New("no membership found")
	}
	return membership, nil
}

func (s *MembershipService) Upgrade(ctx context.Context, userID uint, newPlan entity.MembershipPlan) (*entity.Membership, error) {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, errors.New("no membership found")
	}
	if !membership.IsActive() {
		return nil, errors.New("membership is not active")
	}
	if !isUpgrade(membership.Plan, newPlan) {
		return nil, errors.New("can only upgrade to a higher plan")
	}

	duration := planDurations[newPlan]
	membership.Plan = newPlan
	membership.ExpiresAt = time.Now().Add(duration)
	membership.UpdatedAt = time.Now()

	if err := s.membershipRepo.Update(ctx, membership); err != nil {
		return nil, err
	}
	return membership, nil
}

func (s *MembershipService) Cancel(ctx context.Context, userID uint) error {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if membership == nil {
		return errors.New("no membership found")
	}
	if !membership.IsActive() {
		return errors.New("membership is not active")
	}
	return s.membershipRepo.Cancel(ctx, membership.ID)
}

func (s *MembershipService) ExpireStale(ctx context.Context, userID uint) error {
	membership, err := s.membershipRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if membership == nil {
		return nil
	}
	if membership.IsExpired() && membership.Status == entity.MembershipActive {
		membership.Status = entity.MembershipExpired
		return s.membershipRepo.Update(ctx, membership)
	}
	return nil
}

func isUpgrade(current, next entity.MembershipPlan) bool {
	rank := map[entity.MembershipPlan]int{
		entity.PlanFree:  0,
		entity.PlanBasic: 1,
		entity.PlanPro:   2,
	}
	return rank[next] > rank[current]
}
