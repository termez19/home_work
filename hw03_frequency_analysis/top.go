package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var wordRegex = regexp.MustCompile(`[\p{L}\d]+(?:-[\p{L}\d]+)*`)

type WordFrequency struct {
	Word  string
	Count int
}

func Top10(input string) []string {
	wordCount := make(map[string]int)
	topWords := make([]string, 0, 10)
	words := wordRegex.FindAllString(input, -1)

	for _, word := range words {
		word = strings.ToLower(word)
		wordCount[word]++
	}

	wordFrequencies := make([]WordFrequency, 0, len(wordCount))

	for i, v := range wordCount {
		wordFreq := WordFrequency{
			Word:  i,
			Count: v,
		}
		wordFrequencies = append(wordFrequencies, wordFreq)
	}

	wordFrequencies = sortStructs(wordFrequencies)
	limit := min(len(wordFrequencies), 10)

	for i := range limit {
		topWords = append(topWords, wordFrequencies[i].Word)
	}

	return topWords
}

func sortStructs(words []WordFrequency) (wordsSorted []WordFrequency) {
	sort.Slice(words, func(i, j int) bool {
		if words[i].Count == words[j].Count {
			return words[i].Word < words[j].Word
		}

		return words[i].Count > words[j].Count
	})
	wordsSorted = words
	return wordsSorted
}
