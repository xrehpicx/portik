package ports

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ParseSpec parses "5432,6379,3000-3010" into a sorted, de-duplicated list of ports.
func ParseSpec(spec string) ([]int, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, errors.New("empty ports spec")
	}

	var out []int
	seen := map[int]bool{}

	tokens := strings.Split(spec, ",")
	for _, t := range tokens {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}

		// range: "3000-3010"
		if strings.Contains(t, "-") {
			parts := strings.SplitN(t, "-", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid range token: %q", t)
			}
			a, err := parsePort(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range start %q: %w", parts[0], err)
			}
			b, err := parsePort(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end %q: %w", parts[1], err)
			}
			if a > b {
				a, b = b, a
			}
			for p := a; p <= b; p++ {
				if !seen[p] {
					seen[p] = true
					out = append(out, p)
				}
			}
			continue
		}

		// single
		p, err := parsePort(t)
		if err != nil {
			return nil, err
		}
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}

	if len(out) == 0 {
		return nil, errors.New("no ports parsed from spec")
	}
	sort.Ints(out)
	return out, nil
}

func parsePort(s string) (int, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if n <= 0 || n > 65535 {
		return 0, errors.New("port out of range (1-65535)")
	}
	return n, nil
}
