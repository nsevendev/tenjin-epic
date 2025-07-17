# **Entités principales / Schémas**

---

- [le lien mcd (faudra peut etre creer un compte)](https://dbdiagram.io/d/tenjin-6875d886f413ba3508e07ee5)

## 1. **Institute**

```json
{
  "_id": ObjectId,
  "business_name": String,
  "siret": String,
  "address": String,
  "zip_code": String,
  "city": String,
  "contact_emails": [String],
  "formations": [ObjectId],         // Formations créées/gérées par cet organisme
  "users": [ObjectId],              // Utilisateurs membres de cet organisme
  "created_at": Date,
  "updated_at": Date
}
```

**Description :**

* Identifie chaque organisme de formation.
* Contient les liens vers ses formations et utilisateurs.
* Le SIRET est obligatoire pour assurer la traçabilité légale.

---

## 2. **User**

```json
{
  "_id": ObjectId,
  "firstname": String,
  "lastname": String,
  "email": String,
  "roles": [ "student", "trainer", "manager", "admin", "recruiter" ], // Rôles dans la plateforme
  "organizations": [ObjectId],        // Organismes auxquels l'utilisateur appartient
  "sessions": [ObjectId],             // Sessions suivies ou encadrées par l'utilisateur
  "competence_records": [             // Historique complet des compétences (tous organismes/sessions)
    {
      "competence_id": ObjectId,
      "history": [
        {
          "date": Date,
          "level": String,            // Niveau atteint (ex : bronze, 80/100, etc.)
          "validated_by": ObjectId,   // Org ou utilisateur ayant validé la compétence
          "session_id": ObjectId,     // Session liée à la validation
          "notes": String
        }
      ]
    }
  ],
  "external_experiences": [
    {
      "title": String,
      "description": String,
      "date": Date,
      "proofs": [ObjectId]            // Preuves/diplômes associés
    }
  ],
  "status": String,                   // "training", "employed", "jobseeker", etc.
  "availability": [
    {
      "start_date": Date,
      "end_date": Date,
      "type": String                  // (ex : disponible, indisponible, en congé, etc.)
    }
  ],
  "received_offers": [ObjectId],      // Offres d'emploi reçues (jamais publiques)
  "pending_share_requests": [         // Lorsqu'un recruteur demande des infos supplémentaires
    {
      "offer_id": ObjectId,
      "fields_requested": [String],   // ex : téléphone, email, CV, etc.
      "status": "pending" | "accepted" | "rejected"
    }
  ],
  "quiz_results": [
    {
      "quiz_id": ObjectId,
      "result": String,
      "details": Object
    }
  ],
  "chats": [ObjectId],                // Discussions privées (avec recruteurs après contact)
  "created_at": Date,
  "updated_at": Date
}
```

**Description :**

* **Anonymat** : les infos personnelles (hors prénom/nom et compétences) ne sont partagées qu’après consentement.
* **Compétences** : historique détaillé par compétence et par session/org.
* **Disponibilités** : pour filtrer facilement les profils côté recruteur.

---

## 3. **Formation**

*(Liée aux métiers/compétences du catalogue France Travail)*

```json
{
  "_id": ObjectId,
  "title": String,
  "description": String,
  "organization_id": ObjectId,
  "course_ids": [ObjectId],          // Liste ordonnée des cours/modules
  "competence_ids": [ObjectId],      // Compétences principales visées
  "external_job_ref": String,        // Référence métier France Travail (code ROME ou autre)
  "sessions": [ObjectId],
  "meta": Object,
  "created_at": Date
}
```

**Description :**

* **external\_job\_ref** : lien vers un métier ou une fiche métier du catalogue officiel France Travail.
* Permet d'importer ou d’associer des compétences normalisées.

---

## 4. **Session**

```json
{
  "_id": ObjectId,
  "formation_id": ObjectId,
  "organization_id": ObjectId,
  "title": String,                      // Nom personnalisé de la session si besoin
  "users": [ObjectId],                  // Stagiaires de la session
  "trainers": [ObjectId],               // Formateurs ou intervenants
  "start_date": Date,
  "end_date": Date,
  "course_ids": [ObjectId],             // Ordre ou contenu des cours de cette session
  "resources": [ObjectId],              // Ressources additionnelles spécifiques à la session
  "evaluations": [ObjectId],            // Évaluations réalisées pendant la session
  "quizzes": [ObjectId],                // Quiz associés à la session
  "chats": [ObjectId],                  // Salons/discussions liés à la session
  "calendar_id": ObjectId,              // Calendrier de la session
  "attendance_sheet_id": ObjectId,      // Feuille de présence
  "created_at": Date
}
```

**Description :**

* Point central de toutes les fonctionnalités "live" : calendrier, chat, présence, quiz…
* **course\_ids** : possibilité de personnaliser la structure pour une session donnée.

---

## 5. **Course (Cours)**

```json
{
  "_id": ObjectId,
  "title": String,
  "description": String,
  "content_blocks": [                   // Contenus du cours (texte, fichiers...)
    { "type": "text", "data": String },
    { "type": "pdf",  "url": String, "name": String }
  ],
  "competence_ids": [ObjectId],         // Compétences abordées dans le cours
  "resource_ids": [ObjectId],           // Ressources partagées/réutilisées
  "created_at": Date,
  "updated_at": Date
}
```

**Description :**

* **Réutilisable** entre plusieurs formations/sessions.
* Contenus variés : texte, PDF, liens, etc.

---

## 6. **Resource (Ressource)**

```json
{
  "_id": ObjectId,
  "type": String,                       // "pdf", "video", "image", etc.
  "url": String,
  "name": String,
  "uploaded_by": ObjectId,
  "created_at": Date
}
```

**Description :**

* Fichier ou média attachable à un cours ou une session.
* Permet de partager ou versionner facilement les documents.

---

## 7. **Competence (Compétence)**

```json
{
  "_id": ObjectId,
  "name": String,
  "description": String,
  "domain": String,
  "levels": [ "beginner", "intermediate", "advanced", "expert" ], // Selon l’orga ou importé
  "external_skill_id": String,         // Lien avec une compétence officielle France Travail
  "meta": Object
}
```

**Description :**

* Peut être importée ou synchronisée avec la base France Travail (ROME/Compétences).
* **external\_skill\_id** : référence externe pour garder la cohérence nationale.

---

## 8. **Evaluation (Évaluation)**

```json
{
  "_id": ObjectId,
  "session_id": ObjectId,
  "user_id": ObjectId,
  "course_id": ObjectId,                 // Si l'évaluation porte sur un cours précis
  "competence_id": ObjectId,
  "level": String,                       // Niveau atteint
  "validated": Boolean,
  "validator_id": ObjectId,              // Formateur/gestionnaire validateur
  "date": Date,
  "notes": String
}
```

**Description :**

* Utilisé pour alimenter l’historique de compétences du profil utilisateur.

---

## 9. **AttendanceSheet (Feuille de présence)**

```json
{
  "_id": ObjectId,
  "session_id": ObjectId,
  "records": [
    {
      "user_id": ObjectId,
      "date": Date,
      "status": "present" | "absent" | "late" | "excused",
      "notes": String
    }
  ],
  "created_at": Date
}
```

**Description :**

* Permet un suivi fin des présences sur la durée de la session.
* **status** permet de gérer toutes les situations (retard, absence justifiée…).

---

## 10. **ChatRoom (Salon de discussion)**

```json
{
  "_id": ObjectId,
  "session_id": ObjectId,                // Salon de session ou null pour chat privé recruteur
  "type": "info" | "private" | "moderated" | "recruiter", // Types prédéfinis
  "title": String,
  "participants": [ObjectId],
  "messages": [ObjectId],                // Messages du salon (ou embarqué)
  "created_at": Date,
  "moderators": [ObjectId],              // Gestionnaires/modérateurs
  "auto_delete_at": Date,                // Pour chat privé recruteur (1 mois)
  "status": "active" | "archived" | "deleted"
}
```

**Description :**

* Salons d'information, salons privés, discussions modérées, chats recruteur.
* Les chats privés (recruteur-candidat) sont créés après contact positif et auto-supprimés au bout d’1 mois.

---

## 11. **Recruiter (Recruteur)**

```json
{
  "_id": ObjectId,
  "siret": String,
  "name": String,
  "contact_email": String,
  "authorized": Boolean,                 // Validation par l’admin plateforme
  "search_history": [Object],            // Historique des recherches (pour stats/audit)
  "offers_sent": [ObjectId],             // Offres envoyées (historique)
  "chats": [ObjectId],                   // Discussions actives avec candidats
  "created_at": Date
}
```

**Description :**

* Entité dédiée aux recruteurs, le SIRET garantit leur légitimité.
* Permet de tracer recherches et envois d’offres pour audit.

---

## 12. **Offer (Offre d'emploi)**

```json
{
  "_id": ObjectId,
  "recruiter_id": ObjectId,
  "title": String,
  "message": String,                     // Message personnalisé ou résumé de l’offre
  "attachment_url": String,              // Fichier PDF si fourni
  "candidate_ids": [ObjectId],           // Utilisateurs ciblés par l’offre
  "status": String,                      // "sent", "viewed", "accepted", "declined"
  "quiz_id": ObjectId,                   // Quiz éventuel lié à l’offre
  "responses": [
    {
      "user_id": ObjectId,
      "status": String,                  // Réponse de l’utilisateur
      "shared_fields": [String],         // Infos que l’utilisateur a accepté de partager
      "reply_date": Date
    }
  ],
  "created_at": Date
}
```

**Description :**

* Jamais publique, uniquement envoyée aux candidats choisis.
* Historique de qui a partagé quoi, à quelle date.

---

## 13. **Quiz**

```json
{
  "_id": ObjectId,
  "recruiter_id": ObjectId,
  "title": String,
  "description": String,
  "questions": [
    {
      "question": String,
      "type": "qcm" | "text" | "file",
      "choices": [String],              // Si QCM
      "answer": String
    }
  ],
  "created_at": Date
}
```

**Description :**

* Peut être attaché à une offre ou envoyé directement à un candidat.

---

## 14. **JobSkillCatalog (Service, pas collection)**

* **Rôle :** Service indépendant qui synchronise le catalogue officiel France Travail/Pôle Emploi (métiers, fiches métiers, compétences, codes ROME).
* **Intégration :**

    * Chaque formation ou compétence peut référencer un champ `external_job_ref` ou `external_skill_id` pour garder le lien avec le catalogue.
    * Suggestion automatique de métiers/compétences lors de la création d’une formation.
* **Fonctionnement :**

    * **Synchronisation quotidienne** pour garder la base à jour.
    * **Interface d’admin** pour lier/matcher formations et compétences avec le référentiel.
    * Utilisé comme source de vérité pour la structuration des formations.

---

# **Résumé des relations clés**

* **User** ⇄ **Organization** (appartenance, rôles)
* **User** ⇄ **Session** (participation, évaluation, présence, chat)
* **Session** ⇄ **Formation** (instance, personnalisation possible)
* **Session** ⇄ **Course** (cours planifiés, évalués, enrichis)
* **Course/Session/Formation** ⇄ **Competence** (lié au catalogue France Travail)
* **User** ⇄ **Competence** (historique, validations multi-orgas/sessions)
* **Recruiter** ⇄ **Offer** ⇄ **User** (candidature privée, jamais publique)
* **Recruiter** ⇄ **Quiz** ⇄ **User** (en lien avec offre ou seul)
* **Recruiter** ⇄ **ChatRoom** ⇄ **User** (après premier contact positif)
* **Session** ⇄ **ChatRoom/AttendanceSheet/Calendar** (outils d’organisation)

---

# **Intégration France Travail / JobSkillCatalog**

* **Service indépendant :**

    * Expose des endpoints pour rechercher/importer métiers, fiches métiers, compétences officielles.
    * Permet d’associer une formation ou compétence à un métier/compétence officielle (ex : code ROME).
    * Propose automatiquement la liste officielle lors de la création ou la modification d’une formation.
* **Avantages :**

    * Les formations et compétences sont toujours alignées avec les standards nationaux.
    * Le carnet de compétences utilisateur est valorisable dans tous les contextes officiels.

---
