package auth

type UserConfig struct {
	ConfigId         int `gorm:"primaryKey;autoIncrement;column:config_id" json:"-"`
	UserId           int `gorm:"column:user_id" json:"-"`
	MonthlyStartDate int `gorm:"column:monthly_start_date" json:"monthly_start_date"`
}

func (UserConfig) TableName() string {
	return "user_configs"
}

type UpdateUserConfigPayload struct {
	UserId           int `json:"user_id"`
	MonthlyStartDate int `json:"monthly_start_date"`
}

type UserConfigRepository interface {
	Create(payload *UserConfig) error
	Update(payload *UserConfig) error
	GetByUserId(userId int) *UserConfig
}

func DefaultUserConfig() *UserConfig {
	return &UserConfig{
		MonthlyStartDate: 1,
	}
}

type UserConfigService struct {
	repository UserConfigRepository
}

func NewUserConfigService(repo UserConfigRepository) *UserConfigService {
	return &UserConfigService{repo}
}

func (s *UserConfigService) CreateDefault(userId int) error {
	payload := DefaultUserConfig()
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

func (s *UserConfigService) Update(payload *UpdateUserConfigPayload) error {
	cfg := s.repository.GetByUserId(payload.UserId)
	if cfg == nil {
		defaultCfg := DefaultUserConfig()
		populateUpdateConfig(defaultCfg, payload)
		err := s.repository.Create(defaultCfg)
		if err != nil {
			return err
		}
		return nil
	}

	populateUpdateConfig(cfg, payload)
	err := s.repository.Update(cfg)
	if err != nil {
		return err
	}

	return nil
}

func populateUpdateConfig(current *UserConfig, payload *UpdateUserConfigPayload) {
	current.UserId = payload.UserId
	if payload.MonthlyStartDate != 0 {
		current.MonthlyStartDate = payload.MonthlyStartDate
	}
}
