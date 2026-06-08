DROP DATABASE IF EXISTS RomPelago;

CREATE DATABASE RomPelago
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE RomPelago;

CREATE TABLE Utilisateur (
    id_utilisateur INT AUTO_INCREMENT PRIMARY KEY,
    pseudo         VARCHAR(20)  NOT NULL UNIQUE,
    e_mail         VARCHAR(254) NOT NULL UNIQUE,
    description    VARCHAR(100),
    ban            BOOLEAN      NOT NULL DEFAULT FALSE,
    status         ENUM('user', 'admin') NOT NULL DEFAULT 'user'
    --explication de ENUM: https://dev.mysql.com/doc/refman/8.4/en/enum.html
    --Le type ENUM permet de stocker une valeur unique choisie parmi une liste prédéfinie. C’est très pratique pour des champs comme le statut d’une commande, le type de compte ou la catégorie d’un produit.
) ENGINE=InnoDB;

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

CREATE TABLE Fil_de_discussion (
    id_fil_de_discussion INT         AUTO_INCREMENT PRIMARY KEY,
    fk_utilisateur       INT         NOT NULL,
    name                 VARCHAR(50) NOT NULL,
    description          TEXT,
    open                 ENUM('ouvert', 'ferme', 'archive') NOT NULL DEFAULT 'ouvert',

    CONSTRAINT fk_fil_utilisateur
        FOREIGN KEY (fk_utilisateur)
        REFERENCES Utilisateur(id_utilisateur)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE Type_de_discussion (
    id_type INT         AUTO_INCREMENT PRIMARY KEY,
    name    VARCHAR(50) NOT NULL UNIQUE
) ENGINE=InnoDB;

-- Cette table fait le lien entre un fil et ses catégories, car un fil peut avoir plusieurs catégories et une catégorie peut appartenir à plusieurs fils.
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
        REFERENCES Type_de_discussion(id_type)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE Message (
    id_message           INT      AUTO_INCREMENT PRIMARY KEY,
    fk_fil_de_discussion INT      NOT NULL,
    fk_utilisateur       INT      NOT NULL,
    contenu              TEXT     NOT NULL,
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

CREATE TABLE Score_de_popularite (
    id_score   INT AUTO_INCREMENT PRIMARY KEY,
    fk_message INT NOT NULL UNIQUE,
    nb_like    INT NOT NULL DEFAULT 0,
    nb_dislike INT NOT NULL DEFAULT 0,
    score      INT NOT NULL DEFAULT 0,

    CONSTRAINT fk_score_message
        FOREIGN KEY (fk_message)
        REFERENCES Message(id_message)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE Reaction (
    id_reaction    INT  AUTO_INCREMENT PRIMARY KEY,
    fk_message     INT  NOT NULL,
    fk_utilisateur INT  NOT NULL,
    type_reaction  ENUM('like', 'dislike') NOT NULL,

    UNIQUE KEY unique_reaction (fk_message, fk_utilisateur),

    CONSTRAINT fk_reaction_message
        FOREIGN KEY (fk_message)
        REFERENCES Message(id_message)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_reaction_utilisateur
        FOREIGN KEY (fk_utilisateur)
        REFERENCES Utilisateur(id_utilisateur)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;
