-- +goose Up
-- +goose StatementBegin
CREATE TABLE locations (
    slug TEXT PRIMARY KEY,
    name_de TEXT NOT NULL,
    name_en TEXT NOT NULL,
    provider TEXT NOT NULL
);

CREATE TABLE menus (
    slug TEXT PRIMARY KEY,
    date DATE NOT NULL,
    provider TEXT NOT NULL,
    location_slug TEXT NOT NULL,
    FOREIGN KEY(location_slug) REFERENCES locations(slug)
);

CREATE TABLE meals (
    id TEXT PRIMARY KEY,
    name_de TEXT NOT NULL,
    name_en TEXT NOT NULL,
    category_de TEXT NOT NULL,
    category_en TEXT NOT NULL,
    date DATE NOT NULL,
    menu_slug TEXT NOT NULL,
    location_slug TEXT NOT NULL,
    FOREIGN KEY(menu_slug) REFERENCES menus(slug)
    FOREIGN KEY(location_slug) REFERENCES locations(slug)
);

CREATE TABLE prices (
    id TEXT NOT NULL,
    type_de TEXT NOT NULL,
    type_en TEXT NOT NULL,
    amount FLOAT NOT NULL,
    meal_id TEXT NOT NULL,
    FOREIGN KEY(meal_id) REFERENCES meals(id)
);

CREATE TABLE attributes (
    slug TEXT PRIMARY KEY,
    type TEXT CHECK( type IN ('allergen', 'additive', 'feature') ) NOT NULL,
    name_de TEXT NOT NULL,
    name_en TEXT NOT NULL,
    short_de TEXT NOT NULL,
    short_en TEXT NOT NULL
);

CREATE TABLE meal_attributes (
    meal_id TEXT NOT NULL,
    attribute_slug TEXT NOT NULL,
    FOREIGN KEY(meal_id) REFERENCES meals(id)
    FOREIGN KEY(attribute_slug) REFERENCES attributes(slug)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE locations;
DROP TABLE menus;
DROP TABLE meals;
DROP TABLE menus;
DROP TABLE prices;
DROP TABLE meal_attributes;
-- +goose StatementEnd
