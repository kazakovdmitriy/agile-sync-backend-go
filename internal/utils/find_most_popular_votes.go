package utils

import (
	"backend_go/internal/model"
	"backend_go/internal/model/entitymodel"
)

// FindMostPopularVote возвращает самую популярную оценку.
// При равенстве голосов выбирает значение с наивысшим рангом (последнее в DeckValues).
func FindMostPopularVote(votes []entitymodel.Vote, deckType model.DeckType) *string {
	// Собираем и фильтруем голоса
	counts := make(map[string]int)
	for _, v := range votes {
		if v.Value != "hidden" {
			counts[v.Value]++
		}
	}

	if len(counts) == 0 {
		return nil
	}

	// Находим максимальное количество голосов
	maxCount := 0
	for _, count := range counts {
		if count > maxCount {
			maxCount = count
		}
	}

	// Собираем кандидатов с максимальным количеством
	var candidates []string
	for value, count := range counts {
		if count == maxCount {
			candidates = append(candidates, value)
		}
	}

	if len(candidates) == 1 {
		return &candidates[0]
	}

	// Создаём мапу рангов из порядка в DeckValues
	deckValues, exists := model.DeckValues[deckType]
	if !exists {
		// Если тип колоды неизвестен
		return &candidates[0]
	}

	rankMap := make(map[string]int, len(deckValues))
	for idx, value := range deckValues {
		rankMap[value] = idx
	}

	// Выбираем кандидата с наивысшим рангом
	bestCandidate := candidates[0]
	bestRank := rankMap[bestCandidate]

	for _, candidate := range candidates[1:] {
		// Пропускаем значения, которых нет в колоде
		currentRank, exists := rankMap[candidate]
		if !exists {
			continue
		}

		if currentRank > bestRank {
			bestRank = currentRank
			bestCandidate = candidate
		}
	}

	return &bestCandidate
}
