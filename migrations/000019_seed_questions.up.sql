INSERT INTO questions (id, body, round_type, difficulty, answer_guide, status, is_weekend, source)
VALUES
    (
        'b0000000-0000-4000-8000-000000000001',
        'Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.',
        'dsa',
        'easy',
        'expect hash map approach with O(n) time and O(n) space; mention duplicate handling and index return order',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000002',
        'Implement a function to detect if a linked list has a cycle.',
        'dsa',
        'easy',
        'expect floyd tortoise and hare algorithm; discuss O(n) time and O(1) space',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000003',
        'Design a URL shortening service like bit.ly that supports custom aliases and analytics.',
        'system_design',
        'medium',
        'cover api design, id generation, storage, redirect path, analytics pipeline, cache, rate limits, and availability tradeoffs',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000004',
        'Given a binary tree, return its maximum depth.',
        'dsa',
        'easy',
        'recursive or iterative dfs/bfs acceptable; state time and space complexity',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000005',
        'Design a parking lot system with multiple floors and vehicle types.',
        'lld',
        'medium',
        'define entities parking lot floor spot vehicle ticket pricing strategy and concurrency considerations',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000006',
        'Find the longest substring without repeating characters.',
        'dsa',
        'medium',
        'sliding window with hash set or map; O(n) time expected',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000007',
        'If 5 machines can make 5 widgets in 5 minutes, how long would 100 machines take to make 100 widgets?',
        'aptitude',
        'easy',
        'answer is 5 minutes; explain rate reasoning not proportional trap',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000008',
        'Explain the difference between process and thread with examples relevant to backend systems.',
        'fundamentals',
        'easy',
        'cover isolation memory model scheduling overhead and when to use each',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000009',
        'Tell me about a time you had to deliver under a tight deadline. What tradeoffs did you make?',
        'behavioral',
        'medium',
        'use star format; emphasize prioritization communication and measurable outcome',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000010',
        'Design a rate limiter for an API gateway supporting per-user and per-ip limits.',
        'system_design',
        'hard',
        'compare token bucket leaky bucket sliding window; discuss redis implementation and failure modes',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000011',
        'Given a stream of integers, find the median at any point in O(log n) per insertion.',
        'dsa',
        'hard',
        'two heap approach max heap for lower half min heap for upper half; rebalance invariant',
        'approved',
        false,
        'manual'
    ),
    (
        'b0000000-0000-4000-8000-000000000012',
        'Design a food delivery dispatch system for peak dinner hours.',
        'system_design',
        'hard',
        'cover order matching driver allocation eta prediction hotspots and graceful degradation',
        'approved',
        true,
        'manual'
    );

INSERT INTO question_tags (question_id, company)
VALUES
    ('b0000000-0000-4000-8000-000000000001', 'google'),
    ('b0000000-0000-4000-8000-000000000001', 'amazon'),
    ('b0000000-0000-4000-8000-000000000002', 'meta'),
    ('b0000000-0000-4000-8000-000000000003', 'google'),
    ('b0000000-0000-4000-8000-000000000003', 'microsoft'),
    ('b0000000-0000-4000-8000-000000000004', 'amazon'),
    ('b0000000-0000-4000-8000-000000000005', 'zepto'),
    ('b0000000-0000-4000-8000-000000000006', 'google'),
    ('b0000000-0000-4000-8000-000000000006', 'meta'),
    ('b0000000-0000-4000-8000-000000000007', 'flipkart'),
    ('b0000000-0000-4000-8000-000000000008', 'microsoft'),
    ('b0000000-0000-4000-8000-000000000009', 'amazon'),
    ('b0000000-0000-4000-8000-000000000010', 'google'),
    ('b0000000-0000-4000-8000-000000000011', 'meta'),
    ('b0000000-0000-4000-8000-000000000012', 'zepto'),
    ('b0000000-0000-4000-8000-000000000012', 'swiggy');
