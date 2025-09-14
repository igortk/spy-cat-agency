CREATE TABLE missions (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          cat_id UUID UNIQUE REFERENCES cats(id),
                          is_completed BOOLEAN NOT NULL DEFAULT FALSE
);