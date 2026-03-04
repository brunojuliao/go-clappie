package parties

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var diceRegex = regexp.MustCompile(`^(\d+)d(\d+)([+-]\d+)?$`)

// Roll parses and evaluates dice/pick expressions.
func Roll(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no dice expression provided")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	expr := strings.ToLower(args[0])

	// Coin flip
	if expr == "coin" || expr == "flip" {
		if rng.Intn(2) == 0 {
			return "Heads", nil
		}
		return "Tails", nil
	}

	// Pick from list
	if expr == "pick" {
		if len(args) < 2 {
			return "", fmt.Errorf("pick requires at least one option")
		}
		options := args[1:]
		return options[rng.Intn(len(options))], nil
	}

	// Dice roll: NdS+M
	matches := diceRegex.FindStringSubmatch(expr)
	if matches == nil {
		return "", fmt.Errorf("invalid dice expression: %s (expected NdS, NdS+M, NdS-M, coin, or pick)", expr)
	}

	numDice, _ := strconv.Atoi(matches[1])
	sides, _ := strconv.Atoi(matches[2])
	modifier := 0
	if matches[3] != "" {
		modifier, _ = strconv.Atoi(matches[3])
	}

	if numDice < 1 || numDice > 100 {
		return "", fmt.Errorf("number of dice must be 1-100")
	}
	if sides < 2 || sides > 1000 {
		return "", fmt.Errorf("sides must be 2-1000")
	}

	var rolls []int
	total := 0
	for i := 0; i < numDice; i++ {
		roll := rng.Intn(sides) + 1
		rolls = append(rolls, roll)
		total += roll
	}
	total += modifier

	// Format result
	var rollStrs []string
	for _, r := range rolls {
		rollStrs = append(rollStrs, strconv.Itoa(r))
	}

	result := fmt.Sprintf("%s: [%s]", expr, strings.Join(rollStrs, ", "))
	if modifier != 0 {
		result += fmt.Sprintf(" %+d", modifier)
	}
	result += fmt.Sprintf(" = %d", total)

	return result, nil
}
