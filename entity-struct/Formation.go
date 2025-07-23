// Model représente une formation exercée par un organisme de formation
type Formation struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title           string               `bson:"title" json:"title" validate:"required,min=2,max=200"`
	Description     string               `bson:"description" json:"description" validate:"required,min=10,max=2000"`
	IsActive        bool                 `bson:"is_active" json:"isActive"`
	Duration        int                  `bson:"duration" json:"duration" validate:"required,min=1"` // en heures
	MaxParticipants *int                 `bson:"max_participants" json:"maxParticipants"`
	DocumentUrls    []string             `bson:"document_urls" json:"documentUrls" validate:"omitempty,dive,url"`
	InstituteID     primitive.ObjectID   `bson:"institute_id" json:"instituteId" validate:"required"`
	SessionIDs      []primitive.ObjectID `bson:"sessions" json:"sessionIds"`
	CompetenceIDs   []primitive.ObjectID `bson:"competence_ids" json:"competenceIds" validate:"required,dive"`
	JobID           primitive.ObjectID   `bson:"external_job_ref" json:"jobId" validate:"required"`
	CreatedAt       time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updatedAt"`
}
