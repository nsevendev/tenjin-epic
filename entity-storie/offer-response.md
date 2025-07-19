# OfferResponse

> Qui a accès aux réponses d'offres ?
> - Un candidat authentifié pour ses propres réponses
> - Un recruteur authentifié pour les réponses aux offres qu'il a créées

## Role "recruiter" avec SubRole "recruiter_admin"

- Champs autorisés :
  - tous les champs de la struct OfferResponse
- Lecture :
  - peut voir toutes les réponses d'offres de toutes les offres
  - peut voir tous les détails d'une réponse d'offre
- Écriture :
  - aucune (les réponses ne peuvent être créées que par les candidats)

## Role "recruiter" avec SubRole "recruiter_user"

- Champs autorisés :
  - tous les champs de la struct OfferResponse
- Lecture :
  - peut voir seulement les réponses d'offres liées à ses propres offres
  - peut voir tous les détails d'une réponse d'offre qu'il a créée
- Écriture :
  - aucune (les réponses ne peuvent être créées que par les candidats)

## Role "candidate"

- Champs autorisés :
  - tous les champs sauf RecruiterID et CompanyID (masqués côté client)
- Lecture :
  - peut voir seulement ses propres réponses d'offres
  - peut voir tous les détails de ses réponses d'offres
- Écriture :
  - peut créer une réponse à une offre qu'il a reçue
  - ne peut pas modifier une réponse une fois créée (à discuter)

> La liste des réponses d'offres s'affiche sous forme de tableau, avec les colonnes suivantes :
> - Titre de l'offre
> - Nom du candidat (pour les recruteurs)
> - Statut (accepted, declined) avec indicateur visuel (vert/rouge)
> - Date de réponse
> - Champs partagés (si accepted)
> - Actions (voir, accès chat (si accepted))

> Les détails d'une réponse d'offre s'affichent sous forme de card, et contiennent :
> - Tous les champs de la struct OfferResponse
> - Les informations de l'offre liée
> - Les informations du candidat (nom, prénom, email)
> - Liste détaillée des champs partagés (si accepted)
> - Bouton d'accès au chat privé (si accepted)

## Comportements spécifiques :

- Un candidat ne peut répondre qu'une seule fois à chaque offre (contrainte unique sur OfferSentID)
- Si le statut est `declined`, le tableau `sharedFields` doit être vide
- Si le statut est `accepted`, le candidat peut choisir les champs de profil à partager
- Les champs partagés autorisés : `email`, `phone`, `cv`, `linkedin`, `github`, `skills`, `experience`, `location`, `identity` (à modifier)
- Quand une réponse `accepted` est créée, un ProfileGrant est automatiquement généré
- Quand une réponse `accepted` est créée, un chat privé est ouvert entre le recruteur et le candidat
- Les réponses ne peuvent jamais être supprimées manuellement
- Si une offre est supprimée, archivée ou expirée, toutes les réponses liées sont automatiquement supprimées

## Flows techniques - API

### GET /offer-responses
**Récupération de la liste des réponses d'offres**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toutes les réponses d'offres
  - `recruiter_user` : seulement les réponses liées à ses offres
  - `candidate` : seulement ses propres réponses
- **Query parameters** :
  - `status` : `accepted|declined` (optionnel)
  - `offerId` : filtrer par offre (optionnel)
  - `candidateId` : filtrer par candidat (optionnel, recruiter seulement)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'éléments par page (optionnel)
- **Response** : `200` avec array de réponses d'offres
- **Errors** : `401`, `403`, `500`

### GET /offer-responses/{id}
**Récupération d'une réponse d'offre spécifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute réponse d'offre
  - `recruiter_user` : seulement les réponses liées à ses offres
  - `candidate` : seulement ses propres réponses
- **Response** : `200` avec détails complets de la réponse d'offre
- **Errors** : `401`, `403`, `404`

### POST /offer-responses
**Création d'une nouvelle réponse d'offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - Role `candidate` uniquement
- **Règles** :
  - L'OfferSent doit exister et être destiné au candidat authentifié
  - L'offre doit être active
  - Le candidat ne peut répondre qu'une seule fois (contrainte unique sur OfferSentID)
  - Si status = `declined`, sharedFields doit être vide
  - Si status = `accepted`, sharedFields peut contenir les valeurs autorisées
  - Les valeurs autorisées pour sharedFields : `email`, `phone`, `cv`, `linkedin`, `github`, `skills`, `experience`, `location`, `identity` (à modifier)
- **Body** : DTO `CreateOfferResponse`
exemple json :
```json
{
  "offerSentId": "string",
  "status": "accepted|declined",
  "sharedFields": ["email", "phone", "cv"] // optionnel, seulement si status = accepted
}
```
- **Response** : `201` avec la réponse d'offre créée
- **Errors** : `400`, `401`, `403`, `409` (déjà répondu), `500`

### GET /offers/{offerId}/responses
**Récupération des réponses d'une offre spécifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute offre
  - `recruiter_user` : seulement ses propres offres
- **Query parameters** :
  - `status` : `accepted|declined` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'éléments par page (optionnel)
- **Response** : `200` avec array de réponses d'offres pour cette offre
- **Errors** : `401`, `403`, `404`, `500`

### GET /candidates/{candidateId}/responses
**Récupération des réponses d'un candidat spécifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate` : seulement ses propres réponses
  - `recruiter_admin` : toutes les réponses du candidat
  - `recruiter_user` : seulement les réponses du candidat liées à ses offres
- **Query parameters** :
  - `status` : `accepted|declined` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'éléments par page (optionnel)
- **Response** : `200` avec array de réponses d'offres du candidat
- **Errors** : `401`, `403`, `404`, `500`

### Codes d'erreur communs
- `400 Bad Request` : Données invalides ou manquantes, sharedFields non vide avec status declined
- `401 Unauthorized` : Token manquant ou invalide
- `403 Forbidden` : Permissions insuffisantes
- `404 Not Found` : Ressource non trouvée
- `409 Conflict` : Réponse déjà existante pour cet OfferSent
- `500 Internal Server Error` : Erreur serveur

## Flows techniques - ÉVÉNEMENTS

### Event `offerResponse.accepted`
**Candidat accepte une offre**
- **Déclencheur** : Création d'une OfferResponse avec status `accepted`
- **Action** :
  - Mise à jour du statut OfferSent correspondant à `responded`
  - Mise à jour de `RespondedAt` dans OfferSent
  - Création d'un ProfileGrant avec les champs partagés
  - Ouverture d'un chat privé entre le recruteur et le candidat
  - Notification au recruteur (email + push)
  - Mise à jour des statistiques de l'offre (si il y en a)

### Event `offerResponse.declined`
**Candidat refuse une offre**
- **Déclencheur** : Création d'une OfferResponse avec status `declined`
- **Action** :
  - Mise à jour du statut OfferSent correspondant à `responded`
  - Mise à jour de `RespondedAt` dans OfferSent
  - Notification au recruteur (email + push)
  - Mise à jour des statistiques de l'offre (si il y en a)

### Event `offer.deleted|archived|expired`
**Suppression automatique des réponses**
- **Déclencheur** : Offre supprimée, archivée ou expirée
- **Action** :
  - Suppression de toutes les OfferResponse liées
  - Suppression des ProfileGrant associés
  - Fermeture (suppression) des chats privés associés
  - Notification aux candidats concernés

## Règles de validation

### Contraintes de base de données
- Index unique sur `offer_sent_id` (un candidat ne peut répondre qu'une seule fois)
- Index composé sur `offer_id + status` (pour les statistiques)
- Index sur `candidate_id + created_at` (pour l'historique candidat)