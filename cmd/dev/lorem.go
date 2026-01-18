package dev

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"devkit/internal/output"
)

// loremCmd represents the lorem command
var loremCmd = &cobra.Command{
	Use:   "lorem [type]",
	Short: "Generate Lorem Ipsum text",
	Long: `Generate Lorem Ipsum placeholder text.

Types: word, sentence, paragraph

Examples:
  devkit dev lorem word --count 5
  devkit dev lorem sentence --count 3
  devkit dev lorem paragraph --count 2`,
	RunE: runLorem,
}

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea",
	"commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
	"in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat",
	"nulla", "pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non",
	"proident", "sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit",
	"anim", "id", "est", "laborum",
}

func init() {
	devCmd.AddCommand(loremCmd)

	loremCmd.Flags().IntP("count", "c", 1, "Number of items to generate")
	loremCmd.Flags().StringP("output", "o", "plain", "Output format: plain, json")
}

func runLorem(cmd *cobra.Command, args []string) error {
	count, _ := cmd.Flags().GetInt("count")
	outputFormat, _ := cmd.Flags().GetString("output")
	format := output.OutputFormat(outputFormat)

	if count < 1 {
		return fmt.Errorf("count must be at least 1")
	}

	var loremType string
	if len(args) > 0 {
		loremType = args[0]
	} else {
		loremType = "word"
	}

	var result string
	var results []string

	switch loremType {
	case "word":
		words := generateWords(count)
		result = strings.Join(words, " ")
		results = words
	case "sentence":
		sentences := generateSentences(count)
		result = strings.Join(sentences, " ")
		results = sentences
	case "paragraph":
		paragraphs := generateParagraphs(count)
		result = strings.Join(paragraphs, "\n\n")
		results = paragraphs
	default:
		return fmt.Errorf("invalid type: %s (supported: word, sentence, paragraph)", loremType)
	}

	if format == output.FormatJSON {
		output.PrintSuccess(format, map[string]interface{}{
			"type":   loremType,
			"count":  count,
			"text":   result,
			"items":  results,
		})
	} else {
		output.PrintSuccess(format, result)
	}

	return nil
}

func generateWords(count int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	words := make([]string, count)
	for i := 0; i < count; i++ {
		words[i] = loremWords[r.Intn(len(loremWords))]
	}
	return words
}

func generateSentences(count int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sentences := make([]string, count)
	for i := 0; i < count; i++ {
		wordCount := r.Intn(10) + 5 // 5-15 words per sentence
		words := generateWords(wordCount)
		sentence := strings.Join(words, " ")
		sentence = strings.ToUpper(string(sentence[0])) + sentence[1:] + "."
		sentences[i] = sentence
	}
	return sentences
}

func generateParagraphs(count int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	paragraphs := make([]string, count)
	for i := 0; i < count; i++ {
		sentenceCount := r.Intn(5) + 3 // 3-8 sentences per paragraph
		sentences := generateSentences(sentenceCount)
		paragraphs[i] = strings.Join(sentences, " ")
	}
	return paragraphs
}
