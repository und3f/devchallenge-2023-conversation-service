--liquibase formatted sql

--changeset V3_Sample_Call:insertSampleCall
INSERT INTO calls(processed, name, location, emotional_tone, text)
OVERRIDING SYSTEM VALUE
VALUES (
  TRUE,
  'Sample Call',
  'Kyiv',
  'Neutral',
  'Hello and welcome to out call in Kyiv. I am happy to talk about visa and diplomatic inquries!'
);

--changeset V3_Sample_Call:insertSampleCallCategories
INSERT INTO call_categories(call_id, category_id)
VALUES
  (
    (SELECT id FROM calls WHERE name = 'Sample Call'),
    (SELECT id FROM categories WHERE title = 'Visa and Passport Services')
  ),
  (
    (SELECT id FROM calls WHERE name = 'Sample Call'),
    (SELECT id FROM categories WHERE title = 'Diplomatic Inquiries')
  );
