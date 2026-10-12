ALTER TABLE user_demographics
    ALTER COLUMN gender TYPE SMALLINT USING gender::smallint,
    ALTER COLUMN year_of_birth TYPE INTEGER USING year_of_birth::integer,
    ALTER COLUMN prefecture TYPE varying(10) USING prefecture::varying(10);
