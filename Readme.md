---
Backend Authors : Dujardin Cyprien, Girault Dylan, Morin-Testa Emilien
Frontend Authors :  Jourdain Léo
---
## Projet AI30 - Simumation multi-agent d'une conquête spatiale
**Groupe apprentissage** : Dujardin Cyprien, Girault Dylan, Jourdain Léo, Morin-Testa Emilien

### Requirement
- nodejs
- yarn [(<u>installation</u>)](https://chore-update--yarnpkg.netlify.app/fr/docs/install)
- golang

### Installation 
```bash
git clone
cd frontend && yarn

# OU

cd source_code && cd frontend && yarn
```

### Lancement de la simulation
```bash
# Lancement du backend
cd backend/cmd && go run .

# Lancement du server web (Frontend)
cd frontend && yarn build && yarn preview
```
Le démarrage de la simulation se fait dans le navigateur, on choisi les ressources des pays (`shuffle` pour mettre des ressources aléatoires à chaque pays), puis on click sur `start`

### Enum des stratégies
#### Stratégies d'observation :
- `quick` : le pays observe la première planète de sa liste, s'il a le budget pour y aller, alors il prépare le lancement vers la planète 
- `slow` : le pays observe toutes les planètes une fois, ensuite il observe une nouvelle fois la planète la plus rentable. Si une planète a été observée au moins 2 fois, et qu'elle est la plus rentable, alors il prépare le lancement vers la planète
- `if_lo_quick` : si le pays à moins de 20 000 de budget, alors il opte pour la stratégie `quick`, sinon il suivra la stratégie `slow`
- `if_lo_slow` : si le pays à moins de 20 000 de budget, alors il opte pour la stratégie `slow`, sinon il suivra la stratégie `quick`
- `if_hi_quick` : si le pays à plus de 50 000 de budget, alors il opte pour la stratégie `quick`, sinon il suivra la stratégie `slow`
- `if_hi_slow` : si le pays à plus de 50 000 de budget, alors il opte pour la stratégie `slow`, sinon il suivra la stratégie `quick`

#### Stratégies de lancement (launch) :
- `quick` : le pays prépare l'équipage pour avoir environ 50% de chance de **réussir** la mission
- `medium` : le pays prépare l'équipage pour avoir environ 75% de chance de **réussir** la mission
- `slow` : le pays prépare l'équipage pour avoir environ 95% de chance de **réussir** la mission
- `if_lo_quick` : si le pays un budget **inférieur 2 fois** au coût de la mission, alors il opte pour la stratégie `quick`, sinon il suivra la stratégie `slow`
- `if_lo_medium` : si le pays un budget **inférieur 2 fois** au coût de la mission, alors il opte pour la stratégie `medium`, sinon il suivra la stratégie `slow`
- `if_lo_slow` : si le pays un budget **inférieur 2 fois** au coût de la mission, alors il opte pour la stratégie `slow`, sinon il suivra la stratégie `quick`
- `if_hi_quick` : si le pays un budget **supérieur à 4 fois** le coût de la mission, alors il opte pour la stratégie `quick`, sinon il suivra la stratégie `slow`
- `if_hi_medium` : si le pays un budget **supérieur à 4 fois** le coût de la mission, alors il opte pour la stratégie `medium`, sinon il suivra la stratégie `slow`
- `if_hi_slow` : si le pays un budget **supérieur à 4 fois** le coût de la mission, alors il opte pour la stratégie `slow`, sinon il suivra la stratégie `quick`

#### Stratégies d'abandon (skipping) :
- `always` : Si un pays est déjà sur la planète visée, ou arrive sur la planète pendant que l'on prépare la mission, on abandonne la mission
- `never` : Tant qu'il reste des ressources sur la planète, on peut préparer l'équipage pour y envoyer une fusée
- `if_hi` : Si on a un budget au moins **supérieur à 4 fois** le coût de la mission, alors si une personne est sur la planète, on abandonne la mission
- `if_lo` : Si on a un budget au moins **inférieur à 2 fois** le coût de la mission, alors si une personne est sur la planète, on abandonne la mission

#### Stratégies de réaction :
- `faster` : Si un pays est déjà sur la planète que l'on vise, alors on **diminue** le taux de préparation de l'équipage de 10%
- `slower` : Si un pays est déjà sur la planète que l'on vise, alors on **augmente** le taux de préparation de l'équipage de 10%
- `same` : Si un pays est déjà sur la planète que l'on vise, on ne change rien
- `balanced_lo` : Si un pays est déjà sur la planète que l'on vise, dans le cas ou on est sur une stratégie de lancement conditionnelle (`if_.._..`) alors on **augmente** la préparation de l'équipage si la condition est satisfaite, sinon on la **diminue**.
- `balanced_hi` : Si un pays est déjà sur la planète que l'on vise, dans le cas ou on est sur une stratégie de lancement conditionnelle (`if_.._..`) alors on **diminue** la préparation de l'équipage si la condition est satisfaite, sinon on **l'augmente**.
 

