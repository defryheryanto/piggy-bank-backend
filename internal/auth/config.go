package auth

type UserConfig struct {
	ConfigId         int `gorm:"primaryKey;autoIncrement;column:config_id" json:"config_id"`
	UserId           int `gorm:"column:user_id" json:"user_id"`
	MonthlyStartDate int `gorm:"column:monthly_start_date" json:"monthly_start_date"`
}

func (UserConfig) TableName() string {
	return "user_configs"
}

type UserConfigRepository interface {
	Create(payload *UserConfig) error
}

type UserConfigService struct {
	repository UserConfigRepository
}

func NewUserConfigService(repo UserConfigRepository) *UserConfigService {
	return &UserConfigService{repo}
}

func (s *UserConfigService) CreateDefault(userId int) error {
	payload := &UserConfig{
		UserId:           userId,
		MonthlyStartDate: 1,
	}

	err := s.repository.Create(payload)
	if err != nil {
		return err
	}

	return nil
}
