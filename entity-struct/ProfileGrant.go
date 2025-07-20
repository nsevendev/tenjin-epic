// ProfileGrant - Représente une autorisation accordée à un recruteur ou une entreprise pour accéder aux informations d'un candidat
// il y en aura un creer par defaut pour chaque personne, qui sera sur audience (scopeType) puis au choix du user
// si il veut tous les recruteurs ou toutes les companies ou toutes les intitutions ou tout le monde (public)
// il pourra creer plusieurs grant et donc avoir une granularité plus fine sur les champs autorisés et pour qui
type ProfileGrant struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	CandidateID   primitive.ObjectID  `bson:"candidate_id" validate:"required"` // ID du candidat
	RecruiterID   *primitive.ObjectID `bson:"recruiter_id" json:"recruiterID"`  // ID du recruteur, peut être nil si le grant est pour une entreprise ou une institution ou public
	CompanyID     *primitive.ObjectID `bson:"company_id" json:"companyID"`      // ID de l'entreprise, peut être nil si le grant est pour un recruteur ou une institution ou public
	Audience      *string             `bson:"audience" json:"audience" validate:"omitempty,oneof=recruiter institute company public"`     // cible des groupes (public cible tout le monde)
	ScopeType     string              `bson:"scope_type" json:"scopeType" validate:"required,oneof=recruiter institute company audience"` // donne le type du grant
	GrantedFields []string            `bson:"granted_fields" json:"grantedFields" validate:"required,dive,oneof=email cv"`                // Liste des champs autorisés a modifier
	Revoked       bool                `bson:"revoked" json:"revoked" validate:"required"` // Indique si le grant est révoqué ou non (n'est plus pris en compte pour les droits d'affichage)
	ExpiresAt     *time.Time          `bson:"expires_at" json:"expiresAt"`  // Date d'expiration du grant, peut être nil si le grant n'expire pas
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}
