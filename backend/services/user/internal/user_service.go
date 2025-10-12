package internal

import (
	"context"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/user/models"
	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/user/proto"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.UserProfile, error) {
	user, err := s.repo.GetProfile(req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.UserProfile{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Language: user.Language,
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest) (*proto.UserProfile, error) {
	user := &models.User{
		ID:       req.UserId,
		Name:     req.Name,
		Email:    req.Email,
		Language: req.Language,
	}

	if err := s.repo.UpdateProfile(user); err != nil {
		return nil, err
	}

	return &proto.UserProfile{
		Id:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Language: user.Language,
	}, nil
}

func (s *UserService) GetBookingHistory(ctx context.Context, req *proto.GetBookingHistoryRequest) (*proto.BookingHistoryResponse, error) {
	bookings, err := s.repo.GetBookingHistory(req.UserId)
	if err != nil {
		return nil, err
	}

	var protoBookings []*proto.BookingItem
	for _, b := range bookings {
		protoBookings = append(protoBookings, &proto.BookingItem{
			BookingId: b.ID,
			RoomName:  b.RoomName,
			StartTime: b.StartTime,
			EndTime:   b.EndTime,
			Status:    b.Status,
		})
	}

	return &proto.BookingHistoryResponse{Bookings: protoBookings}, nil
}

func (s *UserService) GetPreferences(ctx context.Context, req *proto.GetPreferencesRequest) (*proto.UserPreferences, error) {
	p, err := s.repo.GetPreferences(req.UserId)
	if err != nil {
		return nil, err
	}

	return &proto.UserPreferences{
		NotificationType: p.NotificationType,
		Language:         p.Language,
	}, nil
}

func (s *UserService) UpdatePreferences(ctx context.Context, req *proto.UpdatePreferencesRequest) (*proto.UserPreferences, error) {
	p := &models.Preferences{
		UserID:           req.UserId,
		NotificationType: req.NotificationType,
		Language:         req.Language,
	}

	if err := s.repo.UpdatePreferences(p); err != nil {
		return nil, err
	}

	return &proto.UserPreferences{
		NotificationType: p.NotificationType,
		Language:         p.Language,
	}, nil
}
