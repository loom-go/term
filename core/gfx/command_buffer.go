package gfx

import (
	"context"
	"sync"

	"github.com/AnatoleLucet/go-opentui"
)

type CommandBuffer struct {
	mu  sync.Mutex
	ctx context.Context

	commands []*Command
}

func NewCommandBuffer(ctx context.Context) *CommandBuffer {
	return &CommandBuffer{ctx: ctx}
}

func (cb *CommandBuffer) Execute(rdr *opentui.Renderer, buffer *opentui.Buffer) error {
	cb.mu.Lock()
	commands := cb.commands
	cb.commands = nil
	cb.mu.Unlock()

	for _, cmd := range commands {
		select {
		case <-cb.ctx.Done():
			return cb.ctx.Err()
		default:
		}

		switch cmd.Type {
		case CmdRender:
			err := cmd.Element.Render(buffer, cmd.Rect)
			if err != nil {
				return err
			}

		case CmdPushOverflowScissors:
			buffer.PushScissorRect(int32(cmd.Scissors.X), int32(cmd.Scissors.Y), uint32(cmd.Scissors.Width), uint32(cmd.Scissors.Height))
		case CmdPopOverflowScissors:
			buffer.PopScissorRect()

		case CmdPushHitGridScissors:
			rdr.HitGridPushScissorRect(int32(cmd.Scissors.X), int32(cmd.Scissors.Y), uint32(cmd.Scissors.Width), uint32(cmd.Scissors.Height))
		case CmdPopHitGridScissors:
			rdr.HitGridPopScissorRect()

		case CmdPushOpacity:
			buffer.PushOpacity(cmd.Opacity)
		case CmdPopOpacity:
			buffer.PopOpacity()
		}

		if cmd.Callback != nil {
			err := cmd.Callback()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (cb *CommandBuffer) Add(cmd *Command) {
	cb.mu.Lock()
	cb.commands = append(cb.commands, cmd)
	cb.mu.Unlock()
}
