package console

import (
	"fmt"
	"github.com/alioth-center/infrastructure/env"
	"github.com/alioth-center/infrastructure/utils/timezone"
	"os"
	"strings"
)

func NewBlock() Block {
	return &block{
		title:    "",
		message:  []*Messages{},
		subBlock: []Block{},
	}
}

type Block interface {
	Title(title string) Block
	Message(messages ...Messages) Block
	SubBlock(blocks ...Block) Block
	print(idx int)
}

type block struct {
	title    string
	message  []*Messages
	subBlock []Block
}

func (b *block) Title(title string) Block {
	b.title = title
	return b
}

func (b *block) Message(messages ...Messages) Block {
	for _, message := range messages {
		if message.builder == nil {
			continue
		}

		b.message = append(b.message, &message)
	}
	return b
}

func (b *block) SubBlock(blocks ...Block) Block {
	for _, sub := range blocks {
		b.subBlock = append(b.subBlock, sub)
	}
	return b
}

func (b *block) print(idx int) {
	prev := strings.Repeat("\t", idx)
	_, _ = fmt.Fprintf(os.Stdout, "%s[%s]: %s\n", prev, b.title, timezone.NowInLocalTime().Format(env.GetEnv(env.AliothFrameworkTimeFormatKey)))
	for _, message := range b.message {
		_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", prev, message.builder.String())
	}
	for _, sub := range b.subBlock {
		sub.print(idx + 1)
	}
	_, _ = fmt.Fprintln(os.Stdout)
}
