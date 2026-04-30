CREATE TABLE site_content (
    `key`      VARCHAR(100) PRIMARY KEY,
    value      TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO site_content (`key`, value) VALUES
('home_intro', 'Hark Horning is a world-renowned artist, software engineer, and premier dirt specialist.'),
('about_content', 'Hark Horning is a world-renowned artist, software engineer, and premier dirt specialist.\n\nHis paintings work across oil, acrylic, watercolor, pastel, and pencil. The subjects vary — what stays consistent is a concern for light, color, form, and composition.\n\nHe has, at various points, tested dirt professionally and taught people to climb walls. Both felt like the right thing to do at the time.\n\nOriginal works and limited prints are available through this site. For inquiries, reach out at howpretty73@gmail.com.');
