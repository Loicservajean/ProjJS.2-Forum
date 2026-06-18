DROP DATABASE IF EXISTS RomPelago;

CREATE DATABASE RomPelago
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE RomPelago;

CREATE TABLE Fil_de_discussion (
    id_fil_de_discussion INT         AUTO_INCREMENT PRIMARY KEY,
    name                 VARCHAR(50) NOT NULL,
    description          TEXT,
    date_creation		 DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fk_utilisateur       INT NULL
) ENGINE=InnoDB;

CREATE TABLE Utilisateur (
    id_utilisateur INT AUTO_INCREMENT PRIMARY KEY,
    fk_fil_de_discussion INT NULL,
    pseudo         VARCHAR(20)  NOT NULL UNIQUE,
    e_mail         VARCHAR(254) NOT NULL UNIQUE,
    description    VARCHAR(100),
    ban            BOOLEAN      NOT NULL DEFAULT FALSE,
    status         ENUM('user', 'admin') NOT NULL DEFAULT 'user',
# explication de ENUM: https://dev.mysql.com/doc/refman/8.4/en/enum.html
# Le type ENUM permet de stocker une valeur unique choisie parmi une liste prédéfinie. C'est très pratique pour des champs comme le statut d'une commande, le type de compte ou la catégorie d'un produit.
    
    CONSTRAINT fk_utilisateur_fil_de_discussion
        FOREIGN KEY (fk_fil_de_discussion)
        REFERENCES Fil_de_discussion(id_fil_de_discussion)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

ALTER TABLE Fil_de_discussion
    ADD CONSTRAINT fk_fil_createur
        FOREIGN KEY (fk_utilisateur)
        REFERENCES Utilisateur(id_utilisateur)
        ON DELETE SET NULL
        ON UPDATE CASCADE;

CREATE TABLE Mot_de_passe (
    id_mot_de_passe INT          AUTO_INCREMENT PRIMARY KEY,
    fk_utilisateur  INT          NOT NULL UNIQUE,
    mot_de_passe    VARCHAR(128) NOT NULL,
    majuscule       BOOLEAN      NOT NULL DEFAULT FALSE,
    carac_spe       VARCHAR(10),
    nb_carac        INT          NOT NULL DEFAULT 0,

    CONSTRAINT fk_mdp_utilisateur
        FOREIGN KEY (fk_utilisateur)
        REFERENCES Utilisateur(id_utilisateur)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;


CREATE TABLE CategoriesDiscussion (
    id_type INT         AUTO_INCREMENT PRIMARY KEY,
    name    VARCHAR(50) NOT NULL UNIQUE,
    Description VARCHAR(100)
) ENGINE=InnoDB;

CREATE TABLE Status (
    id_status INT         AUTO_INCREMENT PRIMARY KEY,
    name    VARCHAR(50) NOT NULL UNIQUE,
    Description VARCHAR(100)
) ENGINE=InnoDB;

CREATE TABLE Fil_status (
    fk_fil  INT NOT NULL,
    fk_status INT NOT NULL,

    PRIMARY KEY (fk_fil, fk_status),

    CONSTRAINT fk_filstatus_fil
        FOREIGN KEY (fk_fil)
        REFERENCES Fil_de_discussion(id_fil_de_discussion)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_filstatus_status
        FOREIGN KEY (fk_status)
        REFERENCES Status(id_status)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;


CREATE TABLE Tag (
    id_tag INT         AUTO_INCREMENT PRIMARY KEY,
    name    VARCHAR(50) NOT NULL UNIQUE,
    Description VARCHAR(100)
) ENGINE=InnoDB;

CREATE TABLE Fil_Tag (
    fk_fil  INT NOT NULL,
    fk_tag INT NOT NULL,

    PRIMARY KEY (fk_fil, fk_tag),

    CONSTRAINT fk_filtag_fil
        FOREIGN KEY (fk_fil)
        REFERENCES Fil_de_discussion(id_fil_de_discussion)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_filtag_tag
        FOREIGN KEY (fk_tag)
        REFERENCES Tag(id_tag)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

# Cette table fait le lien entre un fil et ses catégories, car un fil peut avoir plusieurs catégories et une catégorie peut appartenir à plusieurs fils.
CREATE TABLE Fil_Type (
    fk_fil  INT NOT NULL,
    fk_type INT NOT NULL,

    PRIMARY KEY (fk_fil, fk_type),

    CONSTRAINT fk_filtype_fil
        FOREIGN KEY (fk_fil)
        REFERENCES Fil_de_discussion(id_fil_de_discussion)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_filtype_type
        FOREIGN KEY (fk_type)
        REFERENCES CategoriesDiscussion (id_type)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE Message (
    id_message           INT      AUTO_INCREMENT PRIMARY KEY,
    fk_fil_de_discussion INT      NOT NULL,
    fk_utilisateur       INT      NOT NULL,
    contenu              TEXT     NOT NULL,
    name VARCHAR(50) NOT NULL ,
    nb_like    INT NOT NULL DEFAULT 0,
    nb_dislike INT NOT NULL DEFAULT 0,
    scorepop      INT NOT NULL DEFAULT 0,
    date_envoi           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_message_fil
        FOREIGN KEY (fk_fil_de_discussion)
        REFERENCES Fil_de_discussion(id_fil_de_discussion)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_message_utilisateur
        FOREIGN KEY (fk_utilisateur)
        REFERENCES Utilisateur(id_utilisateur)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE LikeDislike (
    fk_utilisateur INT NOT NULL,
    fk_message INT NOT NULL,
    type_vote VARCHAR(10) NOT NULL, 
    PRIMARY KEY (fk_utilisateur, fk_message),
    FOREIGN KEY (fk_utilisateur) REFERENCES Utilisateur(id_utilisateur),
    FOREIGN KEY (fk_message) REFERENCES Message(id_message)
) ENGINE=InnoDB;
INSERT INTO Status (name, description) VALUES ("Ouvert", "Fil disponible et réponses disponibles");
INSERT INTO Status (name, description) VALUES ("Fermé", "Fil disponible et réponses indisponibles");
INSERT INTO Status (name, description) VALUES ("Archivé", "Fil indisponible à part pour le créateur");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("games", "Discussions autour des jeux vidéo");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("Archipelago links", "Liens vers le site de l'Archipelago");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("Core-Vérified games", "Jeux officiellement supportés et intégrés par Archipelago");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("New Games", "Propositions de nouveaux jeux");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("FAQ", "Foire Aux Questions");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("Games with Energy link", "Jeux concernés par le système Energy Link");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("Helping", "Comment puis-je vous aider ?");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("Général topic", "Sujets généraux");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("proposed modification", "Propositions de modifications");
INSERT INTO CategoriesDiscussion (name, description) VALUES ("Uncategorized categories", "Catégories non classées");
