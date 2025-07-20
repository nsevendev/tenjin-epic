// Offer - Template/modèle d'offre créé par le recruteur
type Offer struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	RecruiterID   primitive.ObjectID  `bson:"recruiter_id" json:"recruiterId" binding:"required" validate:"required"`  // ID du recruteur qui a créé l'offre
	CompanyID     primitive.ObjectID  `bson:"company_id" json:"companyId" binding:"required" validate:"required"`  // ID de l'entreprise pour laquelle l'offre est créée
	Title         string              `bson:"title" json:"title" binding:"required" validate:"required"`  // Titre de l'offre
	Description   string              `bson:"description" json:"description" binding:"required" validate:"required"`  // Description détaillée de l'offre
	AttachmentURL *string             `bson:"attachment_url" json:"attachmentUrl"`  // URL de l'attachement (par exemple, un document PDF ou un lien vers une page web)
	QuizID        *primitive.ObjectID `bson:"quiz_id" json:"quizId"`  // ID du quiz associé à l'offre, peut être nil si aucun quiz n'est requis
	QuizRequired  *bool               `bson:"quiz_required" json:"quizRequired" validate:"required"`  // Indique si le quiz est requis pour postuler à l'offre (defaut à false)
	StartDateJob  time.Time           `bson:"start_date_job" json:"startDateJob" binding:"required" validate:"required"`  // Date de début du travail pour l'offre
	Salary        string              `bson:"salary" json:"salary" binding:"required" validate:"required"`  // Salaire proposé pour l'offre, peut être un montant fixe ou une fourchette
	ExpiredAt     time.Time           `bson:"expired_at" json:"expiredAt" binding:"required" validate:"required"`  // Date d'expiration de l'offre, après laquelle elle n'est plus valide
	EndDate       time.Time           `bson:"end_date" json:"endDate" binding:"required" validate:"required"`  // Date de fin de l'offre, peut être la même que la date d'expiration ou différente
	Status        string              `bson:"status" json:"status" binding:"required" validate:"required,oneof=enabled expired disabled archived"`  // Statut de l'offre (enabled, expired, disabled, archived)
	EmploiType    string              `bson:"emploi_type" json:"emploiType" binding:"required" validate:"required,oneof=CDI CDD Alternance Stage Freelance"`  // Type d'emploi (CDI, CDD, Alternance, Stage, Freelance)
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}

// OfferSent - Table intermédiaire : représente l'envoi d'une offre à un candidat spécifique
type OfferSent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OfferID     primitive.ObjectID `bson:"offer_id" json:"offerId" binding:"required" validate:"required"`  // ID de l'offre envoyée
	CandidateID primitive.ObjectID `bson:"candidate_id" json:"candidateId" binding:"required" validate:"required"`  // ID du candidat à qui l'offre est envoyée
	RecruiterID primitive.ObjectID `bson:"recruiter_id" json:"recruiterId" binding:"required" validate:"required"`  // ID du recruteur qui a envoyé l'offre
	CompanyID   primitive.ObjectID `bson:"company_id" json:"companyId" binding:"required" validate:"required"`  // ID de l'entreprise pour laquelle l'offre est envoyée
	Status      string             `bson:"status" json:"status" binding:"required" validate:"required,oneof=sent viewed responded"`  // Statut de l'envoi de l'offre (sent, viewed, responded)
	SentAt      time.Time          `bson:"sent_at" json:"sentAt" binding:"required" validate:"required"`  // Date d'envoi de l'offre au candidat
	ViewedAt    *time.Time         `bson:"viewed_at" json:"viewedAt"`  // Date à laquelle le candidat a consulté l'offre, peut être nil si non consultée
	RespondedAt *time.Time         `bson:"responded_at" json:"respondedAt"`  // Date à laquelle le candidat a répondu à l'offre, peut être nil si non répondu
	Message     *string            `bson:"message" json:"message"`  // Message personnalisé envoyé avec l'offre, peut être nil si aucun message n'est envoyé
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// OfferResponse - Réponse du candidat à une offre spécifique
type OfferResponse struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OfferSentID  primitive.ObjectID `bson:"offer_sent_id" json:"offerSentId" binding:"required" validate:"required"`  // ID de l'envoi de l'offre auquel le candidat répond
	OfferID      primitive.ObjectID `bson:"offer_id" json:"offerId" binding:"required" validate:"required"`  // ID de l'offre à laquelle le candidat répond
	CompanyID    primitive.ObjectID `bson:"company_id" json:"companyId" binding:"required" validate:"required"`  // ID de l'entreprise pour laquelle l'offre est envoyée
	CandidateID  primitive.ObjectID `bson:"candidate_id" json:"candidateId" binding:"required" validate:"required"`  // ID du candidat qui répond à l'offre
	RecruiterID  primitive.ObjectID `bson:"recruiter_id" json:"recruiterId" binding:"required" validate:"required"`  // ID du recruteur qui a envoyé l'offre
	Status       string             `bson:"status" json:"status" validate:"required,oneof=accepted declined"`  // Statut de la réponse du candidat (accepted, declined)
	SharedFields []string           `bson:"shared_fields" json:"sharedFields" validate:"dive,oneof=email phone cv linkedin github skills experience location identity"` // a modifer
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
}
