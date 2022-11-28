/*Matt Martin 
initial tables for social media project
*/


DROP TABLE IF EXISTS post;
CREATE TABLE post (
    postid    INT AUTO_INCREMENT NOT NULL,
    userid    INTEGER NOT NULL,
    tstamp TIMESTAMP,
    txt      VARCHAR(300),
    PRIMARY KEY (`postid`)
);



DROP TABLE IF EXISTS SocUser;
CREATE TABLE SocUser (
    userid      INT AUTO_INCREMENT NOT NULL,
 --   pword       VARCHAR(255),
    fname       VARCHAR(255),
    lname       VARCHAR(255),
    displayname VARCHAR(255),
    email       VARCHAR(255),
    picurl    VARCHAR(255),
    PRIMARY KEY (`userid`)
);

DROP TABLE IF EXISTS SecurePass;
CREATE TABLE SecurePass (
    userid      INT NOT NULL,
    pword       VARCHAR(255),
    PRIMARY KEY (`userid`)
);


DROP TABLE IF EXISTS FriendList;
CREATE TABLE FriendList (
    userid      INT NOT NULL,
    friendid    INT NOT NULL  
);

ALTER TABLE post
    ADD CONSTRAINT post_user_fk FOREIGN KEY ( userid )
        REFERENCES SocUser ( userid );

ALTER TABLE SecurePass
    ADD CONSTRAINT user_pass_fk FOREIGN KEY ( userid )
        REFERENCES SocUser ( userid );



INSERT INTO SocUser
    (fname, lname, displayname, email, picurl)
VALUES
('Matt', 'Martin', 'Jug', 'matrmart@iu.edu', 'https://espinvr.com/wp-content/uploads/2020/10/Matt.png'),
('Michaelina', 'Magnuson', 'MickyFree', 'Mdwewe@msu.edu', 'https://espinvr.com/wp-content/uploads/2020/10/Micky.png'),
('Autumn', 'Ramirez', 'AutoMir', 'AutMonkey@gmail.com','https://espinvr.com/wp-content/uploads/2020/10/Autumn.png'),
('Ruben', 'Saldivar', 'RuTheMan', 'RuTheMan@gmail.com','https://espinvr.com/wp-content/uploads/2022/11/Ruben.jpg'),
('Jonathan', 'Zook', 'ReverendHowitzer', 'BivBiy@gmail.com','https://espinvr.com/wp-content/uploads/2020/10/Jonathan.png'),
('Christian', 'Palmer', 'DuffelBlog', 'LtColDuffelBlog@gmail.com','https://espinvr.com/wp-content/uploads/2022/11/Palmer.jpg'),
('Alex', 'Jones', 'Samsquatch', 'Beardy@gmail.com','https://espinvr.com/wp-content/uploads/2020/10/Alex.png');

INSERT INTO SecurePass
    (userid, pword)
VALUES
(1,'Password1'),
(2,'Password2'),
(3,'Password3'),
(4,'Password4'),
(5,'Password5'),
(6,'Password6'),
(7,'Password7');

INSERT INTO post
    (userid, tstamp, txt)
VALUES
(1, '2022-11-24 22:05:00', 'What would you trade me for my Margarita?'),
(1, '2022-11-23 13:15:00', 'Nerf''s Battleaxe has to be the greatest invention ever.'),
(1, '2022-11-20 13:15:00', 'Wow... This is just great! -- Vlad D. when given his first shishkebab'),
(1, '2022-11-19 16:25:00', 'So we get in the truck to take the monkeys home when Autumn say She forgot her keys.  She ran frantically back inside while we waited. Many minutes later she calls exasperarted from the house phone.." Dad... I can''t find them... oh wait, I''m wearing them!". You''d think a lanyard would have helped... '),
(2, '2022-11-23 19:25:00',  'hope to see everyone Friday night!  Looks like I''ll be working the pow-wow stand before fitting in some good belly laughs from Don Burnstick!'),
(2, '2022-11-21 09:15:00', 'Spent the weekend with my sweetheart, his extended family, my sisters, and my mother.  I am blessed!  Hope everyone else had as great a day.'),
(2, '2022-11-20 19:15:00', 'Got my butt kicked in Scrabble. My ego is shot! lol'),
(2, '2022-11-18 17:15:00', '2.5 hours for winter tires, but at least I have them now. Going to find some slushy roads to break them in!'),
(3, '2022-11-24 22:15:00', 'the sadness you feel after finishing the series xc ugh...'),
(3, '2022-11-23 20:45:00', 'perks of being sick around Thanksgiving: I DON''T HAVE TO HELP COOK ANYTHING xD lol'),
(3, '2022-11-22 12:55:00', 'you know what? I got pringles and Usher I''m set for the rest of my life xD'),
(4, '2022-11-24 11:15:00', 'My wife just said the nicest thing ever.... I asked her if she wanted me to loose weight and be all muscle... She said " No because then other women will look at you...."'),
(4, '2022-11-19 14:30:00', 'Anyone know about a cage fight in Dowagiac on the 27th??? On Rudy rd I think...'),
(4, '2022-11-17 02:15:00', 'Accident 94w near hartford careful'),
(5, '2022-11-23 16:50:00', 'Today at work we''ve been playing that fun game, "How long will it take my co-worker to give me the flu?"  You can still see the skidmarks on the wall from where he coughed up a lung multiple times today.'),
(5, '2022-11-18 17:30:00', 'I have good news and bad news.  The good news is that last night at my company''s party my left leg popped back into my hip.  The bad news is that I no longer walk like Long John Silver.'),
(5, '2022-11-15 18:15:00', 'Tonight''s Movie at our regular open-house is John Carter, starring the Princess of Mars.  Yes, Virginia, there is a Princess of Mars.'),
(6, '2022-11-24 11:15:00', 'So, I just watched the Hobbit.  I think you could watch this same movie and assume that this is like the tenth group of Dwarves Gandalf''s led to their deaths in like a year.'),
(6, '2022-11-23 19:00:00', 'I wonder if anyone''s ever invoked Gresham''s law to explain why we lose so many of our good Marines to the private sector.'),
(6, '2022-11-12 12:15:00', 'Denial, anger...just two stages of grieving left before the fringe can join the rest of us in planning for the midterms.'),
(7, '2022-11-23 23:15:00', 'When life get''s ya down, make a comforter...'),
(7, '2022-11-20 07:45:00', 'Everything is perfect about the past except that it led to the present - Homer Simpson.'),
(7, '2022-11-14 21:55:00', 'As of next week, passwords will be entered in Morse code!'),
(1, '2022-11-26 23:27:59', 'Here I go working on a group project...'),
(7, '2022-11-26 23:39:20', 'Honk if you like bigfoots... bigfeet... like the plural of saquatches n stuff.');

INSERT INTO FriendList
    (userid, friendid)
VALUES
    (1,2),
    (1,3),
    (1,4),
    (1,7),
    (2,1),
    (2,5),
    (2,6),
    (3,1),
    (3,5),
    (3,6),
    (3,7),
    (4,1),
    (4,5),
    (4,6),
    (4,7),
    (5,2),
    (5,3),
    (5,4),
    (6,2),
    (6,3),
    (6,4),
    (7,1),
    (7,3),
    (7,4);
