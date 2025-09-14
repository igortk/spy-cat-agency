CREATE TABLE targets (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         mission_id UUID REFERENCES missions(id) ON DELETE CASCADE,
                         name VARCHAR(255) NOT NULL,
                         country VARCHAR(255) NOT NULL,
                         notes JSONB NOT NULL DEFAULT '[]',
                         is_completed BOOLEAN NOT NULL DEFAULT FALSE,
                         CONSTRAINT unique_target_name_mission UNIQUE (mission_id, name)
);