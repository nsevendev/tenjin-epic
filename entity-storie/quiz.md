fait des # Quiz

> Qui a acc�s aux quiz ?
> - Un recruteur authentifi� pour les quiz d'offres qu'il a cr��s
> - Un teacher/formateur authentifi� pour les quiz de cours qu'il a cr��s 
> - Un candidat/stagiaire authentifi� pour r�pondre aux quiz qui lui sont destin�s
> - Un admin pour tous les quiz

## Role "recruiter" avec SubRole "recruiter_admin"

- Champs autoris�s :
  - tous les champs de la struct Quiz (type "offer" uniquement)
  - tous les champs de QuizSession, QuizSubmission et QuizStats
- Lecture :
  - peut voir tous les quiz de type "offer"
  - peut voir toutes les sessions et soumissions des quiz d'offres
  - peut voir tous les d�tails d'un quiz d'offre
- �criture :
  - peut cr�er des quiz de type "offer"
  - peut modifier tous les quiz de type "offer"
  - peut activer/d�sactiver tous les quiz d'offres
  - peut supprimer tous les quiz d'offres
  - peut �valuer les soumissions de quiz (avec feedback)

## Role "recruiter" avec SubRole "recruiter_user"

- Champs autoris�s :
  - tous les champs de la struct Quiz (type "offer" uniquement)
  - tous les champs de QuizSession, QuizSubmission et QuizStats
- Lecture :
  - peut voir seulement ses quiz de type "offer"
  - peut voir seulement les sessions et soumissions de ses quiz
  - peut voir tous les d�tails de ses quiz
- �criture :
  - peut cr�er des quiz de type "offer" seulement pour ses offres
  - peut modifier seulement ses quiz
  - peut activer/d�sactiver seulement ses quiz
  - peut �valuer les soumissions de ses quiz (avec feedback)

## Role "teacher" 

- Champs autoris�s :
  - tous les champs de la struct Quiz (type "course" uniquement)
  - tous les champs de QuizSession, QuizSubmission et QuizStats
- Lecture :
  - peut voir seulement ses quiz de type "course"
  - peut voir les sessions et soumissions de ses quiz de cours
  - peut voir tous les d�tails de ses quiz
- �criture :
  - peut cr�er des quiz de type "course" pour ses cours
  - peut modifier seulement ses quiz
  - peut rendre obligatoire/optionnel ses quiz (IsRequired)
  - peut activer/d�sactiver seulement ses quiz
  - peut �valuer les soumissions de ses quiz (avec feedback)

## Role "candidate" et "stagiaire"

- Champs autoris�s :
  - tous les champs sauf les r�ponses correctes (Answer masqu� dans Question)
  - ses propres QuizSession et QuizSubmission uniquement
- Lecture :
  - peut voir les quiz li�s aux cours auxquels il participe
  - peut voir ses propres sessions et soumissions
  - peut voir ses scores et feedback re�us
- �criture :
  - peut d�marrer un quiz de cours (cr�ation QuizSession)
  - peut soumettre ses r�ponses au quiz
  - ne peut pas modifier une soumission une fois finalis�e

> La liste des quiz s'affiche sous forme de tableau, avec les colonnes suivantes :
> - Titre
> - Type (offer/course)
> - Statut (active/inactive) - modifiable directement avec un toggle
> - Nombre de questions
> - Temps limite (si d�fini)
> - Score minimum requis (si d�fini)
> - Date de cr�ation
> - Statistiques (nombre de soumissions, taux de r�ussite)
> - Actions (modifier, voir stats, supprimer)

> Les d�tails d'un quiz s'affichent sous forme de card, et contiennent :
> - Tous les champs de la struct Quiz
> - La liste des questions avec leurs types et points
> - Les statistiques d�taill�es (QuizStats)
> - L'historique des sessions et soumissions r�centes
> - Boutons d'action selon les droits

## Comportements sp�cifiques :

### Pour les Quiz d'Offres (type="offer")
- Un quiz d'offre est **OBLIGATOIRE** : si une offre a un QuizID, le candidat DOIT le compl�ter avant de pouvoir valider sa OfferResponse
- Le quiz est automatiquement assign� lors de la r�ception d'une OfferSent
- Une seule tentative par candidat et par quiz d'offre

### Pour les Quiz de Cours (type="course")
- Un quiz de cours peut �tre obligatoire (IsRequired=true) ou optionnel
- Si obligatoire, tous les stagiaires de la session DOIVENT le compl�ter
- Si le stagiaire n'a pas fait le quizz OBLIGATOIRE dans le temps imparti alors il aura 0 ou le niveau le plus bas

### R�gles g�n�rales
- Un quiz ne peut �tre supprim� ou d�sactiv� que s'il n'a aucune session en cours
sauf dans le cas d'une suppression de cours ou suppression d'offre 
- Les statistiques sont mises � jour automatiquement apr�s chaque soumission
- Si une offre est archiv�e/expir�e, il ne sera plus possible de creer une session sur les quiz associ�

## Flows techniques - API

### GET /quizzes
**R�cup�ration de la liste des quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tous les quiz de type "offer"
  - `recruiter_user` : seulement ses quiz de type "offer"
  - `teacher` : seulement ses quiz de type "course"
  - `candidate/student` : quiz accessibles selon les offres/cours
- **Query parameters** :
  - `type` : `offer|course` (optionnel)
  - `status` : `active|inactive` (optionnel)
  - `offerId` : filtrer par offre (optionnel)
  - `courseId` : filtrer par cours (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'�l�ments par page (optionnel)
- **Response** : `200` avec array de quiz
- **Errors** : `401`, `403`, `500`

### GET /quizzes/{id}
**R�cup�ration d'un quiz sp�cifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tout quiz de type "offer"
  - `recruiter_user` : seulement ses quiz de type "offer"
  - `teacher` : seulement ses quiz de type "course"
  - `candidate/student` : quiz accessibles selon droits
- **Response** : `200` avec d�tails complets du quiz
- **Errors** : `401`, `403`, `404`

### POST /quizzes
**Cr�ation d'un nouveau quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - Role `recruiter` pour type "offer"
  - Role `teacher` pour type "course"
- **R�gles** :
  - Type "offer" : doit �tre li� � une offre existante appartenant au recruteur
  - Type "course" : doit �tre li� � un cours existant enseign� par le teacher
  - Au moins une question requise
  - TimeLimit optionnel (en secondes)
  - PassingScore optionnel mais recommand�
  - Pour questions QCM : au moins 2 choix requis
  - Pour questions File : pas de r�ponse correcte
- **Body** : DTO `CreateQuiz`
exemple json :
```json
{
  "title": "string",
  "description": "string", 
  "type": "offer|course",
  "offerId": "string", // requis si type="offer"
  "courseId": "string", // requis si type="course"
  "questions": [
    {
      "question": "string",
      "type": "qcm|text|file",
      "choices": ["string"], // requis si type="qcm", au moins 2 choix
      "answer": "string", // requis si type="qcm" ou "text", null si type="file"
      "point": 10
    }
  ],
  "timeLimit": 1800, // optionnel, en secondes
  "isRequired": true, // pour type="course"
  "passingScore": 70, // optionnel
  "visibilityPassingScore": true
}
```
- **Response** : `201` avec le quiz cr��
- **Errors** : `400`, `401`, `403`, `500`

### POST /quizzes/{id}
**Modification compl�te d'un quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tout quiz de type "offer"
  - `recruiter_user` : seulement ses quiz de type "offer"
  - `teacher` : seulement ses quiz de type "course"
- **R�gles** :
  - Ne peut pas modifier un quiz s'il a d�j� des soumissions
  - Toutes les r�gles de `UpdateQuiz`
- **Body** : DTO `UpdateQuiz`
- **Response** : `200` avec le quiz mis � jour
- **Errors** : `400`, `401`, `403`, `404`, `409` (soumissions existantes), `500`

### POST /quizzes/{id}/toggle-status
**Activation/D�sactivation d'un quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tout quiz de type "offer"
  - `recruiter_user` : seulement ses quiz de type "offer"  
  - `teacher` : seulement ses quiz de type "course"
- **Regles** :
  - Impossible d'effectuer l'action si une session existe sur le quizz
- **Body** : pas de body
- **Response** : `200`
- **Errors** : `401`, `403`, `404`, `500`

### POST /quizzes/{id}/start
**D�marrage d'un quiz par un utilisateur (cr�ation d'une QuizSession)**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate` : pour quiz d'offres re�ues
  - `student` : pour quiz de cours auxquels il participe
- **R�gles** :
  - Le quiz doit �tre actif (IsActive=true)
  - L'utilisateur doit avoir acc�s au quiz (via offre re�ue ou cours)
  - Une seule session active par utilisateur par quiz
  - Pour quiz d'offre : doit avoir re�u l'OfferSent correspondante
  - Pour quiz de cours : doit �tre inscrit dans la session du cours
- **Body** : `{ "type": "offer|course" }` + DTO `createQuizSession`
- **Response** : `201` avec QuizSession cr��e (status="inprogress")
- **Errors** : `400`, `401`, `403`, `409` (d�j� commenc�), `412` (pas d'acc�s), `500`

### POST /quiz-sessions/{id}/save
**Sauvegarde des r�ponses temporaires pendant le quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate/student` : seulement ses propres sessions
- **R�gles** :
  - La QuizSession doit avoir status="inprogress"
  - Mise � jour de TimeWorking et TimeRemaining
  - Validation des types de r�ponses
  - Enregistrement automatique toutes les 30s (en plus du save manuel)
- **Body** : DTO `SaveQuizAnswers`
```json
{
  "answers": [
    {
      "questionId": "string",
      "answer": "string", // optionnel pour type file
      "fileUrl": "string" // optionnel pour type file
    }
  ]
}
```
- **Response** : `200` avec QuizSession mise � jour
- **Errors** : `400`, `401`, `403`, `404`, `408` (session expir�e), `500`

### POST /quiz-sessions/{id}/pause
**Mise en pause d'une session de quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate/student` : seulement ses propres sessions
- **R�gles** :
  - La QuizSession doit avoir status="inprogress"
  - Sauvegarde le temps travaill� et temps restant
- **Body** : DTO `pausedQuizSession`
- **Response** : `200` avec QuizSession (status="paused")
- **Errors** : `401`, `403`, `404`, `409` (d�j� en pause), `500`

### POST /quiz-sessions/{id}/resume
**Reprise d'une session de quiz en pause**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate/student` : seulement ses propres sessions
- **R�gles** :
  - La QuizSession doit avoir status="paused"
  - Red�marre le timer avec le temps restant
- **Body** : pas de body
- **Response** : `200` avec QuizSession (status="inprogress")
- **Errors** : `401`, `403`, `404`, `409` (pas en pause), `500`

### POST /quiz-sessions/{id}/submit
**Soumission finale d'un quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate/student` : seulement ses propres sessions
- **R�gles** :
  - La QuizSession doit avoir status="inprogress", "paused" ou "expired"
  - Toutes les questions doivent avoir une r�ponse
  - Validation des types de r�ponses (QCM, texte, fichier)
  - Calcul automatique du score pour QCM et text
  - Status "evaluated" si toutes questions auto-�valu�es, "submitted" si questions file reste � �valuer
- **Body** : DTO `SubmitQuiz`
```json
{
  "finalAnswers": [
    {
      "questionId": "string",
      "answer": "string",
      "fileUrl": "string" // optionnel pour type file
    }
  ]
}
```
- **Response** : `201` avec QuizSubmission cr��e et QuizSession (status="completed")
- **Errors** : `400`, `401`, `403`, `404`, `500`

### GET /quizzes/{id}/sessions
**R�cup�ration des sessions d'un quiz (affichage d'information simple)**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toutes sessions quiz type "offer"
  - `recruiter_user` : sessions de ses quiz type "offer"
  - `teacher` : sessions de ses quiz type "course"
- **Query parameters** :
  - `status` : `inprogress|paused|expired|completed` (optionnel)
  - `userId` : filtrer par utilisateur (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'�l�ments par page (optionnel)
- **Response** : `200` avec array de QuizSession
- **Errors** : `401`, `403`, `404`, `500`

### GET /quizzes/{id}/submissions
**R�cup�ration des soumissions d'un quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toutes soumissions quiz type "offer"
  - `recruiter_user` : soumissions de ses quiz type "offer"
  - `teacher` : soumissions de ses quiz type "course"
- **Query parameters** :
  - `status` : `submitted|evaluated` (optionnel)
  - `userId` : filtrer par utilisateur (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'�l�ments par page (optionnel)
- **Response** : `200` avec array de QuizSubmission
- **Errors** : `401`, `403`, `404`, `500`

### POST /quiz-submissions/{id}/evaluate
**�valuation d'une soumission de quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute soumission quiz type "offer"
  - `recruiter_user` : soumissions de ses quiz type "offer"
  - `teacher` : soumissions de ses quiz type "course"
- **R�gles** :
  - La soumission doit avoir status="submitted"
  - �valuation manuelle pour questions de type "file"
  - Calcul automatique du score pour QCM et text
  et rajout de l'�valuation des questions files
  - Passage automatique � "evaluated" apr�s validation de l'�valuation
- **Body** : DTO `EvaluateSubmission`
```json
{
  "score": 85, // optionnel si calcul auto
  "feedback": "string", // optionnel
  "passed": true // optionnel si PassingScore d�fini
}
```
- **Response** : `200` avec QuizSubmission �valu�e (status="evaluated")
- **Errors** : `400`, `401`, `403`, `404`, `500`

### GET /quizzes/{id}/stats
**R�cup�ration des statistiques d'un quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : stats de tous quiz type "offer"
  - `recruiter_user` : stats de ses quiz type "offer"
  - `teacher` : stats de ses quiz type "course"
- **Response** : `200` avec QuizStats compl�tes
- **Errors** : `401`, `403`, `404`, `500`

### GET /my-quizzes
**R�cup�ration des quiz d'un candidat/�tudiant**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate` : quiz d'offres re�ues
  - `student` : quiz de cours auxquels il participe
- **Query parameters** :
  - `type` : `offer|course` (optionnel)
  - `status` : `pending|completed` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'�l�ments par page (optionnel)
- **Response** : `200` avec quiz accessibles et leur statut de progression
- **Errors** : `401`, `500`

### GET /my-sessions
**R�cup�ration des sessions d'un utilisateur**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate/student` : seulement ses propres sessions
- **Query parameters** :
  - `quizId` : filtrer par quiz (optionnel)
  - `status` : `inprogress|paused|expired|completed` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'�l�ments par page (optionnel)
- **Response** : `200` avec array de ses QuizSession
- **Errors** : `401`, `500`

### GET /my-submissions
**R�cup�ration des soumissions d'un utilisateur**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate/student` : seulement ses propres soumissions
- **Query parameters** :
  - `quizId` : filtrer par quiz (optionnel)
  - `status` : `submitted|evaluated` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'�l�ments par page (optionnel)
- **Response** : `200` avec array de ses QuizSubmission
- **Errors** : `401`, `500`

### POST /quizzes/{id}/delete
**Suppression d'un quiz**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tout quiz type "offer"
  - `teacher` : seulement ses quiz type "course"
- **R�gles** :
  - Ne peut supprimer que si aucune session n'existe
  sauf si cette suppression viens de sont parent course ou offer
- **Response** : `204`
- **Errors** : `401`, `403`, `404`, `409` (soumissions existantes), `500`

### Codes d'erreur communs
- `400 Bad Request` : Donn�es invalides, r�ponses manquantes
- `401 Unauthorized` : Token manquant ou invalide
- `403 Forbidden` : Permissions insuffisantes
- `404 Not Found` : Ressource non trouv�e
- `408 Request Timeout` : Limite de temps d�pass�e
- `409 Conflict` : Quiz d�j� commenc�, soumissions existantes
- `412 Precondition Failed` : Conditions d'acc�s non remplies
- `500 Internal Server Error` : Erreur serveur

## Flows techniques - �V�NEMENTS

### Event `quiz.created`
**Cr�ation d'un quiz**
- **D�clencheur** : Cr�ation d'un Quiz
- **Action** :
  - Initialisation QuizStats avec valeurs par d�faut
  - Si type="offer" : rattacher � l'offer
  - Si type="course" et IsRequired=true : notification aux �tudiants

### Event `quizSession.started`
**D�marrage d'un quiz par un utilisateur**
- **D�clencheur** : Cr�ation QuizSession avec status="inprogress"
- **Action** :
  - D�marrage du timer si TimeLimit d�fini
  - Notification au cr�ateur du quiz
  - Initialisation du TimeWorking et TimeRemaining

### Event `quizSession.completed`
**Fin d'une session de quiz**
- **D�clencheur** : QuizSession passe � status="completed"
- **Action** :
  - Cr�ation automatique de QuizSubmission
  - Calcul automatique du score pour questions QCM et text
  - Notification au cr�ateur du quiz (recruiter/teacher)

### Event `quizSubmission.submitted`
**Soumission finale cr��e**
- **D�clencheur** : Cr�ation QuizSubmission avec status="submitted"
- **Action** :
  - Mise � jour des QuizStats
  - Si quiz d'offre : d�bloque la possibilit� de cr�er OfferResponse
  - Si quiz de cours obligatoire : marque comme compl�t�

### Event `quizSubmission.evaluated`
**�valuation d'une soumission**
- **D�clencheur** : QuizSubmission passe � status="evaluated"
- **Action** :
  - Notification � l'utilisateur avec r�sultat et feedback
  - Mise � jour finale des QuizStats

### Event `offer.deleted`
**Impact sur les quiz d'offres**
- **D�clencheur** : Offre supprim�e, archiv�e ou expir�e
- **Action** :
  - Quiz associ� sont supprimer
  - Sessions de quiz supprim�es
  - Notification aux candidats ayant des sessions en cours ou soumission
  - QuizStats supprimer

### Event `offer.archived|expired`
**Impact sur les quiz d'offres**
- **D�clencheur** : Offre supprim�e, archiv�e ou expir�e
- **Action** :
  - Quiz impossible de creer de nouvelle session (IsActive=false)
  - Sessions de quiz en cours (status="inprogress|paused|expired") seront supprim�es
  - Notification aux candidats ayant des sessions en cours

### Event `course.deleted`
**Impact sur les quiz de cours**
- **D�clencheur** : Cours supprim� ou archiv�
- **Action** :
  - Quiz associ� sont supprimer
  - Sessions de quiz supprim�es
  - Notification aux candidats ayant des sessions en cours ou soumission
  - QuizStats supprimer

### Event `course.archived`
**Impact sur les quiz de cours**
- **D�clencheur** : Cours supprim� ou archiv�
- **Action** :
  - Quiz impossible de creer de nouvelle session (IsActive=false)
  - Sessions de quiz en cours (status="inprogress|paused|expired") seront supprim�es
  - Notification aux candidats ayant des sessions en cours

## R�gles de validation

### Contraintes de base de donn�es
- Index unique sur `quiz_id + user_id` pour QuizSession (une session par utilisateur par quiz)
- Index unique sur `quiz_session_id` pour QuizSubmission (une soumission par session)
- Index compos� sur `quiz_id + status` pour les requ�tes de sessions/soumissions
- Index sur `user_id + created_at` pour l'historique utilisateur
- Index sur `offer_id` et `course_id` pour les liaisons
- Index sur `type + is_active` pour les recherches de quiz

### Contraintes m�tier
- Type "offer" : OfferID requis, CourseID interdit, RecruiterID requis
- Type "course" : CourseID requis, OfferID interdit, TeacherID requis
- Questions : au moins une question par quiz
- Questions QCM : au moins 2 choix requis
- Questions File : Answer doit �tre null
- TimeLimit : si d�fini, doit �tre > 0
- PassingScore : si d�fini, doit �tre entre 0 et la somme des points des questions
- QuizSession : une seule session active par utilisateur par quiz
- QuizSubmission : une seule soumission par session de quiz
