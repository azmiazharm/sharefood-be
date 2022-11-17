-- +goose Up
-- +goose StatementBegin
-- please delete this migration cause for example purpose
CREATE TABLE IF NOT EXISTS example (
    id INT UNSIGNED AUTO_INCREMENT,
    name VARCHAR (100) DEFAULT NULL,
    email VARCHAR (50) NOT NULL,
    phone VARCHAR (20) DEFAULT '' ,
    address TEXT DEFAULT '',
    PRIMARY KEY id (id),
    UNIQUE email(email)

)ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
-- +goose StatementEnd
INSERT INTO example (name, email, phone, address) VALUES ("jhon doe", "jhon.doe@mail.com","0821111110","jl merdeka raya");

-- +goose Down
-- +goose StatementBegin
SELECT 'do nothing';
-- +goose StatementEnd
