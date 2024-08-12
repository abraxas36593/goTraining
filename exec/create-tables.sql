DROP TABLE IF EXISTS library;
CREATE TABLE library (
  id      INT AUTO_INCREMENT NOT NULL,
  title   VARCHAR(128) NOT NULL,
  author  VARCHAR(255) NOT NULL,
  checked INT,
  PRIMARY KEY (`id`)
);

INSERT INTO library
  (title, author, checked)
VALUES
  ('Magick Without Tears', 'Aleister Crowley', 1),
  ('Paradise Lost', 'John Milton', 0),
  ('Liber Null And Psychonaut', 'Peter Carroll', 1),
  ('The Earthsea Trilogy', 'Ursela Laguin', 0),
  ('1984', 'George Orwell', 1),
  ('Neuromancer', 'William Gibson', 0),
  ('The Valis Trilogy', 'Phillip K Dick', 1);
