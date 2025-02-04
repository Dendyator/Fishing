-- migrations/apply_migrations.sql

CREATE TABLE IF NOT EXISTS fishing_methods (
                                               id SERIAL PRIMARY KEY,
                                               name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS fish_species (
                                            id SERIAL PRIMARY KEY,
                                            name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS baits (
                                     id SERIAL PRIMARY KEY,
                                     name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS fishing_data (
                                            id SERIAL PRIMARY KEY,
                                            date DATE NOT NULL,
                                            location TEXT,
                                            coordinates TEXT,
                                            comment TEXT,
                                            fishing_methods INTEGER[],
                                            caught_fish JSONB,
                                            trophy_fish JSONB
);