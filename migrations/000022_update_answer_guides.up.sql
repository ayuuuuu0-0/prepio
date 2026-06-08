-- Add explicit required concepts to each answer guide for concept-level evaluation.
UPDATE questions SET answer_guide = 'concepts:hash map|O(n) time|O(n) space|two sum|duplicate handling
expect hash map approach with O(n) time and O(n) space; mention duplicate handling and index return order'
WHERE id = 'b0000000-0000-4000-8000-000000000001';

UPDATE questions SET answer_guide = 'concepts:floyd|cycle detection|tortoise|hare|O(n) time|O(1) space
expect floyd tortoise and hare algorithm; discuss O(n) time and O(1) space'
WHERE id = 'b0000000-0000-4000-8000-000000000002';

UPDATE questions SET answer_guide = 'concepts:api design|id generation|storage|redirect|analytics|cache|rate limit|availability
cover api design, id generation, storage, redirect path, analytics pipeline, cache, rate limits, and availability tradeoffs'
WHERE id = 'b0000000-0000-4000-8000-000000000003';

UPDATE questions SET answer_guide = 'concepts:recursive|iterative|dfs|bfs|maximum depth|time complexity|space complexity
recursive or iterative dfs/bfs acceptable; state time and space complexity'
WHERE id = 'b0000000-0000-4000-8000-000000000004';

UPDATE questions SET answer_guide = 'concepts:parking lot|floor|spot|vehicle|ticket|pricing|concurrency
define entities parking lot floor spot vehicle ticket pricing strategy and concurrency considerations'
WHERE id = 'b0000000-0000-4000-8000-000000000005';

UPDATE questions SET answer_guide = 'concepts:sliding window|hash set|hash map|O(n) time|distinct characters
sliding window with hash set or map; O(n) time expected'
WHERE id = 'b0000000-0000-4000-8000-000000000006';

UPDATE questions SET answer_guide = 'concepts:5 minutes|rate|proportional|widgets|machines
answer is 5 minutes; explain rate reasoning not proportional trap'
WHERE id = 'b0000000-0000-4000-8000-000000000007';

UPDATE questions SET answer_guide = 'concepts:process|thread|isolation|memory|scheduling|overhead
cover isolation memory model scheduling overhead and when to use each'
WHERE id = 'b0000000-0000-4000-8000-000000000008';

UPDATE questions SET answer_guide = 'concepts:star format|situation|task|action|result|tradeoff
use star format; emphasize prioritization communication and measurable outcome'
WHERE id = 'b0000000-0000-4000-8000-000000000009';

UPDATE questions SET answer_guide = 'concepts:token bucket|leaky bucket|sliding window|redis|rate limit
compare token bucket leaky bucket sliding window; discuss redis implementation and failure modes'
WHERE id = 'b0000000-0000-4000-8000-000000000010';

UPDATE questions SET answer_guide = 'concepts:two heap|max heap|min heap|median|O(log n)
two heap approach max heap for lower half min heap for upper half; rebalance invariant'
WHERE id = 'b0000000-0000-4000-8000-000000000011';

UPDATE questions SET answer_guide = 'concepts:dispatch|driver allocation|eta|hotspot|degradation|peak hours
cover order matching driver allocation eta prediction hotspots and graceful degradation'
WHERE id = 'b0000000-0000-4000-8000-000000000012';
