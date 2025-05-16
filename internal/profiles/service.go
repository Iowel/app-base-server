package profiles

type ProfileService struct {
	profileRepo *ProfileRepository
}

func NewProfileService(profile *ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profile,
	}
}

// func (p *ProfileService) GetUserStatus(ctx context.Context, userID int64) (string, error) {
// 	status, err := p.profileRepo.GetUserStatus(ctx, userID)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get user status for user %d: %w", userID, err)
// 	}

// 	return status, nil

// }
