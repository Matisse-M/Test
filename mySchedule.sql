CREATE DATABASE IF NOT EXISTS mySchedule;

USE mySchedule;

CREATE TABLE IF NOT EXISTS employee (
    id INT AUTO_INCREMENT PRIMARY KEY,
    lastname VARCHAR(50) NOT NULL,
    firstname VARCHAR(50) NOT NULL,
    monday VARCHAR(50),
    tuesday VARCHAR(50),
    wednesday VARCHAR(50),
    thursday VARCHAR(50),
    friday VARCHAR(50),
    saturday VARCHAR(50),
    out_of_office BOOLEAN DEFAULT FALSE,
    sick BOOLEAN DEFAULT FALSE
);
