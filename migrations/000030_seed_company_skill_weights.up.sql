-- Company skill weight profiles (sum = 100 per company).
-- Google
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'google', id, 15 FROM skills WHERE slug = 'arrays';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'google', id, 15 FROM skills WHERE slug = 'trees';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'google', id, 15 FROM skills WHERE slug = 'graphs';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'google', id, 20 FROM skills WHERE slug = 'dynamic-programming';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'google', id, 15 FROM skills WHERE slug = 'system-design-scaling';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'google', id, 10 FROM skills WHERE slug = 'behavioral-star';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'google', id, 10 FROM skills WHERE slug = 'communication-clarity';

-- Amazon
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'amazon', id, 15 FROM skills WHERE slug = 'arrays';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'amazon', id, 20 FROM skills WHERE slug = 'behavioral-leadership';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'amazon', id, 15 FROM skills WHERE slug = 'behavioral-star';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'amazon', id, 15 FROM skills WHERE slug = 'lld-patterns';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'amazon', id, 20 FROM skills WHERE slug = 'problem-solving';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'amazon', id, 15 FROM skills WHERE slug = 'system-design-data';

-- Meta
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'meta', id, 15 FROM skills WHERE slug = 'arrays';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'meta', id, 15 FROM skills WHERE slug = 'dynamic-programming';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'meta', id, 20 FROM skills WHERE slug = 'system-design-scaling';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'meta', id, 15 FROM skills WHERE slug = 'behavioral-star';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'meta', id, 15 FROM skills WHERE slug = 'communication-structure';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'meta', id, 20 FROM skills WHERE slug = 'graphs';

-- Uber
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'uber', id, 25 FROM skills WHERE slug = 'system-design-scaling';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'uber', id, 20 FROM skills WHERE slug = 'graphs';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'uber', id, 15 FROM skills WHERE slug = 'arrays';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'uber', id, 15 FROM skills WHERE slug = 'behavioral-star';
INSERT INTO company_skill_weights (company, skill_id, weight)
SELECT 'uber', id, 25 FROM skills WHERE slug = 'problem-solving';
