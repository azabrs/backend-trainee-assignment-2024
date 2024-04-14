INSERT INTO tags (name) VALUES 
('red'), 
('blue'), 
('green'),
('yellow'),
('orange'),
('purple'),
('pink'),
('brown'),
('black'),
('white'),
('gray'),
('cyan'),
('magenta'),
('turquoise'),
('lime');
INSERT INTO banners_data (id, title, text_content, url_content) VALUES
(101, 'FIRST title', 'FIRST BANNER', 'FIRST URL'),
(102, 'SECOND title', 'SECOND BANNER', 'SECOND URL'),
(103, 'THIRD title', 'THIRD BANNER', 'THIRD URL');
INSERT INTO banners (feature_id, data_id, is_active) VALUES 
(1, 101, true),
(2, 102, false),
(2, 103, true);
INSERT INTO banner_tags (banner_id, tag_id) VALUES 
(1, 1), 
(1, 2), 
(2, 1), 
(3, 3), 
(3, 2);