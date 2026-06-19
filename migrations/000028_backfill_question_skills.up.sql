-- Enrich seed questions with user-facing metadata and explicit skill mappings.

-- Q1: Two Sum — arrays + hash maps
UPDATE questions SET
    explanation = 'Use a hash map to store each value and its index while scanning. For each element, check if target minus current exists in the map.',
    hints = '[{"order": 1, "text": "Can you avoid checking every pair?"}, {"order": 2, "text": "Store seen values in a hash map as you iterate."}]'::jsonb,
    solution = 'Iterate nums with index i. If target - nums[i] is in the map, return [map[target-nums[i]], i]. Otherwise store nums[i] -> i. O(n) time, O(n) space.'
WHERE id = 'b0000000-0000-4000-8000-000000000001';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000001', 'b2000001-0000-4000-8000-000000000002', 'c3000001-0000-4000-8000-000000000005', 0.600),
    ('b0000000-0000-4000-8000-000000000001', 'b2000001-0000-4000-8000-000000000004', 'c3000001-0000-4000-8000-000000000011', 0.400);

-- Q2: Linked list cycle — linked lists
UPDATE questions SET
    explanation = 'Use Floyd''s tortoise and hare: slow moves one step, fast moves two. If they meet, a cycle exists.',
    hints = '[{"order": 1, "text": "Two pointers at different speeds can detect cycles."}]'::jsonb,
    solution = 'Initialize slow and fast at head. While fast and fast.next exist, advance slow by 1 and fast by 2. Return true if they meet. O(n) time, O(1) space.'
WHERE id = 'b0000000-0000-4000-8000-000000000002';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000002', 'b2000001-0000-4000-8000-000000000005', 'c3000006-0000-4000-8000-000000000002', 1.000);

-- Q3: URL shortener — system design fundamentals
UPDATE questions SET
    explanation = 'Design a scalable redirect service with unique short codes, durable storage, analytics, and caching for hot URLs.',
    hints = '[{"order": 1, "text": "Separate read-heavy redirect path from write path."}, {"order": 2, "text": "Consider base62 encoding for short codes."}]'::jsonb,
    solution = 'API: POST /shorten, GET /{code}. Use distributed ID generation or hash+salt. Store mappings in SQL/NoSQL with cache layer. Track click analytics asynchronously. Rate limit creation.'
WHERE id = 'b0000000-0000-4000-8000-000000000003';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000003', 'b2000001-0000-4000-8000-000000000016', 'c3000003-0000-4000-8000-000000000002', 0.500),
    ('b0000000-0000-4000-8000-000000000003', 'b2000001-0000-4000-8000-000000000016', 'c3000003-0000-4000-8000-000000000001', 0.500);

-- Q4: Max depth binary tree — trees DFS
UPDATE questions SET
    explanation = 'Recursively compute 1 + max depth of left and right subtrees, or use iterative BFS level counting.',
    hints = '[{"order": 1, "text": "Base case: null node has depth 0."}]'::jsonb,
    solution = 'Recursive: if root is nil return 0; return 1 + max(maxDepth(left), maxDepth(right)). O(n) time, O(h) stack space.'
WHERE id = 'b0000000-0000-4000-8000-000000000004';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000004', 'b2000001-0000-4000-8000-000000000008', 'c3000001-0000-4000-8000-000000000013', 1.000);

-- Q5: Parking lot LLD
UPDATE questions SET
    explanation = 'Model floors, spots by vehicle type, tickets, and entry/exit with pricing strategy. Consider concurrency for spot allocation.',
    hints = '[{"order": 1, "text": "Identify core entities first: Vehicle, Spot, Ticket, ParkingLot."}]'::jsonb,
    solution = 'Classes: ParkingLot, Floor, Spot (type, occupied), Vehicle (type), Ticket (entryTime, spot). Strategy pattern for pricing. Synchronized spot allocation or optimistic locking.'
WHERE id = 'b0000000-0000-4000-8000-000000000005';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000005', 'b2000001-0000-4000-8000-000000000020', 'c3000006-0000-4000-8000-000000000010', 1.000);

-- Q6: Longest substring without repeating — strings sliding window
UPDATE questions SET
    explanation = 'Expand window with a hash set/map tracking last seen index. Shrink from left when duplicate found.',
    hints = '[{"order": 1, "text": "Sliding window over the string with a frequency map."}]'::jsonb,
    solution = 'Use map char->index. For each char, if seen within window move left pointer. Track max window size. O(n) time.'
WHERE id = 'b0000000-0000-4000-8000-000000000006';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000006', 'b2000001-0000-4000-8000-000000000003', 'c3000001-0000-4000-8000-000000000009', 1.000);

-- Q7: Aptitude widgets — problem solving
UPDATE questions SET
    explanation = 'Each machine makes 1 widget per minute regardless of count. 100 machines still take 5 minutes for 100 widgets.',
    hints = '[{"order": 1, "text": "Calculate per-machine rate, not total proportional scaling."}]'::jsonb,
    solution = '5 machines → 5 widgets in 5 min means 1 widget/machine/min. 100 machines → 100 widgets in 5 minutes.'
WHERE id = 'b0000000-0000-4000-8000-000000000007';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000007', 'b2000001-0000-4000-8000-000000000028', 'c3000006-0000-4000-8000-000000000013', 1.000);

-- Q8: Process vs thread — programming fundamentals
UPDATE questions SET
    explanation = 'Processes have isolated memory; threads share address space within a process. Threads are lighter but need synchronization.',
    hints = '[{"order": 1, "text": "Compare isolation, memory model, and scheduling overhead."}]'::jsonb,
    solution = 'Process: own memory, heavier context switch, crash isolated. Thread: shared heap, lighter switch, needs locks. Use threads for I/O-bound concurrency within a service.'
WHERE id = 'b0000000-0000-4000-8000-000000000008';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000008', 'b2000001-0000-4000-8000-000000000001', 'c3000006-0000-4000-8000-000000000001', 1.000);

-- Q9: Behavioral deadline — STAR
UPDATE questions SET
    explanation = 'Structure with STAR: describe situation, your task, actions taken, and measurable result. Emphasize prioritization and communication.',
    hints = '[{"order": 1, "text": "Use STAR format with a quantifiable outcome."}]'::jsonb,
    solution = 'Situation: tight release deadline. Task: deliver critical feature. Action: prioritized scope, communicated tradeoffs, parallelized testing. Result: shipped on time with X% defect reduction.'
WHERE id = 'b0000000-0000-4000-8000-000000000009';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000009', 'b2000001-0000-4000-8000-000000000023', 'c3000004-0000-4000-8000-000000000003', 1.000);

-- Q10: Rate limiter — system design scaling
UPDATE questions SET
    explanation = 'Compare token bucket, leaky bucket, and sliding window. Use Redis for distributed counters with TTL.',
    hints = '[{"order": 1, "text": "Consider per-user and per-IP keys separately."}]'::jsonb,
    solution = 'Sliding window log or token bucket in Redis. Key: rl:{scope}:{id}. Return 429 when exceeded. Handle Redis failure with fail-open or local fallback.'
WHERE id = 'b0000000-0000-4000-8000-000000000010';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000010', 'b2000001-0000-4000-8000-000000000017', 'c3000003-0000-4000-8000-000000000004', 0.500),
    ('b0000000-0000-4000-8000-000000000010', 'b2000001-0000-4000-8000-000000000016', 'c3000003-0000-4000-8000-000000000002', 0.500);

-- Q11: Median stream — heaps
UPDATE questions SET
    explanation = 'Maintain max-heap for lower half and min-heap for upper half. Rebalance so sizes differ by at most 1.',
    hints = '[{"order": 1, "text": "Two heaps keep the median at the top of one or both heaps."}]'::jsonb,
    solution = 'On insert, add to appropriate heap and rebalance. Median is top of max-heap (odd count) or average of both tops (even). O(log n) per insertion.'
WHERE id = 'b0000000-0000-4000-8000-000000000011';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000011', 'b2000001-0000-4000-8000-000000000010', 'c3000006-0000-4000-8000-000000000005', 1.000);

-- Q12: Food delivery dispatch — system design scaling
UPDATE questions SET
    explanation = 'Design order matching, driver allocation, ETA prediction, hotspot handling, and graceful degradation during peak load.',
    hints = '[{"order": 1, "text": "Separate hot path (dispatch) from analytics."}]'::jsonb,
    solution = 'Services: Order, Dispatch, Location, ETA. Geospatial index for drivers. Queue-based matching with priority for peak zones. Cache ETAs. Degrade to longer wait estimates under load.'
WHERE id = 'b0000000-0000-4000-8000-000000000012';

INSERT INTO question_skills (question_id, skill_id, subskill_id, weight) VALUES
    ('b0000000-0000-4000-8000-000000000012', 'b2000001-0000-4000-8000-000000000017', 'c3000003-0000-4000-8000-000000000006', 0.500),
    ('b0000000-0000-4000-8000-000000000012', 'b2000001-0000-4000-8000-000000000016', 'c3000003-0000-4000-8000-000000000001', 0.500);
