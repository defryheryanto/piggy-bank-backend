package auth

type UserConfig struct {
	ConfigId         int `gorm:"primaryKey;autoIncrement;column:config_id" json:"config_id"`
	UserId           int `gorm:"column:user_id" json:"user_id"`
	MonthlyStartDate int `gorm:"column:monthly_start_date" json:"monthly_start_date"`
}

func (UserConfig) TableName() string {
	return "user_configs"
}

var DefaultUserConfig = UserConfig{
	MonthlyStartDate: 1,
}

type UserConfigRepository interface {
	Create(payload *UserConfig) error
	GetByUserId(userId int) *UserConfig
}

type UserConfigService struct {
	repository UserConfigRepository
}

func NewUserConfigService(repo UserConfigRepository) *UserConfigService {
	return &UserConfigService{repo}
}

func (s *UserConfigService) CreateDefault(userId int) error {
	payload := &DefaultUserConfig
	payload.UserId = userId

	err := s.repository.Create(payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserConfigService) GetByUserId(userId int) (*UserConfig, error) {
	cfg := s.repository.GetByUserId(userId)
	if cfg == nil {
		err := s.CreateDefault(userId)
		if err != nil {
			return nil, err
		}
		cfg = s.repository.GetByUserId(userId)
	}

	return cfg, nil
}
