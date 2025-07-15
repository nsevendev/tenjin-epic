# **User Stories & Workflows**

Voici les **principaux parcours d’utilisation** (user stories) couvrant tout le périmètre métier :

---

### A. **Création d’une formation (manager/admin)**

1. L’administrateur d’un organisme crée une nouvelle formation.
2. Il recherche un métier dans le catalogue France Travail (via JobSkillCatalog).
3. Il sélectionne les compétences associées automatiquement ou en ajoute d’autres si besoin.
4. Il sélectionne/associe les cours (modules) existants ou en crée de nouveaux.
5. Il enregistre la formation, qui devient disponible pour la création de sessions.

---

### B. **Création d’une session de formation**

1. Le manager/gestionnaire choisit une formation et crée une session (dates, formateurs, etc.).
2. Il sélectionne les cours/modules (et peut réordonner/ajouter/supprimer pour cette session).
3. Il affecte les stagiaires à la session (import, inscription…).
4. Il ajoute les ressources (PDF, vidéos, etc.), les quiz et crée les salons de discussion (général, privé…).
5. Le calendrier de la session est renseigné (dates de cours, examens, etc.).

---

### C. **Suivi des présences**

1. À chaque session, le formateur/membre staff marque la présence/absence/retard pour chaque stagiaire.
2. Les feuilles de présence sont consultables et exportables.

---

### D. **Gestion du carnet de compétences**

1. Lorsqu’un stagiaire suit un cours/session, il passe des évaluations (quiz, devoirs…).
2. Le formateur valide les compétences atteintes avec un niveau (note, palier, etc.), selon la grille de l’organisme ou celle du métier.
3. L’historique de toutes les validations, niveaux, dates et organismes est stocké sur le profil de l’utilisateur.

---

### E. **Ajout d’une expérience externe**

1. Un utilisateur ajoute une formation suivie hors de la plateforme (titre, description, date…).
2. Il télécharge une ou plusieurs preuves (PDF, diplômes…).
3. L’admin de l’organisme peut, s’il le souhaite, valider ou compléter la compétence associée à cette expérience.

---

### F. **Recherche de profils par un recruteur**

1. Un recruteur (ayant validé son SIRET) accède à la recherche avancée.
2. Il filtre les profils sur la base des compétences, niveaux, statut, etc.
3. Pour chaque profil trouvé, il ne voit que le prénom, le nom, les compétences et leur niveau/validation.
4. Il sélectionne un ou plusieurs candidats pour leur envoyer une offre d’emploi (texte ou PDF, quiz éventuel).

---

### G. **Envoi et gestion d’une offre par le recruteur**

1. Le recruteur compose une offre (message personnalisé, fichier joint éventuel, quiz éventuel).
2. Il sélectionne un ou plusieurs candidats ciblés.
3. Les utilisateurs reçoivent une notification.
4. L’utilisateur peut accepter ou refuser l’offre et choisir les informations personnelles à partager (ex : téléphone, email, CV…).
5. Si l’utilisateur accepte, un salon de discussion privé est ouvert entre le recruteur et lui.
   → Ce salon est auto-supprimé au bout d’un mois d’inactivité.
6. Le recruteur peut suivre l’état de chaque offre (envoyée, vue, acceptée, refusée…).

---

### H. **Quiz côté recruteur**

1. Le recruteur crée un quiz (QCM, texte libre, fichier…).
2. Il l’associe à une offre ou l’envoie indépendamment à des candidats sélectionnés.
3. Les candidats répondent ; le recruteur visualise les résultats dans son espace.

---

### I. **Chat et collaboration**

1. Les sessions ont des salons dédiés (général, modéré, etc.), visibles selon le rôle.
2. Les discussions privées entre utilisateurs et recruteurs sont créées uniquement après consentement et supprimées après 1 mois sans réponse.
3. Les formateurs et gestionnaires peuvent modérer les salons de session.

---

### J. **Synchronisation du catalogue France Travail**

1. Un service autonome synchronise régulièrement le catalogue France Travail (métiers, fiches, compétences…).
2. Lorsqu’une formation ou une compétence est créée ou modifiée, le système propose d’associer à une fiche/compétence officielle.
3. Les formations et compétences restent toujours alignées avec le référentiel national.

---

### K. **Export/Valorisation du carnet de compétences**

1. Un utilisateur peut exporter à tout moment son carnet de compétences, avec tout l’historique des validations et preuves.
2. Ce document est valorisable pour un recruteur, un organisme externe, ou pour les démarches administratives.

---

### L. **Gestion de la confidentialité**

1. Lorsqu’un recruteur envoie une offre, l’utilisateur garde le contrôle sur les infos personnelles à partager.
2. Rien n’est visible au recruteur sans consentement explicite.
3. Un utilisateur peut révoquer à tout moment un partage d’info ou fermer un chat privé.
