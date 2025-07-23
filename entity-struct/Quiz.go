// Question d'un quizz
type Question struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Question string   `bson:"question" json:"question" validate:"required"`                // LA question elle-même
	Type     string   `bson:"type" json:"type" validate:"required,oneof=qcm text file"`    // quel est le genre de la question : il peut etre qcm, text ou file ceci determine le type de reponse
	Choices  []string `bson:"choices" json:"choices"`                                      // plusieurs choix possibles, si type est qcm alors il faut au moins 2 choix, si le type est text et file alors il n'y a pas de choix
	Answer   *string  `bson:"answer" json:"answer"`                                        // reponse du qcm ou du text, si type est file alors pas de good answer, tester le type avant enregistrement
	Point    int  `bson:"point" json:"point" validate:"min=1"`                         // point attribuer à la question
	CreatedAt time.Time  `bson:"created_at" json:"createdAt"`
    UpdatedAt time.Time  `bson:"updated_at" json:"updatedAt"`
}

// Quiz represents the main quiz entity
type Quiz struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	RecruiterID *primitive.ObjectID `bson:"recruiter_id" json:"recruiterId"`              // personne qui creer le quiz si type = offer
	TeacherID   *primitive.ObjectID `bson:"teacher_id" json:"teacherId"`                  // personne qui creer le quiz si type = course
	OfferID     *primitive.ObjectID `bson:"offer_id" json:"offerId"`                      // si recruiterId et type offer alors offerId
	CourseID    *primitive.ObjectID `bson:"course_id" json:"courseId"`                    // si teacherId et type course alors courseId
	Title       string   `bson:"title" json:"title" validate:"required"`                     // titre du quiz
	Description string   `bson:"description" json:"description" validate:"required"`         // description du quiz
	Type        string   `bson:"type" json:"type" validate:"required,oneof=offer course"`    // type du quiz (offer ou course)
	Questions   []Question  `bson:"questions" json:"questions" validate:"required,dive"`        // list des questions du quiz
	TimeLimit   *int        `bson:"time_limit" json:"timeLimit"`                                // Temps limite en secondes pour le quiz (c'est le temps total pour le quiz, pas par question)
	IsRequired  bool        `bson:"is_required" json:"isRequired"`                              // definit si le quiz est obligatoire pour passer à une étape suivante ou autre (true par défaut pour offer)
	IsActive    bool        `bson:"is_active" json:"isActive" validate:"required"`              // quiz est-il actif ? (false par défaut) (defini si visible ou pas)
	PassingScore *int       `bson:"passing_score" json:"passingScore" validate:"required"`                           // score minimal pour réussir le quiz, calculé sur le total des points des questions
	VisibilityPassingScore bool `bson:"visibility_passing_score" json:"visibilityPassingScore" validate:"required"`       // si le score de passage est visible par les utilisateurs (true par défaut)
	CreatedAt   time.Time       `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time       `bson:"updated_at" json:"updatedAt"`
}

// QuizSession le user (candidate ou student) est en train de faire le quiz
type QuizSession struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuizID          primitive.ObjectID `bson:"quiz_id" json:"quizId" validate:"required"`
	UserID          primitive.ObjectID `bson:"user_id" json:"userId" validate:"required"`
	Type            string             `bson:"type" json:"type" validate:"required,oneof=offer course"`   // le type de la session en fonction du type du quiz
	Answers         []QuizAnswer       `bson:"answers" json:"Answers"`           // réponses de l'utilisateur, peut être vide si le quiz n'est pas commencé
	StartedAt       *time.Time         `bson:"started_at" json:"startedAt" validate:"required"`      // date de début de la session
	FinishedAt      *time.Time         `bson:"finished_at" json:"finishedAt"`    // date de fin de la session, peut être vide si le quiz n'est pas terminé
	TimeWorking     *int               `bson:"time_working" json:"timeWorking" validate:"required"`  // temps passé en secondes sur le quiz (des que la session est en cours (inprogress) chaque seconde viens s'ajouter)
	// mise a jour à dans la bdd à chaque passage à pause, completed (enregistrement auto en bdd toutes les 30s, si depassement du temps du quiz alors pas d'enregistrement auto)
	// enregistrement dans le local storage en temps reel avec objet "timeWorking" => id quzSession et time en seconde (attention stockage sous forme de string, donc stringifier l'objet)
	TimeRemaining   *int               `bson:"time_remaining" json:"timeRemaining" validate:"required"` // temps restant en seconde, carlculé à partir de la durée totale du quiz moins le temps déjà travaillé
	// mise a jour à dans la bdd à chaque passage à pause, completed (enregistrement auto en bdd toutes les 30s, si depassement du temps du quiz alors pas d'enregistrement auto)
	// enregistrement dans le local storage en temps reel avec objet "timeRemaining" => id quzSession et time en seconde (attention stockage sous forme de string, donc stringifier l'objet)
	Status          string             `bson:"status" json:"status" validate:"required,oneof=inprogress started paused expired completed"` // état de la session (started, paused, expired, completed)
	// meme si le quizz est expired, la personne peut le terminer, et le soumettre, le receveur verra qu'il n'a pas fini dans le temps imparti
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}

// QuizAnswer represente les reponses au question d'un quiz dans une session
// sont donnée au quizsession
type QuizAnswer struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuestionID string `bson:"question_id" json:"questionId" validate:"required"`     // reference de la question
	Answer     *string `bson:"answer" json:"answer"`                                 // la réponse du user, obligatoire si type de question est text ou qcm, sinon pas obligatoire
	FileURL    *string `bson:"file_url" json:"fileUrl"`                              // si type de la question est file, l'url du fichier uploadé
}

// QuizSubmission represents the final submission of a completed quiz
type QuizSubmission struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	QuizSessionID primitive.ObjectID  `bson:"quiz_session_id" json:"quizSessionId" validate:"required"`   // id de la session de quiz de laquelle est creer cette soumission
	QuizID        primitive.ObjectID  `bson:"quiz_id" json:"quizId" validate:"required"`   // id du quiz auquel cette soumission est liée
	UserID        primitive.ObjectID  `bson:"user_id" json:"userId" validate:"required"`   // id de l'utilisateur qui a soumis le quiz
	OfferID       *primitive.ObjectID `bson:"offer_id" json:"offerId"`      // id de l'offre si lié à une offre en fonction du type du quiz
	CourseID      *primitive.ObjectID `bson:"course_id" json:"courseId"`    // id du cours si lié à un cours en fonction du type du quiz
	FinalAnswers  []QuizAnswer        `bson:"final_answers" json:"finalAnswers" validate:"required,dive"`   // réponses finales de l'utilisateur, viens de la session de quiz
	SubmittedAt   time.Time           `bson:"submitted_at" json:"submittedAt" validate:"required"`                  // date de soumission du quiz
	Score         *int                `bson:"score" json:"score" validate:"required"`           // score obtenu par l'utilisateur depuis la session de quiz
	Percentage    *float64            `bson:"percentage" json:"percentage" validate:"required"` // calcul du pourcentage de réussite en fonction du total score du quiz et le résultat de la session de quiz
	Passed        *bool               `bson:"passed" json:"passed" validate:"required"`         // au dela de 50%% true sinon false
	Status        string              `bson:"status" json:"status" validate:"required,oneof=submitted evaluated"` // deux etats submitted toute les questions sont automatiquement évaluées (qcm, text),
	// evaluated lorsque qu'il reste des questions file et que le teacher doit les évaluer manuellement
	ReviewedBy    primitive.ObjectID `bson:"reviewed_by" json:"reviewedBy" validate:"required"` // la personne qui a créé le quiz
	ReviewedAt    *time.Time          `bson:"reviewed_at" json:"reviewedAt"`                    // date quand la personne a commencer à évaluer les questions file
	Feedback      *string             `bson:"feedback" json:"feedback"`     // message de la personne qui a créé le quiz (retour sur le quiz, sur les réponses, etc.)
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}

// QuizStats represents aggregated statistics for a quiz
type QuizStats struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuizID            primitive.ObjectID `bson:"quiz_id" json:"quizId" validate:"required"` // id du quiz pour lequel les stats sont calculées
	TotalSubmissions  int                `bson:"total_submissions" json:"totalSubmissions" validate:"required"` // nombre total de soumissions pour ce quiz
	TotalPassed       int                `bson:"total_passed" json:"totalPassed" validate:"required"` // nombre total de soumissions réussies
	TotalFailed       int                `bson:"total_failed" json:"totalFailed" validate:"required"` // nombre total de soumissions échouées
	AverageScore      *float64           `bson:"average_score" json:"averageScore" validate:"required"` // score moyen des soumissions, calculé sur les scores de toutes les soumissions
	AveragePercentage *float64           `bson:"average_percentage" json:"averagePercentage" validate:"required"` // pourcentage moyen de réussite, calculé sur les pourcentages de toutes les soumissions
	AverageTimeSpent  *int               `bson:"average_time_spent" json:"averageTimeSpent" validate:"required"`  // temps moyen passé sur le quiz en secondes, calculé sur le temps total passé par tous les utilisateurs
	CreatedAt         time.Time          `bson:"created_at" json:"createdAt"`
    UpdatedAt         time.Time          `bson:"updated_at" json:"updatedAt"`
}
