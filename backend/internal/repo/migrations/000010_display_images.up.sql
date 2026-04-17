CREATE TABLE art_display_images (
    art_tile_id INT NOT NULL,
    image_id    INT NOT NULL,
    PRIMARY KEY (art_tile_id, image_id),
    FOREIGN KEY (art_tile_id) REFERENCES art_tiles(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id)    REFERENCES images(id)    ON DELETE CASCADE
);

CREATE TABLE print_display_images (
    print_id INT NOT NULL,
    image_id INT NOT NULL,
    PRIMARY KEY (print_id, image_id),
    FOREIGN KEY (print_id) REFERENCES prints(id)  ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES images(id)  ON DELETE CASCADE
);
