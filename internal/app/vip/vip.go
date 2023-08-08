package vip

import (
	"time"
	"wxcloudrun-golang/internal/pkg/model"
)

type Service struct {
	VipDao   *model.Vip
	OrderDao *model.Order
}

func NewService() *Service {
	return &Service{
		VipDao:   &model.Vip{},
		OrderDao: &model.Order{},
	}
}

// GetRemainingCount 获取剩余次数
func (s *Service) GetRemainingCount(openID string) (int32, error) {
	vip, err := s.VipDao.GetByOpenID(openID)
	if err != nil {
		return 0, err
	}
	return vip.Count, nil
}

// UpdateRemainingCount 更新剩余次数
func (s *Service) UpdateRemainingCount(openID string, countToAdd int32) (*model.Vip, error) {
	vip, err := s.VipDao.UpdateCountByOpenID(openID, countToAdd)
	if err != nil {
		return nil, err
	}
	return vip, nil
}

// CreateOrder 创建订单
func (s *Service) CreateOrder(openID string, orderType int32, cost float64) (int32, error) {
	order, err := s.OrderDao.Create(&model.Order{
		OpenID:      openID,
		OrderType:   orderType,
		Cost:        cost,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	})
	if err != nil {
		return 0, err
	}
	return order.ID, nil
}

// GetOrdersByOpenID 根据订单ID获取订单
func (s *Service) GetOrdersByOpenID(openID string) ([]*model.Order, error) {
	orders, err := s.OrderDao.GetByOpenID(openID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
