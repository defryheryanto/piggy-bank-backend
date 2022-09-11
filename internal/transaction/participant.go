package transaction

type Participant struct {
	ParticipantId int     `gorm:"primaryKey;autoIncrement;column:participant_id" json:"participant_id"`
	TransactionId int     `gorm:"column:transaction_id" json:"transaction_id"`
	Name          string  `gorm:"column:name" json:"name"`
	Amount        float64 `gorm:"column:amount" json:"amount"`
}

func (Participant) TableName() string {
	return "transaction_participants"
}

type CreateParticipantPayload struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

func (p *CreateParticipantPayload) Validate() error {
	if p.Name == "" {
		return ErrEmptyParticipantName
	}
	if p.Amount < 0 {
		return ErrInvalidAmount
	}

	return nil
}

type ParticipantRepository interface {
	BulkCreate(payloads []*Participant) error
}

type ParticipantService struct {
	repository ParticipantRepository
}

func NewParticipantService(repo ParticipantRepository) *ParticipantService {
	return &ParticipantService{repo}
}

func (s *ParticipantService) BulkCreate(transactionId int, payloads []*CreateParticipantPayload) error {
	participants := []*Participant{}

	for _, payload := range payloads {
		err := payload.Validate()
		if err != nil {
			return err
		}
		participants = append(participants, &Participant{
			TransactionId: transactionId,
			Name:          payload.Name,
			Amount:        payload.Amount,
		})
	}

	err := s.repository.BulkCreate(participants)
	if err != nil {
		return err
	}

	return nil
}
