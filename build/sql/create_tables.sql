CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    description TEXT
        CONSTRAINT description_length CHECK (char_length(username) <= 255),
    username TEXT
        NOT NULL
        UNIQUE
        CONSTRAINT name_length CHECK (char_length(username) <= 255),
    password_hash TEXT
        NOT NULL
        CONSTRAINT password_hash_length CHECK (char_length(password_hash) <= 511),
    create_time TIMESTAMP
        NOT NULL,
    image_path TEXT DEFAULT ('default.jpg')
        NOT NULL
        CONSTRAINT image_path_length CHECK (char_length(image_path) <= 255)
);

CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY,
    data JSON,
    create_time TIMESTAMP
        NOT NULL,
    update_time TIMESTAMP,
    owner_id UUID REFERENCES users (id)
        NOT NULL
);


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION add_draft_note()
    RETURNS trigger
    LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
    name_text json;
BEGIN
    IF NEW.id IS NOT NULL AND NEW.username IS NOT NULL THEN
        name_text := format('{
            "title": "You-note ❤️",
            "content": "Привет, %s!"
        }', NEW.username);
        INSERT INTO notes (id, data, create_time, update_time, owner_id)
        VALUES (uuid_generate_v4(), name_text, CURRENT_TIMESTAMP, NULL, NEW.id);
    END IF;

    RETURN NEW;
END;
$BODY$;

CREATE OR REPLACE TRIGGER add_note_on_new_user
    AFTER INSERT
    ON users
    FOR EACH ROW
    EXECUTE FUNCTION add_draft_note();


INSERT INTO users (id, description, username, password_hash, create_time)
    VALUES ('3d67ba58-9023-42b5-a059-d626d7587f1e', 'У меня много заметок!', 'alladan', 'f969248d621bcded4a3582a1c3b17a71eedfefa9120c36ee3bd1957438cd55b9', CURRENT_TIMESTAMP);
INSERT INTO notes (id, data, create_time, update_time, owner_id)
    VALUES ('f732e6ae-07fe-4f20-aca1-37e8f115d082', '{
            "title": "Список покупок 😇",
            "content": "Полка для овощей и фруктов:\nКартофель\nМорковь\nЛук\nЧеснок\nПетрушка\n\nУкроп\nЯблоки/бананы\nЛимон\n\nПолка для молочных продуктов:\nМасло сливочное\n\nКефир\nМолоко детское (7 упаковок)\nСметана\nТворог\nСыр"
        }', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '3d67ba58-9023-42b5-a059-d626d7587f1e'),
            ('81e769de-4cd3-499b-90a7-98b4aca6f9d2', '{
            "title": "Вино из одуванчиков ☀️ (Рэй Брэдбери)",
            "content": "Нет! Дуглас приказал себе ни о чем другом не думать. Нельзя! А потом вдруг… Можно. Почему нет? Да! Потасовка, соударение тел, падение на землю не отпугнуло нахлынувшее море, которое все затопило и вынесло их на травянистый берег в чащу леса. Костяшки пальцев стукнули его по зубам, он ощутил во рту ржавый теплый привкус крови. Он крепко-накрепко сграбастал Тома, и они лежали в тишине. Сердечки колотились, ноздри сопели. И наконец, медленно, опасаясь, что ничего не обнаружится, Дуглас приоткрыл один глаз.\nИ все, все-превсе оказалось на своем месте.\nПодобно огромной радужке гигантского глаза, который тоже только что раскрылся и раздался вширь, чтобы вобрать в себя все на свете, на него уставилась Вселенная.\nИ он понял: то, что нахлынуло на него, останется с ним навсегда и никуда больше не сбежит.\n«Я – живу», – подумал он."
        }', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '3d67ba58-9023-42b5-a059-d626d7587f1e'),
            ('6e246734-e712-4182-ad59-c5cd8debff8b', '{
            "title": "Испанский 🇪🇸",
            "content": "Si estás estudiando español, tenemos buenas noticias para ti. Debes saber que estás aprendiendo la segunda lengua del mundo con el mayor\nnúmero de hablantes nativos. Y si a esto sumamos a todos aquellos que, como tú, estudian español, el resultado está cerca de los 600 millones de personas. Recuerda que la mayoría del continente americano, una parte de Europa y una pequeña parte de África usa la lengua española a diario. Es decir, se trata de un idioma universal."
        }', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '3d67ba58-9023-42b5-a059-d626d7587f1e');

INSERT INTO users (id, description, username, password_hash, create_time)
    VALUES ('57247dd2-b768-4665-a290-8dad9506616a', 'Заметок не делаю, всё держу в голове!', 'mizhgun', 'f969248d621bcded4a3582a1c3b17a71eedfefa9120c36ee3bd1957438cd55b9', CURRENT_TIMESTAMP);
DELETE FROM notes WHERE owner_id = '57247dd2-b768-4665-a290-8dad9506616a';

INSERT INTO users (id, description, username, password_hash, create_time)
VALUES ('17248dd2-b768-4665-a290-8dad9506616a', 'アニメが好きです', 'japanman', 'f969248d621bcded4a3582a1c3b17a71eedfefa9120c36ee3bd1957438cd55b9', CURRENT_TIMESTAMP);
INSERT INTO notes (id, data, create_time, update_time, owner_id)
VALUES ('a732e7ae-07fe-4f20-aca1-37e8f115d082', '{
  "title": "私は本物の日本人です 🥷",
  "content": "- 山田さん、どこへ行きますか。\n- きっぷを買いに行きます。\n- コンサ- トのチケットですか。\n\n- いいえ、ひこうきのチケットです。\n来週ドイツへ行きます。\n\n- かんこうに行きますか。\n- ちがいます。\n- 友だちに会いに行きますか。\n- ちがいます。\n- 何をしにいきますか。\n- べんきょうに行きます。\n- 何のべんきょうに行きますか。\n- 医学のべんきょうに行きます。\n- りゅう学をしますか。\n\n - そうです。\nベルリン医科大学にりゅう学したいです。\nよい医者になりたいです。\n\n- ぼくも若いときロンドンへりゅう学くしました。\nほかの国のりゅう学生も多かったです。\nたくさんの友だちができました。\n先生がたはしんせつでべんきょうは面白かったです。\nいい思い出がたくさんあります。\n休みに国へかえりますか。\n\n- もちろん、りょうしんに会いに帰ります。\n冬休みに帰りませんが、夏やすみに帰ります。\n\n- そうですか。\nじゅぎょうは何月からですか。\n\n- ９月からです。\nこれからしけんをうけに行きます。\n\n- じゃ、がんばってくださいね。\n- がんばります。"
}', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '17248dd2-b768-4665-a290-8dad9506616a');
