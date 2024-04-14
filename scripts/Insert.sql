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
INSERT INTO banners_data (id, content) VALUES
(101, 'FIRST BANNER'),
(102, 'SECOND BANNER'),
(103, 'THIRD BANNER');
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