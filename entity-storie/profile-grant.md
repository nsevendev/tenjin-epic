# ProfileGrant

> Système de permissions granulaires pour contrôler l'accès aux données de profil des candidats/stagiaires
> 
> **Principe** : Un candidat peut créer plusieurs ProfileGrant pour donner accès à différents champs de son profil à différents types d'entités (recruteurs, entreprises, instituts, audiences publiques)

## Vue d'ensemble du système

### Création automatique
- **À l'inscription** : Un ProfileGrant par défaut est créé automatiquement
- **Lors d'une OfferResponse accepted** : Un ProfileGrant spécifique au recruteur est créé avec les `SharedFields` sélectionnés

### Types de scope (ScopeType)
1. **`recruiter`** - Permission pour un recruteur spécifique (RequesterID obligatoire)
2. **`company`** - Permission pour une entreprise entière (CompanyID obligatoire)  
3. **`institute`** - Permission pour un institut de formation (InstituteID obligatoire)
4. **`audience`** - Permission pour un groupe d'utilisateurs (Audience obligatoire)

### Types d'audience (pour ScopeType = "audience")
1. **`recruiter`** - Tous les recruteurs
2. **`institute`** - Tous les instituts
3. **`company`** - Toutes les entreprises
4. **`public`** - Accès public (tous les utilisateurs authentifiés)

## Logique de permissions

### Qui a accès aux ProfileGrant ?
> - Le candidat propriétaire du profil (lecture/écriture complète)
> - Les entités ayant des grants valides (lecture des champs accordés uniquement)

## Role "candidate" (propriétaire)

- **Champs autorisés** : tous les champs de la struct ProfileGrant
- **Lecture** : 
  - peut voir tous ses ProfileGrant
  - peut voir l'historique de création/modification
- **Écriture** :
  - peut créer de nouveaux ProfileGrant
  - peut modifier les GrantedFields de ses ProfileGrant existants
  - peut révoquer (Revoked = true) ses ProfileGrant
  - peut définir des dates d'expiration
  - ne peut pas supprimer définitivement (audit trail)

## Role "recruiter" avec SubRole "recruiter_admin"

- **Champs autorisés** : lecture seule des ProfileGrant qui les concernent
- **Lecture** :
  - peut voir les ProfileGrant où ils sont bénéficiaires
  - peut voir les champs accordés dans GrantedFields
- **Écriture** : aucune (seul le candidat contrôle ses permissions)

## Role "recruiter" avec SubRole "recruiter_user"

- **Champs autorisés** : lecture seule des ProfileGrant qui les concernent
- **Lecture** :
  - peut voir les ProfileGrant où ils sont bénéficiaires
  - peut voir les champs accordés dans GrantedFields
- **Écriture** : aucune

## Règles de validation

### Champs accordés disponibles
Les `GrantedFields` peuvent contenir :
- `email` - Adresse email du candidat
- `phone` - Numéro de téléphone
- `cv` - Curriculum vitae (fichier/lien)
- `linkedin` - Profil LinkedIn
- `github` - Profil GitHub
- `skills` - Compétences techniques
- `experience` - Expérience professionnelle
- `location` - Localisation géographique
- `identity` - Informations d'identité (nom, prénom)
> completer cette liste

## Comportements spécifiques

### Création automatique lors d'OfferResponse
Quand un candidat accepte une offre (`OfferResponse.status = "accepted"`), un ProfileGrant est automatiquement créé :
```json
{
  "candidateId": "candidate_id_from_response",
  "recruiterId": "recruiter_id_from_offer", 
  "scopeType": "recruiter",
  "grantedFields": ["..."] // Copie de OfferResponse.sharedFields
}
```

### Hiérarchie d'accès
- Tout les grant s'override, en partant du plus petit (donc specifique à une personne)  
si y a pas on passe au dessus donc le groupe

### Contrôle d'accès au profil
Quand un utilisateur consulte un profil candidat :
1. Le système vérifie tous les ProfileGrant actifs pour ce candidat
2. Détermine les champs visibles selon le rôle de l'utilisateur
3. Filtre les données affichées selon les GrantedFields

## Flows techniques - API

### GET /profile-grants
**Récupération des grants du candidat authentifié**
- **Authentification** : Token JWT requis
- **Autorisation** : Role `candidate` uniquement pour ses propres grants
- **Query parameters** :
  - `scopeType` : `recruiter|company|institute|audience` (optionnel)
  - `revoked` : `true|false` (optionnel)
  - `expired` : `true|false` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'éléments par page (optionnel)
- **Response** : `200` avec array de ProfileGrant
- **Errors** : `401`, `403`, `500`

### GET /profile-grants/{id}
**Récupération d'un grant spécifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `candidate` : seulement ses propres grants
  - `recruiter` : seulement les grants qui les concernent
- **Response** : `200` avec détails complets du ProfileGrant
- **Errors** : `401`, `403`, `404`

### POST /profile-grants
**Création d'un nouveau ProfileGrant**
- **Authentification** : Token JWT requis
- **Autorisation** : Role `candidate` uniquement
- **Règles** :
  - Validation conditionnelle selon ScopeType
  - GrantedFields doivent être des valeurs autorisées
  - Si ExpiresAt est défini, doit être dans le futur
- **Body** : DTO `CreateProfileGrant`
exemple json :
```json
{
  "scopeType": "recruiter|company|institute|audience",
  "recruiterId": "string", // si scopeType = recruiter
  "companyId": "string", // si scopeType = company  
  "audience": "recruiter|institute|company|public", // si scopeType = audience
  "grantedFields": ["email", "cv", "skills"],
  "expiresAt": "2024-12-31T23:59:59Z" // optionnel
}
```
- **Response** : `201` avec le ProfileGrant créé
- **Errors** : `400`, `401`, `403`, `500`

### PUT /profile-grants/{id}
**Modification d'un ProfileGrant**
- **Authentification** : Token JWT requis
- **Autorisation** : Role `candidate` uniquement pour ses propres grants
- **Règles** :
  - Ne peut pas modifier ScopeType, RecruiterID, CompanyID après création
  - Peut modifier GrantedFields et ExpiresAt
  - Peut révoquer (Revoked = true) mais pas annuler la révocation
- **Body** : DTO `UpdateProfileGrant`
exemple json :
```json
{
  "grantedFields": ["email", "phone"],
  "expiresAt": "2025-01-31T23:59:59Z",
  "revoked": true
}
```
- **Response** : `200` avec le ProfileGrant mis à jour
- **Errors** : `400`, `401`, `403`, `404`, `500`

### POST /profile-grants/{id}/revoke
**Révocation d'un ProfileGrant**
- **Authentification** : Token JWT requis
- **Autorisation** : Role `candidate` uniquement pour ses propres grants
- **Body** : pas de body
- **Response** : `200`
- **Errors** : `401`, `403`, `404`, `500`

### GET /candidates/{candidateId}/accessible-fields
**Récupération des champs accessibles d'un profil candidat**
- **Authentification** : Token JWT requis
- **Autorisation** : Selon les ProfileGrant du candidat
- **Response** : `200` avec liste des champs accessibles pour l'utilisateur authentifié
exemple json :
```json
{
  "accessibleFields": ["email", "skills"],
  "grantSource": "recruiter|company|institute|audience",
  "grantId": "string"
}
```
- **Errors** : `401`, `403`, `404`

### Codes d'erreur communs
- `400 Bad Request` : Validation des champs conditionnels échouée
- `401 Unauthorized` : Token manquant ou invalide
- `403 Forbidden` : Permissions insuffisantes
- `404 Not Found` : Ressource non trouvée
- `409 Conflict` : Grant similaire déjà existant
- `500 Internal Server Error` : Erreur serveur

## Flows techniques - ÉVÉNEMENTS

### Event `profileGrant.created`
**Création d'un nouveau grant**
- **Déclencheur** : Création d'un ProfileGrant
- **Action** :
  - Notification à l'entité bénéficiaire (si applicable)
  - Mise à jour du cache de permissions
  - Audit log de la création

### Event `profileGrant.updated`
**Modification d'un grant**
- **Déclencheur** : Modification des GrantedFields ou ExpiresAt
- **Action** :
  - Notification à l'entité bénéficiaire des changements
  - Invalidation du cache de permissions
  - Audit log de la modification

### Event `profileGrant.revoked`
**Révocation d'un grant**
- **Déclencheur** : Revoked passe à true
- **Action** :
  - Notification à l'entité bénéficiaire
  - Invalidation immédiate du cache de permissions
  - Fermeture des accès en cours
  - Audit log de la révocation

### Event `profileGrant.expired`
**Expiration automatique d'un grant**
- **Déclencheur** : ExpiresAt dépassé (commande quotidienne)
- **Action** :
  - Marquer comme expiré (Revoked = true)
  - Notification à l'entité bénéficiaire
  - Invalidation du cache de permissions

### Event `offerResponse.accepted`
**Création automatique depuis OfferResponse**
- **Déclencheur** : OfferResponse avec status "accepted"
- **Action** :
  - Création d'un ProfileGrant avec scopeType "recruiter"
  - RecruiterID = celui de l'offre
  - GrantedFields = copie de OfferResponse.sharedFields

### Event `offer.deleted|archived|expired`
**Suppression automatique des ProfileGrant liés**
- **Déclencheur** : Offre supprimée, archivée ou expirée
- **Action** :
  - Recherche des ProfileGrant créés via OfferResponse pour cette offre
  - Suppression définitive des ProfileGrant associés
  - Notification aux candidats concernés de la révocation d'accès
  - Invalidation du cache de permissions (si cache)

## Cas d'usage typiques

### 1. Candidat créé un grant public
```json
{
  "scopeType": "audience",
  "audience": "public", 
  "grantedFields": ["email", "skills"]
}
```
→ Tous les utilisateurs authentifiés peuvent voir email et skills

### 2. Candidat donne accès à une entreprise spécifique
```json
{
  "scopeType": "company",
  "companyId": "company_123",
  "grantedFields": ["cv", "experience"],
  "expiresAt": "2024-06-30T23:59:59Z"
}
```
→ Tous les recruteurs de cette entreprise peuvent voir CV et expérience jusqu'au 30/06/2024

### 3. Grant automatique suite à acceptation d'offre
```json
{
  "scopeType": "recruiter", 
  "recruiterId": "recruiter_456",
  "grantedFields": ["email", "phone", "linkedin"]
}
```
→ Le recruteur spécifique peut voir les champs partagés lors de l'acceptation

## Commandes automatiques

### command `expireProfileGrants`
**Expiration automatique des grants**
- **Fréquence** : Tous les jours à 01:00
- **Action** :
  - Vérifie tous les ProfileGrant avec ExpiresAt dépassé
  - Met Revoked = true pour les grants expirés
  - Déclenche les événements d'expiration