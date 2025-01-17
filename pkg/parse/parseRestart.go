package parse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func parseRestart(p *Parser) (parseFn, error) {
	block, err := p.acceptRestart()
	if err != nil {
		return nil, err
	}
	p.cfg.Rotation = append(p.cfg.Rotation, block)
	return parseRows, nil
}

func (p *Parser) acceptRestart() (core.ActionBlock, error) {
	block := core.ActionBlock{
		Type: core.ActionBlockTypeCalcRestart,
	}
	//next token should be ;
	n := p.next()
	if n.typ != itemTerminateLine {
		return block, fmt.Errorf("reset_limit expecting ; got %v at line %v", n, p.tokens)
	}

	return block, nil
}
