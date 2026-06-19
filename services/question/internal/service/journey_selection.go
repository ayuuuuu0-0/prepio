package service

import (
	"context"
	"hash/fnv"
	"math/rand/v2"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/store"
)

// nodeAssignment holds question IDs selected for a journey node.
type nodeAssignment struct {
	questionIDs []string
	fromPool    bool
}

// resolveNodeAssignment picks questions for a node via pools or index fallback.
func (s *QuestionService) resolveNodeAssignment(
	ctx context.Context,
	poolEnabled bool,
	userID, sessionID string,
	node store.JourneyNode,
	nodeIndex int,
	paper *dto.DailyPaperResponse,
	seen map[string]bool,
) nodeAssignment {
	if !poolEnabled {
		return indexAssignment(node, nodeIndex, paper)
	}

	pools, err := s.content.ListNodePools(ctx, node.ID)
	if err != nil || len(pools) == 0 {
		return indexAssignment(node, nodeIndex, paper)
	}

	seed := selectionSeed(userID, sessionID, node.ID)
	if isBossMixed(pools) {
		ids := selectBossMixed(ctx, s.content, pools, seen, seed)
		if len(ids) == 0 {
			return indexAssignment(node, nodeIndex, paper)
		}
		return nodeAssignment{questionIDs: ids, fromPool: true}
	}

	pool := pools[0]
	poolQuestionIDs, err := s.content.ListPoolQuestionIDs(ctx, pool.PoolID)
	if err != nil || len(poolQuestionIDs) == 0 {
		return indexAssignment(node, nodeIndex, paper)
	}

	selected, ok := selectFromPool(pool.SelectionStrategy, poolQuestionIDs, seen, seed)
	if !ok {
		return indexAssignment(node, nodeIndex, paper)
	}
	return nodeAssignment{questionIDs: []string{selected}, fromPool: true}
}

func indexAssignment(node store.JourneyNode, nodeIndex int, paper *dto.DailyPaperResponse) nodeAssignment {
	if nodeIndex < len(paper.Questions) {
		return nodeAssignment{questionIDs: []string{paper.Questions[nodeIndex].ID}}
	}
	return nodeAssignment{}
}

func isBossMixed(pools []store.NodePoolBinding) bool {
	for _, pool := range pools {
		if pool.SelectionStrategy == constants.PoolSelectionBossMixed {
			return true
		}
	}
	return len(pools) > 1
}

func selectBossMixed(
	ctx context.Context,
	content *store.ContentStore,
	pools []store.NodePoolBinding,
	seen map[string]bool,
	baseSeed uint64,
) []string {
	ids := make([]string, 0, len(pools))
	for i, pool := range pools {
		poolQuestionIDs, err := content.ListPoolQuestionIDs(ctx, pool.PoolID)
		if err != nil || len(poolQuestionIDs) == 0 {
			continue
		}
		strategy := pool.SelectionStrategy
		if len(strategy) == 0 || strategy == constants.PoolSelectionBossMixed {
			strategy = constants.PoolSelectionRandomUnseen
		}
		selected, ok := selectFromPool(strategy, poolQuestionIDs, seen, baseSeed+uint64(i))
		if ok {
			ids = append(ids, selected)
		}
	}
	return ids
}

func selectFromPool(strategy string, poolQuestionIDs []string, seen map[string]bool, seed uint64) (string, bool) {
	if len(poolQuestionIDs) == 0 {
		return "", false
	}

	switch strategy {
	case constants.PoolSelectionSequential:
		return selectSequential(poolQuestionIDs, seen), true
	case constants.PoolSelectionRandomUnseen, constants.PoolSelectionBossMixed, "":
		return selectRandomUnseen(poolQuestionIDs, seen, seed), true
	default:
		return selectRandomUnseen(poolQuestionIDs, seen, seed), true
	}
}

func selectSequential(poolQuestionIDs []string, seen map[string]bool) string {
	for _, id := range poolQuestionIDs {
		if !seen[id] {
			return id
		}
	}
	return poolQuestionIDs[0]
}

func selectRandomUnseen(poolQuestionIDs []string, seen map[string]bool, seed uint64) string {
	unseen := make([]string, 0, len(poolQuestionIDs))
	for _, id := range poolQuestionIDs {
		if !seen[id] {
			unseen = append(unseen, id)
		}
	}
	candidates := unseen
	if len(candidates) == 0 {
		candidates = poolQuestionIDs
	}
	rng := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))
	return candidates[rng.IntN(len(candidates))]
}

func selectionSeed(userID, sessionID, nodeID string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(userID))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(sessionID))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(nodeID))
	return h.Sum64()
}

func allAnsweredInSession(questionIDs []string, sessionAnswered map[string]bool) bool {
	if len(questionIDs) == 0 {
		return false
	}
	for _, id := range questionIDs {
		if !sessionAnswered[id] {
			return false
		}
	}
	return true
}

func displayQuestionID(questionIDs []string, sessionAnswered map[string]bool) string {
	for _, id := range questionIDs {
		if !sessionAnswered[id] {
			return id
		}
	}
	if len(questionIDs) > 0 {
		return questionIDs[len(questionIDs)-1]
	}
	return ""
}

func nodeUnlocked(nodeIndex int, assignments []nodeAssignment, sessionAnswered map[string]bool) bool {
	if nodeIndex == 0 {
		return true
	}
	for i := 0; i < nodeIndex; i++ {
		if !allAnsweredInSession(assignments[i].questionIDs, sessionAnswered) {
			return false
		}
	}
	return true
}

func bossUnlockedIndex(node store.JourneyNode, paper *dto.DailyPaperResponse, sessionAnswered map[string]bool) bool {
	if node.NodeType != "boss" || len(paper.Questions) == 0 {
		return false
	}
	for _, q := range paper.Questions {
		if !sessionAnswered[q.ID] {
			return false
		}
	}
	return true
}
