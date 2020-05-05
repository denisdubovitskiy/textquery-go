package query

import (
	"errors"
	"sort"
	"strings"
)

type PartType int

const (
	PartTypeParenthesis PartType = iota
	PartTypeOperator
	PartTypeData
)

type Parser struct {
	openParens, closeParens string
	operatorsOrCloseParens  map[string]bool
	operatorStartsWith      map[string]string
	parentheses             map[string]bool
	fieldDelimiter          string
	fieldCloseParens        string
	fieldOpenParens         string
	operators               map[string]bool
	allOpenParens           map[string]bool
	closeToOpenParens       map[string]string
	allCloseParens          map[string]bool
	modifierCloseParens     string
	modifierOpenParens      string
}

type Options struct {
	openParens          string
	closeParens         string
	fieldOpenParens     string
	fieldCloseParens    string
	fieldDelimiter      string
	operators           []string
	modifierOpenParens  string
	modifierCloseParens string
}

type Option func(o *Options)

func New(opts ...Option) *Parser {
	options := Options{
		openParens:  "(",
		closeParens: ")",

		modifierOpenParens:  "{",
		modifierCloseParens: "}",

		fieldDelimiter:   ":",
		fieldOpenParens:  "[",
		fieldCloseParens: "]",
		operators:        []string{"AND", "OR", "NOT"},
	}

	for _, opt := range opts {
		opt(&options)
	}

	operatorsOrCloseParens := convertToSet(options.operators)
	operatorsOrCloseParens[options.closeParens] = true

	return &Parser{
		openParens:  options.openParens,
		closeParens: options.closeParens,

		operators: convertToSet(options.operators),

		operatorStartsWith: convertOperators(options.operators),

		parentheses: map[string]bool{
			options.openParens:  true,
			options.closeParens: true,
		},

		fieldDelimiter:   options.fieldDelimiter,
		fieldOpenParens:  options.fieldOpenParens,
		fieldCloseParens: options.fieldCloseParens,

		modifierOpenParens:  options.modifierOpenParens,
		modifierCloseParens: options.modifierCloseParens,

		operatorsOrCloseParens: operatorsOrCloseParens,

		allOpenParens: map[string]bool{
			options.openParens:         true,
			options.fieldOpenParens:    true,
			options.modifierOpenParens: true,
		},

		allCloseParens: map[string]bool{
			options.closeParens:         true,
			options.fieldCloseParens:    true,
			options.modifierCloseParens: true,
		},

		closeToOpenParens: map[string]string{
			options.closeParens:         options.openParens,
			options.fieldCloseParens:    options.fieldOpenParens,
			options.modifierCloseParens: options.modifierOpenParens,
		},
	}
}

func convertOperators(operators []string) map[string]string {
	mp := make(map[string]string, len(operators))

	for _, op := range operators {
		mp[string(op[0])] = op
	}

	return mp
}

func convertToSet(operators []string) map[string]bool {
	mp := make(map[string]bool, len(operators))

	for _, op := range operators {
		mp[op] = true
	}

	return mp
}

func (p *Parser) checkNeedsWrapping(query string) bool {
	chunks := make([][]int, 0, strings.Count(query, p.openParens))
	currentChunk := []int{0, 0}

	for i, char := range query {
		letter := string(char)

		if letter == p.openParens {
			currentChunk[0] = i
		}

		if letter == p.closeParens {
			currentChunk[1] = i
			chunks = append(chunks, currentChunk)
			currentChunk = []int{0, 0}
		}
	}

	if len(chunks) == 0 {
		return true
	}

	sort.Slice(chunks, func(i, j int) bool {
		leftChunk := chunks[i]
		rightChunk := chunks[j]

		leftChunkLen := leftChunk[1] - leftChunk[0]
		rightChunkLen := rightChunk[1] - rightChunk[0]

		return leftChunkLen > rightChunkLen
	})

	largestChunk := chunks[0]
	largestChunkLen := largestChunk[1] - largestChunk[0]
	return largestChunkLen != len(query)-1
}

type Part struct {
	partType PartType
	field    string
	operator string
	data     string
	modifier string
	position int
}

func (p *Parser) splitQuery(query string) []Part {
	if p.checkNeedsWrapping(query) {
		query = "(" + query + ")"
	}

	charIndex := 0

	parts := make([]Part, 0, p.countPartsCapacity(query))

	for charIndex < len(query) {
		character := string(query[charIndex])

		if operator, ok := p.operatorStartsWith[character]; ok {
			if nextChars := query[charIndex : charIndex+len(operator)]; nextChars == operator {
				parts = append(parts, Part{data: operator, partType: PartTypeOperator})
				charIndex += len(operator)
				continue
			}
		}

		if p.parentheses[character] {
			parts = append(parts, Part{data: character, position: charIndex})
			charIndex++
			continue
		}

		field, fieldOperator, modifier, word := "", "", "", ""
		partIndex := charIndex
		for partIndex < len(query) {
			character = string(query[partIndex])

			if p.parentheses[character] {
				word = strings.TrimSpace(word)
				if len(word) > 0 {
					parts = append(parts, Part{
						field:    field,
						operator: fieldOperator,
						data:     word,
						modifier: modifier,
						partType: PartTypeData,
					})
				}

				field, fieldOperator, modifier, word = "", "", "", ""
				charIndex = partIndex - 1
				break
			}

			if operator, ok := p.operatorStartsWith[character]; ok {
				if nextChars := query[partIndex : partIndex+len(operator)]; nextChars == operator {
					word = strings.TrimSpace(word)
					if len(word) > 0 {
						parts = append(parts, Part{
							field:    field,
							operator: fieldOperator,
							modifier: modifier,
							data:     word,
							partType: PartTypeData,
						})
					}

					field, fieldOperator, modifier, word = "", "", "", ""
					charIndex = partIndex - 1

					break
				}
			}

			if character == p.fieldDelimiter {
				field = word

				if strings.HasSuffix(field, p.fieldCloseParens) {
					openParensIndex := strings.Index(field, p.fieldOpenParens)
					fieldOperator = field[openParensIndex+1 : len(field)-1]
					field = field[:openParensIndex]
				}

				if strings.Contains(field, p.modifierCloseParens) {
					openParensIndex := strings.Index(field, p.modifierOpenParens)
					closeParensIndex := strings.Index(field, p.modifierCloseParens)
					modifier = field[openParensIndex+1 : closeParensIndex]
					field = field[:openParensIndex]
				}

				word = ""
				partIndex++
				continue
			}

			word += character
			partIndex++
		}
		charIndex++
	}

	return parts
}

func (p *Parser) Parse(query string) (*Node, error) {
	if err := p.isBalanced(query); err != nil {
		return nil, err
	}

	splittedQuery := p.splitQuery(query)
	stack := newNodeStack(p.countNodeStackCapacity(splittedQuery))

	tree := &Node{}
	stack.push(tree)

	currentNode := tree

	for _, part := range splittedQuery {
		if part.data == p.openParens {
			currentNode.insertLeft()
			stack.push(currentNode)
			currentNode = currentNode.accessLeft()
		} else if !p.operatorsOrCloseParens[part.data] {
			currentNode.Data = part
			parentNode, ok := stack.pop()
			if !ok {
				return nil, errors.New("pop from empty stack")
			}

			currentNode = parentNode
		} else if p.operators[part.data] {
			currentNode.Data = part
			currentNode.insertRight()
			stack.push(currentNode)
			currentNode = currentNode.accessRight()
		} else if part.data == p.closeParens {
			parentNode, ok := stack.pop()
			if !ok {
				return nil, errors.New("pop from empty stack")
			}

			currentNode = parentNode
		} else {
			return nil, errors.New("error")
		}
	}

	if err := p.checkForErrors(tree); err != nil {
		return nil, err
	}

	if len(splittedQuery) == 3 {
		return tree.accessLeft(), nil
	}

	return tree, nil
}

func (p *Parser) countNodeStackCapacity(query []Part) int {
	capacity := 0

	for _, part := range query {
		if part.partType == PartTypeData || part.partType == PartTypeOperator {
			capacity++
		}
	}

	return capacity
}

func (p *Parser) countPartsCapacity(query string) int {
	capacity := 0

	for op := range p.operators {
		capacity += strings.Count(query, op)
	}

	for p := range p.allOpenParens {
		capacity += strings.Count(query, p)
	}

	for p := range p.allCloseParens {
		capacity += strings.Count(query, p)
	}

	return capacity * 2
}

func (p *Parser) checkForErrors(node *Node) error {
	if node.Data.partType == PartTypeParenthesis {
		return ValidationError{
			message:  "invalid query syntax",
			position: node.Data.position,
		}
	}

	for _, child := range node.Children {
		if err := p.checkForErrors(child); err != nil {
			return err
		}
	}

	return nil
}
