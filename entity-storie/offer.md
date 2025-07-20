# Offer  

> Qui a accès aux offres ? 
> - Un recruteur authentifié
> - Un utilisateur authentifié qui a reçu l'offre en question

## Role "recruiter" avec SubRole "recruiter_admin"

- Champs authorisés : 
  - tous les champs de la struct Offer
- Lecture :  
  - peut voir toutes les offres créées par tout le monde
  - peut voir toutes les offres désactivées
  - peut voir toutes les offres archivées
  - peut voir tous les détails d'une seule offre
- Ecriture : 
  - peut créer une offre
  - peut modifier toutes les offres
  - peut désactiver toutes les offres
  - peut archiver toutes les offres
  - peut supprimer toutes les offres
  - peut transférer une offre à un autre recruteur

## Role "recruiter" avec SubRole "recruiter_user"

- Champs authorisés :
    - tous les champs de la struct Offer
- Lecture :
    - peut voir seulement ses offres
    - peut voir seulement ses offres désactivées
    - peut voir seulement ses offres archivées
    - peut voir tous les détails d'une seule offre
- Ecriture :
    - peut créer une offre
    - peut modifier seulement ses offres
    - peut désactiver seulement ses offres
    - peut archiver seulement ses offres

> La liste des offres s'affiche sous forme de tableau, avec les colonnes suivantes :  
> - Titre
> - Statut (active, expired, archivée, désactivée) se modifie en fonction des actions suivantes : (modifiable directement avec un select)
>   - Création (active)
>   - Désactivation (désactivée)
>   - Archivage (archivée)
>   - depassement de la date d'expiration (expired)
> - Date de création
> - Date d'expiration (modifiable directement)
> - Date StartDateJob (modifiable directement)
> - Quiz (oui / non) si "oui" cliquable direction vers le quiz si "non" cliquable pour en créer un
> - Actions (modifier, supprimer)

> Les détails d'une offre s'affiche sous forme de card, et contient toutes les informations de l'offre
> - Tous les champs de la struct Offer
> - la partie "Quiz" contient un lien vers le quiz s'il existe, sinon un lien pour accèder à la création d'un quiz
> - le nombre de jours que l'offre existe (calculé à partir de la date de création)
> - le nombre de jours avant la date d'expiration (calculé à partir de la date actuelle)
> - le nombre de jours avant la date de début de travail (calculé à partir de la date actuelle)
> - le nombre de condidats qui ont accepté l'offre (si l'offre est active)
> - le nombre de condidats qui ont refusé l'offre (si l'offre est active)
> - affiche une liste courtes des condidats qui ont répondu à l'offre (si l'offre est active)
>   - colonnes : nom prénom domaine 
>   - possibilité de cliquer sur un candidat pour voir son profil (en fonction des droits activés par le candidat)

## Comportement specifiques : 

- si l'annonce est expirée, archivée, supprimée, le candidat ne verra plus l'offre et tous les envois d'offres liés à cette offre seront supprimés

## Flows techniques - API

### GET /offers
**Récupération de la liste des offres**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toutes les offres (actives, désactivées, archivées)
  - `recruiter_user` : seulement ses propres offres
  - `candidate` : seulement si `OfferSend` existe pour cette utilisateur et que l'offre a le statut `active`
- **Query parameters** :
  - `status` : `active|disabled|archived` (optionnel)
  - `page` : pagination (optionnel)
  - `limit` : nombre d'éléments par page (optionnel)
- **Response** : `200` avec array d'offres
- **Errors** : `401`, `403` `500`

### GET /offers/{id}
**Récupération d'une offre spécifique**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute offre
  - `recruiter_user` : seulement ses propres offres
  - `candidate` : seulement si `OfferSend` existe pour cette utilisateur et que l'offre a le statut `active`
- **Response** : `200` avec détails complets de l'offre
- **Errors** : `401`, `403`, `404`

### POST /offers
**Création d'une nouvelle offre**
- **Authentification** : Token JWT requis
- **Autorisation** : 
  - Role `recruiter` requis
- **Regles** :
  - Toute les regles de `CreateOffer` 
  - la date d'expiration doit etre supérieure à la date actuelle et à la date de début de travail
  - la date de début de travail doit etre supérieure à la date actuelle et inférieure à la date d'expiration
  - il sera possible de créé un quiz seulement apres que l'offre soit créée
- **Body** : DTO `CreateOffer`
- **Response** : `201` avec l'offre créée
- **Errors** : `400`, `401`, `403`, `500`

### POST /offers/{id}
**Modification complète d'une offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute offre
  - `recruiter_user` : seulement ses propres offres
- **Regles** :
  - Toute les regles de `UpdateOffer`
  - la date d'expiration doit etre supérieure à la date actuelle et à la date de début de travail
  - la date de début de travail doit etre supérieure à la date actuelle et inférieure à la date d'expiration
- **Body** : DTO `UpdateOffer`
- **Response** : `200` avec l'offre mise à jour
- **Errors** : `400`, `401`, `403`, `404`, `500`

### POST /offers/{id}/enabled
**Activation d'une offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute offre
  - `recruiter_user` : seulement ses propres offres
- **Body** : pas de body
- **Response** : `200`
- **Errors** : `401`, `403`, `404`, `500`

### POST /offers/{id}/disabled
**Désactivation d'une offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute offre
  - `recruiter_user` : seulement ses propres offres
- **Body** : pas de body
- **Response** : `200`
- **Errors** : `401`, `403`, `404`, `500`

### POST /offers/{id}/archived
**Archivage d'une offre**
- **Authentification** : Token JWT requis
- **Autorisation** :
  - `recruiter_admin` : toute offre
  - `recruiter_user` : seulement ses propres offres
- **Body** : pas de body
- **Response** : `200`
- **Errors** : `401`, `403`, `404`, `500`

### POST /offers/{id}/delete
**Suppression définitive d'une offre**
- **Authentification** : Token JWT requis
- **Autorisation** : 
  - `recruiter_admin` uniquement
- **Body** : pas de body
- **Response** : `204`
- **Errors** : `401`, `403`, `404`, `500`

### POST /offers/{id}/transfer
**Transfert d'une offre à un autre recruteur**
- **Authentification** : Token JWT requis
- **Autorisation** : 
  - `recruiter_admin` uniquement
- **Regles** :
  - L'offre doit etre active
  - Le nouveau recruteur doit exister et être un utilisateur de la meme companie
- **Body** : `{ "newRecruiterId": "string" }`
- **Response** : `200`
- **Errors** : `400`, `401`, `403`, `404`, `500`

### Codes d'erreur communs
- `400 Bad Request` : Données invalides ou manquantes
- `401 Unauthorized` : Token manquant ou invalide
- `403 Forbidden` : Permissions insuffisantes
- `404 Not Found` : Ressource non trouvée
- `500 Internal Server Error` : Erreur serveur

## Flows techniques - COMMANDE

### command `checkDateExpiredOffers`
**Vérification des champs "expiredDate" des offres**
- **Fréquence** : Tout les jours à 00:00
- **Action** : 
  - Vérifie toutes les offres dont la date d'expiration est dépassée
  - Met à jour le statut de ces offres à `expired`

### command `checkDateStartDateJobOffers`
**Vérification des champs "startDateJob" des offres**
- **Fréquence** : Tout les jours à 00:00
- **Action** : 
  - Vérifie toutes les offres dont la date de début de travail est dépassée
  - Met à jour le statut de ces offres à `expired`

### Event `offer.deleted|archived|expired`
**Impact sur les ProfileGrant associés et Quiz**
- **Déclencheur** : Offre supprimée, archivée ou expirée
- **Action** :
  - Suppression automatique des ProfileGrant créés via OfferResponse pour cette offre
  - Quiz se referer au fhcier quiz.md
  - Notification aux candidats concernés de la révocation d'accès
  - Notification aux candidats ayant des soumissions de quiz en cours
  - Invalidation du cache de permissions ProfileGrant (si il y a)

### Event `offer.archived|expired`
**Impact sur les ProfileGrant associés et Quiz**
- **Déclencheur** : Offre supprimée, archivée ou expirée
- **Action** :
    - Suppression automatique des ProfileGrant créés via OfferResponse pour cette offre
    - Quiz se referer au fichier quiz.md
    - Notification aux candidats concernés de la révocation d'accès
    - Notification aux candidats ayant des soumissions de quiz en cours
    - Invalidation du cache de permissions ProfileGrant (si il y a)
