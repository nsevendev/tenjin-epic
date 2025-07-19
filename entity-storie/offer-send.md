# OfferSent

> Qui a accès aux envois d'offres ?
> - Un recruteur authentifié qui a créé l'offre ou l'envoi
> - Un candidat authentifié qui a reçu l'offre

## Role "recruiter" avec SubRole "recruiter_admin"

- Champs autorisés :
  - tous les champs de la struct OfferSent
- Lecture :
  - peut voir tous les envois d'offres de toutes les offres
  - peut voir tous les détails d'un envoi d'offre
- Écriture :
  - peut créer tous les envois d'offre à un ou plusieurs candidats de toutes les offres
  - peut modifier le statut des envois d'offres
  - peut supprimer tous les envois d'offres

## Role "recruiter" avec SubRole "recruiter_user"

- Champs autorisés :
  - tous les champs de la struct OfferSent
- Lecture :
  - peut voir seulement les envois d'offres liés à ses propres offres
  - peut voir tous les détails d'un envoi d'offre qu'il a créé
- Écriture :
  - peut créer des envois d'offre seulement pour ses propres offres
  - peut modifier le statut des envois d'offres qu'il a créés

## Role "candidate"

- Champs autorisés :
  - tous les champs sauf RecruiterID et CompanyID (masqués côté client)
- Lecture :
  - peut voir seulement les envois d'offres qui lui sont destinés
  - peut voir tous les détails d'un envoi d'offre qui lui est destiné
- Écriture :
  - peut modifier le statut de ses envois d'offres (`viewed`, `responded`)

> La liste des envois d'offres s'affiche sous forme de tableau, avec les colonnes suivantes :
> - Titre de l'offre
> - Nom du candidat
> - Statut (sent, viewed, responded) - modifiable avec un select pour les recruteurs
> - Date d'envoi
> - Date de visualisation
> - Date de réponse
> - Message personnalisé (aperçu)
> - Actions (voir, supprimer)

> Les détails d'un envoi d'offre s'affichent sous forme de card, et contiennent :
> - Tous les champs de la struct OfferSent
> - Les informations de l'offre liée
> - Les informations du candidat (nom, prénom, email)
> - Le message personnalisé complet
> - L'historique des changements de statut

## Comportements spécifiques :

- Quand un candidat visualise une offre, le statut passe automatiquement à `viewed` et `ViewedAt` est mis à jour
- Quand un candidat répond à une offre, le statut passe à `responded` et `RespondedAt` est mis à jour
- Si une offre est supprimée, archivée ou expirée, tous les envois liés sont automatiquement supprimés
- Un même candidat ne peut recevoir qu'un seul envoi par offre (contrainte unique sur OfferID + CandidateID)

## Flows techniques - API

### GET /offer-sends
**Récupération de la liste des envois d'offres**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tous les envois d'offres
  - `recruiter_user` : seulement les envois liés à ses offres
  - `candidate` : seulement ses propres envois reçus
- **Query parameters** :
  - `status` : `sent|viewed|responded` (optionnel)
  - `offerId` : filtrer par offre (optionnel)
  - `candidateId` : filtrer par candidat (optionnel, recruiter seulement)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'éléments par page (optionnel)
- **Response** : `200` avec array d'envois d'offres
- **Errors** : `401`, `403`, `500`

### GET /offer-sends/{id}
**Récupération d'un envoi d'offre spécifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tout envoi d'offre
  - `recruiter_user` : seulement les envois liés à ses offres
  - `candidate` : seulement ses propres envois reçus
- **Response** : `200` avec détails complets de l'envoi d'offre
- **Errors** : `401`, `403`, `404`, `500`

### POST /offer-sends
**Création d'un nouvel envoi d'offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - Role `recruiter` requis
- **Règles** :
  - L'offre doit exister et être active
  - Le candidat doit exister
  - Un seul envoi par candidat et par offre (contrainte unique)
  - Le recruteur doit être propriétaire de l'offre (sauf recruiter_admin)
- **Body** : DTO `CreateOfferSent`
exemple court :
```json
{
  "offerId": "string",
  "candidateIds": "string",
  "message": "string (optionnel)"
}
```
- **Response** : `201` avec l'envoi d'offre créé
- **Errors** : `400`, `401`, `403`, `409` (déjà envoyé), `500`

### POST /offer-sends/bulk
**Création d'envois d'offres en masse**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - Role `recruiter` requis
- **Règles** :
  - Mêmes règles que l'envoi simple
  - Maximum 50 candidats par envoi en masse
- **Body** : DTO `CreateBulkOfferSent`
exemple court :
```json
{
  "offerId": "string",
  "candidateIds": ["string"],
  "message": "string (optionnel)"
}
```
- **Response** : `201` avec la liste des envois créés et les échecs
- **Errors** : `400`, `401`, `403`, `500`

### POST /offer-sends/{id}/status
**Modification du statut d'un envoi d'offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tout envoi d'offre
  - `recruiter_user` : seulement les envois liés à ses offres
  - `candidate` : seulement ses propres envois (pour viewed/responded)
- **Règles** :
  - Les candidats ne peuvent que passer de `sent` à `viewed` ou de `viewed` à `responded`
  - Le status est a `send` à l'envoie de l'offre 
- **Body** : DTO `UpdatedStatusOfferSend`
exemple
```json
{
  "status": "viewed|responded"
}
```
- **Response** : `200` avec l'envoi d'offre mis à jour
- **Errors** : `400`, `401`, `403`, `404`, `500`

### DELETE /offer-sends/{id}
**Suppression d'un envoi d'offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : tout envoi d'offre
  - `recruiter_user` : seulement les envois liés à ses offres
  - `candidate` : a réflechir mais l'idee c'est que si la personne n'est plus disponible il peut supprimer l'offersend
- **Response** : `204`
- **Errors** : `401`, `403`, `404`, `500`

### GET /offers/{offerId}/sends
**Récupération des envois d'une offre spécifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute offre
  - `recruiter_user` : seulement ses propres offres
- **Query parameters** :
  - `status` : `sent|viewed|responded` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'éléments par page (optionnel)
- **Response** : `200` avec array d'envois d'offres pour cette offre
- **Errors** : `401`, `403`, `404`, `500`

### Codes d'erreur communs
- `400 Bad Request` : Données invalides ou manquantes
- `401 Unauthorized` : Token manquant ou invalide
- `403 Forbidden` : Permissions insuffisantes
- `404 Not Found` : Ressource non trouvée
- `409 Conflict` : Envoi déjà existant pour ce candidat/offre
- `500 Internal Server Error` : Erreur serveur

## Flows techniques - ÉVÉNEMENTS

### Event `offerSent.created`
**Envoi d'une offre à un candidat**
- **Déclencheur** : Création d'un OfferSent
- **Action** :
  - Notification email au candidat
  - Notification push si activée
  - Mise à jour des statistiques de l'offre (si il y en a)

### Event `offerSent.viewed`
**Candidat visualise une offre**
- **Déclencheur** : Statut passe à `viewed`
- **Action** :
  - Notification au recruteur
  - Mise à jour des statistiques de l'offre

### Event `offerSent.responded`
**Candidat répond à une offre**
- **Déclencheur** : Statut passe à `responded`
- **Action** :
  - Notification au recruteur
  - Mise à jour des statistiques de l'offre (si il y en a)
  - Déclenchement du processus de suivi (ouverture de chat privé entre recruteur et candidat)

### Event `offer.deleted|archived|expired`
**Suppression automatique des envois**
- **Déclencheur** : Offre supprimée, archivée ou expirée
- **Action** :
  - Suppression de tous les OfferSent liés
  - Notification aux candidats concernés
  - suppression des chats privé
