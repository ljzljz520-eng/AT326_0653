package transport

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/service"
)

var ErrUsage = errors.New("usage: register-player, register-machine, score, board, or dashboard")

type Command struct {
	Name string
	Args []string
}

func Parse(line string) (Command, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return Command{}, ErrUsage
	}
	return Command{Name: strings.ToLower(parts[0]), Args: parts[1:]}, nil
}

func Execute(svc *service.ArcadeService, command Command) (string, error) {
	if svc == nil {
		return "", fmt.Errorf("service unavailable")
	}
	switch command.Name {
	case "register-player":
		if len(command.Args) < 2 {
			return "", ErrUsage
		}
		player, err := svc.RegisterPlayer(command.Args[0], command.Args[1], optional(command.Args, 2), optional(command.Args, 3))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("player %s registered as %s", player.ID, player.Nickname), nil
	case "register-machine":
		if len(command.Args) < 5 {
			return "", ErrUsage
		}
		machine := model.ArcadeMachine{ID: command.Args[0], Name: command.Args[1], Title: command.Args[2], Venue: command.Args[3], Difficulty: command.Args[4]}
		if _, err := svc.RegisterMachine(machine); err != nil {
			return "", err
		}
		if _, err := svc.PublishMachine(machine.ID); err != nil {
			return "", err
		}
		return fmt.Sprintf("machine %s published", machine.ID), nil
	case "score":
		if len(command.Args) < 5 {
			return "", ErrUsage
		}
		value, err := strconv.ParseInt(command.Args[2], 10, 64)
		if err != nil {
			return "", err
		}
		receipt, err := svc.SubmitScore(command.Args[0], command.Args[1], command.Args[1], value, command.Args[3], command.Args[4])
		if err != nil {
			return "", err
		}
		return receipt.Message + ": " + receipt.ScoreID, nil
	case "board":
		if len(command.Args) < 1 {
			return "", ErrUsage
		}
		entries, err := svc.Leaderboard(command.Args[0], 10)
		if err != nil {
			return "", err
		}
		return RenderBoard(entries), nil
	case "dashboard":
		machineID := optional(command.Args, 0)
		dashboard, err := svc.Dashboard(machineID, 10)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("players=%d matches=%d prizes=%d\n%s", dashboard.Players, dashboard.Matches, dashboard.ActivePrizes, RenderBoard(dashboard.TopScores)), nil
	default:
		return "", ErrUsage
	}
}

func optional(values []string, index int) string {
	if index >= 0 && index < len(values) {
		return values[index]
	}
	return ""
}
