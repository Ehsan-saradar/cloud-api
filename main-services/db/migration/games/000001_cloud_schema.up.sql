CREATE TABLE "games" (
                            "id" serial PRIMARY KEY,
                            "rank" bigint NOT NULL,
                            "name" varchar NOT NULL,
                            "platform" varchar NOT NULL,
                            "year" bigint NOT NULL,
                            "genre" varchar NOT NULL,
                            "publisher" varchar NOT NULL,
                            "na_sale" real NOT NULL,
                            "eu_sale" real NOT NULL,
                            "jp_sale" real NOT NULL,
                            "other_sale" real NOT NULL,
                            "global_sale" real NOT NULL
);