package api

import (
	"math/big"
	"sort"
)

const (
	sortKeyViews          = "views"
	sortKeyRelevanceScore = "relevanceScore"
)

func sortLinksBySortKey(links []Link, sortKey string) {
	if sortKey == sortKeyRelevanceScore {
		sortByRelevanceScore(links)
		return
	}
	sortByViews(links)
}

func sortByRelevanceScore(links []Link) {
	sort.SliceStable(links, func(i, j int) bool {
		lBigFloat := big.NewFloat(float64(links[i].RelevanceScore))
		rBigFloat := big.NewFloat(float64(links[j].RelevanceScore))
		// lBigFloat.Cmp(rBigFloat) returns -1 if lBigFloat < rBigFloat,
		// 0 if lBigFloat == rBigFloat,
		// 1 if lBigFloat > rBigFloat
		cmp := lBigFloat.Cmp(rBigFloat)
		// float equivalent for links[i].RelevanceScore > links[j].RelevanceScore
		return cmp == 1
	})
}

func sortByViews(links []Link) {
	sort.SliceStable(links, func(i, j int) bool {
		return links[i].Views > links[j].Views
	})
}
