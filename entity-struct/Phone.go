// Model représente un numéro de téléphone d'un organisme
type Model struct {
	Type   constantes.PhoneNumberType `bson:"type" json:"type" validate:"required,oneof=mobile fix fax autre"`
	Number string                     `bson:"number" json:"number" validate:"required,e164"` // Format international
	Label  constantes.LabelPhone      `bson:"label" json:"label" validate:"required,oneof=directeur secretaire reception administration autre technique support"`
}