-- Phase L3: the lesson platform replaces the daily-paper question loop.
-- Dependents first, then the tables they reference. Skills, subskills, user_skill_scores,
-- worlds, journey_nodes, and the lesson tables are kept.
DROP TABLE IF EXISTS daily_paper_questions;
DROP TABLE IF EXISTS daily_papers;
DROP TABLE IF EXISTS user_question_history;
DROP TABLE IF EXISTS pool_questions;
DROP TABLE IF EXISTS node_pools;
DROP TABLE IF EXISTS node_skills;
DROP TABLE IF EXISTS question_pools;
DROP TABLE IF EXISTS question_skills;
DROP TABLE IF EXISTS user_journey_progress;
DROP TABLE IF EXISTS questions;
