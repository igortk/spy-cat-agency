CREATE TABLE cats (
                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                      name VARCHAR(255) NOT NULL,
                      years_of_experience INTEGER NOT NULL,
                      breed VARCHAR(255) NOT NULL,
                      salary DECIMAL(10, 2) NOT NULL,
                      status VARCHAR(50) NOT NULL DEFAULT 'Available',
                      CONSTRAINT chk_years_of_experience CHECK (years_of_experience >= 0),
                      CONSTRAINT chk_salary CHECK (salary >= 0)
);