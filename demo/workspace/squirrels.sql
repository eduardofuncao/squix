-- Create squirrels database schema with expanded tables
DROP TABLE IF EXISTS squirrel_cache_locations;
DROP TABLE IF EXISTS squirrel_family_tree;
DROP TABLE IF EXISTS predators;
DROP TABLE IF EXISTS squirrel_diets;
DROP TABLE IF EXISTS sightings;
DROP TABLE IF EXISTS squirrels;

CREATE TABLE squirrels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    species TEXT NOT NULL,
    age_years INTEGER,
    favorite_food TEXT,
    park_location TEXT,
    weight_kg REAL,
    health_status TEXT,
    tag_id TEXT UNIQUE
);

CREATE TABLE sightings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    squirrel_id INTEGER NOT NULL,
    date TEXT NOT NULL,
    behavior TEXT,
    notes TEXT,
    temperature_c REAL,
    weather_condition TEXT,
    observer_name TEXT,
    FOREIGN KEY (squirrel_id) REFERENCES squirrels(id)
);

CREATE TABLE squirrel_diets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    food_name TEXT NOT NULL UNIQUE,
    energy_per_100g REAL,
    protein_g REAL,
    fat_g REAL,
    carbohydrates_g REAL,
    description TEXT
);

CREATE TABLE predators (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    species TEXT NOT NULL,
    danger_level TEXT,
    primary_hunting_time TEXT,
    description TEXT
);

CREATE TABLE squirrel_cache_locations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    squirrel_id INTEGER NOT NULL,
    location_name TEXT NOT NULL,
    cached_food TEXT,
    cache_date TEXT,
    retrieval_date TEXT,
    notes TEXT,
    FOREIGN KEY (squirrel_id) REFERENCES squirrels(id)
);

CREATE TABLE squirrel_family_tree (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER,
    child_id INTEGER NOT NULL,
    relationship_type TEXT,
    FOREIGN KEY (parent_id) REFERENCES squirrels(id),
    FOREIGN KEY (child_id) REFERENCES squirrels(id)
);

-- Insert sample squirrels with expanded data
INSERT INTO squirrels (name, species, age_years, favorite_food, park_location, weight_kg, health_status, tag_id) VALUES
    ('Nutty', 'Eastern Gray Squirrel', 3, 'acorns', 'Central Park North', 0.45, 'Healthy', 'SQ001'),
    ('Bushy', 'Fox Squirrel', 2, 'walnuts', 'Riverside Park', 0.52, 'Healthy', 'SQ002'),
    ('Squeaky', 'Red Squirrel', 1, 'pine seeds', 'Prospect Park', 0.28, 'Underweight', 'SQ003'),
    ('Shadow', 'Eastern Gray Squirrel', 4, 'peanuts', 'Battery Park', 0.48, 'Healthy', 'SQ004'),
    ('Acorn', 'Eastern Gray Squirrel', 2, 'hazelnuts', 'Golden Gate Park', 0.41, 'Healthy', 'SQ005'),
    ('Rusty', 'Fox Squirrel', 3, 'corn', 'Lincoln Park', 0.55, 'Healthy', 'SQ006'),
    ('Peanut', 'Eastern Gray Squirrel', 1, 'almonds', 'Boston Common', 0.32, 'Recovering', 'SQ007'),
    ('Chippy', 'American Red Squirrel', 2, 'sunflower seeds', 'Grant Park', 0.30, 'Healthy', 'SQ008'),
    ('Swift', 'Eastern Gray Squirrel', 3, 'bird seed', 'Bryant Park', 0.44, 'Healthy', 'SQ009'),
    ('Bandit', 'Fox Squirrel', 4, 'pizza crust', 'Millennium Park', 0.58, 'Overweight', 'SQ010'),
    ('Snowball', 'Eastern Gray Squirrel', 2, 'pecans', 'Central Park South', 0.43, 'Healthy', 'SQ011'),
    ('Ginger', 'Red Squirrel', 1, 'beechnuts', 'Riverside Park', 0.26, 'Healthy', 'SQ012'),
    ('Rocky', 'Fox Squirrel', 3, 'hickory nuts', 'Golden Gate Park', 0.51, 'Healthy', 'SQ013');

-- Insert sample sightings with expanded data
INSERT INTO sightings (squirrel_id, date, behavior, notes, temperature_c, weather_condition, observer_name) VALUES
    (1, '2025-01-15', 'foraging', 'Found three large acorns', 5.5, 'Sunny', 'Park Ranger Jane'),
    (1, '2025-02-10', 'nesting', 'Building nest in oak tree', 8.2, 'Cloudy', 'Visitor Tom'),
    (2, '2025-01-20', 'foraging', 'Buried walnuts near bench', 12.0, 'Sunny', 'Bird Watcher Bill'),
    (3, '2025-03-01', 'chattering', 'Warning call at hawk', 15.5, 'Clear', 'Naturalist Lisa'),
    (4, '2025-01-25', 'bathing', 'Dusting in dirt path', 10.0, 'Partly Cloudy', 'Runner Rob'),
    (5, '2025-02-14', 'foraging', 'Stole hazelnut from tourist', 18.0, 'Sunny', 'Tourist Maria'),
    (6, '2025-03-05', 'resting', 'Napping on park bench', 20.0, 'Warm', 'Jogger Jake'),
    (7, '2025-01-30', 'begging', 'Following joggers', 7.5, 'Overcast', 'Dog Walker Dan'),
    (8, '2025-02-20', 'foraging', 'Collecting seeds', 14.0, 'Sunny', 'Photographer Pam'),
    (9, '2025-03-10', 'playing', 'Chasing Blue Jay', 17.5, 'Clear', 'Student Sam'),
    (10, '2025-01-18', 'eating', 'Found discarded pizza', 22.0, 'Hot', 'Urban Explorer Ursula'),
    (1, '2025-03-15', 'running', 'Racing with other squirrel', 19.0, 'Sunny', 'Child Chris'),
    (2, '2025-02-28', 'foraging', 'Found bird feeder stash', 16.0, 'Partly Cloudy', 'Home Owner Hannah'),
    (5, '2025-03-08', 'grooming', 'Cleaning tail', 21.0, 'Warm', 'Bicyclist Bob'),
    (7, '2025-02-05', 'climbing', 'Scaling small maple tree', 11.0, 'Breezy', 'Tree Climber Tim'),
    (10, '2025-02-25', 'sleeping', 'Napping on branch', 13.0, 'Cloudy', 'Nature Writer Nora'),
    (3, '2025-03-12', 'digging', 'Recovering cached food', 9.0, 'Cool', 'Researcher Rachel'),
    (4, '2025-01-22', 'exploring', 'Checking trash cans', 6.0, 'Foggy', 'Sanitation Worker Steve'),
    (6, '2025-03-02', 'chasing', 'Pursuing another squirrel', 23.0, 'Hot', 'Playground Parent Paula'),
    (9, '2025-02-18', 'drinking', 'From puddle', 12.0, 'Rainy', 'Rain Walker Rick'),
    (11, '2025-03-20', 'mating', 'Courting display', 25.0, 'Sunny', 'Biologist Brian'),
    (12, '2025-02-22', 'storing', 'Caching beechnuts', 14.0, 'Overcast', 'Forager Fred'),
    (13, '2025-03-01', 'territorial', 'Chasing intruder', 19.0, 'Clear', 'Territory Tracker Ted'),
    (1, '2025-03-25', 'teaching', 'Showing kit how to crack nuts', 21.0, 'Warm', 'Wildlife Watcher Wendy');

-- Insert squirrel diet information
INSERT INTO squirrel_diets (food_name, energy_per_100g, protein_g, fat_g, carbohydrates_g, description) VALUES
    ('acorns', 387, 8, 24, 40, 'High in tannins, staple for gray squirrels'),
    ('walnuts', 654, 15, 65, 14, 'High fat content, prized food source'),
    ('pine seeds', 673, 14, 68, 13, 'Extremely high fat, red squirrel favorite'),
    ('peanuts', 567, 26, 49, 16, 'Not actually a nut, but loved by squirrels'),
    ('hazelnuts', 628, 15, 61, 17, 'Rich in vitamin E and healthy fats'),
    ('corn', 86, 3, 1, 19, 'Common urban food source'),
    ('almonds', 579, 21, 50, 22, 'High protein content'),
    ('sunflower seeds', 584, 21, 51, 20, 'Very popular, easy to open'),
    ('bird seed', 450, 12, 18, 60, 'Mixed seeds from bird feeders'),
    ('pecans', 691, 9, 72, 14, 'Highest fat content of common nuts'),
    ('beechnuts', 576, 7, 50, 33, 'Small but energy-dense'),
    ('hickory nuts', 657, 12, 64, 18, 'Very hard shell, high reward'),
    ('pizza crust', 265, 9, 9, 41, 'Urban junk food, high in carbs');

-- Insert predator information
INSERT INTO predators (name, species, danger_level, primary_hunting_time, description) VALUES
    ('Red-tailed Hawk', 'Buteo jamaicensis', 'High', 'Daytime', 'Large hawk that hunts from above'),
    ('Cooper''s Hawk', 'Accipiter cooperii', 'High', 'Daytime', 'Agile hawk that navigates through trees'),
    ('Great Horned Owl', 'Bubo virginianus', 'High', 'Nighttime', 'Powerful nocturnal predator'),
    ('Domestic Cat', 'Felis catus', 'Medium', 'Day/Night', 'Urban and suburban threat'),
    ('Red Fox', 'Vulpes vulpes', 'Medium', 'Dawn/Dusk', 'Cunning canine predator'),
    ('Eastern Rat Snake', 'Pantherophis alleghaniensis', 'Medium', 'Daytime', 'Climbs trees to raid nests'),
    ('Barred Owl', 'Strix varia', 'Medium', 'Nighttime', 'Smaller owl, still a threat'),
    ('Feral Dog', 'Canis lupus familiaris', 'Low-Medium', 'Daytime', 'Occasional threat in parks'),
    ('Blue Jay', 'Cyanocitta cristata', 'Low', 'Daytime', 'May raid nests, harass squirrels'),
    ('American Crow', 'Corvus brachyrhynchos', 'Low', 'Daytime', 'Mobbing behavior, nest raider');

-- Insert cache location data
INSERT INTO squirrel_cache_locations (squirrel_id, location_name, cached_food, cache_date, retrieval_date, notes) VALUES
    (1, 'Under oak tree #45', 'acorns', '2025-01-16', '2025-02-15', 'Retrieved after 30 days'),
    (1, 'North garden bed', 'walnuts', '2025-01-25', NULL, 'Not yet retrieved'),
    (2, 'Behind bench #12', 'mixed nuts', '2025-02-01', '2025-02-20', 'Some stolen by jay'),
    (3, 'Pine tree hollow', 'pine seeds', '2025-02-10', '2025-03-05', 'Multiple caches nearby'),
    (4, 'Flower pot #7', 'peanuts', '2025-01-28', '2025-02-18', 'Urban caching behavior'),
    (5, 'Near Japanese garden', 'hazelnuts', '2025-02-15', NULL, 'Strategic location'),
    (6, 'Under maple #33', 'corn cobs', '2025-02-20', '2025-03-10', 'Buried deep'),
    (7, 'Playground sandbox', 'almonds', '2025-02-05', '2025-02-28', 'Forgot location initially'),
    (8, 'Bird feeder area', 'sunflower seeds', '2025-02-12', '2025-03-01', 'Temporary scatter hoard'),
    (9, 'Bryant Park fountain', 'bird seed', '2025-02-25', NULL, 'High traffic area'),
    (10, 'Near cafe patio', 'pizza crust', '2025-01-19', '2025-01-20', 'Quick stash'),
    (11, 'Central Park sheep meadow', 'pecans', '2025-03-01', NULL, 'Open area cache'),
    (12, 'Riverside path edge', 'beechnuts', '2025-02-18', '2025-03-08', 'Red squirrel style'),
    (13, 'Conservatory edge', 'hickory nuts', '2025-02-28', NULL, 'Hard shell, protected'),
    (1, 'Rock outcropping', 'acorns', '2025-03-10', NULL, 'Backup cache');

-- Insert family tree data
INSERT INTO squirrel_family_tree (parent_id, child_id, relationship_type) VALUES
    (NULL, 1, 'unknown'),
    (NULL, 2, 'unknown'),
    (NULL, 3, 'unknown'),
    (NULL, 4, 'unknown'),
    (1, 7, 'parent'),
    (2, 10, 'parent'),
    (3, 12, 'parent'),
    (4, 11, 'parent'),
    (5, 9, 'sibling'),
    (6, 13, 'sibling'),
    (8, 12, 'sibling'),
    (1, 7, 'parent'),
    (7, 14, 'parent'),
    (2, 10, 'parent'),
    (10, 15, 'parent');
