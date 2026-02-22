package flow

import (
	"fmt"
	"regexp"
	"strings"
)

func (s *CodeFlowManageService) CleanFallbackBranches() error {
	branches, err := s.clients.Git.ListBranchesWithPrefix("kk-fallback")
	if err != nil {
		return err
	}

	if len(branches) == 0 {
		s.logger.Info("Clean", "message", "no fallback branches found")
		return nil
	}

	for _, branch := range branches {
		if err := s.clients.Git.DeleteLocalBranch(branch); err != nil {
			return err
		}
		s.logger.Info("Clean", "deleted", branch)
	}

	s.logger.Info("Clean", "total deleted", len(branches))
	return nil
}

func (s *CodeFlowManageService) Commit() error {
	branch, err := s.clients.Git.GetCurrentBranchName()
	if err != nil {
		return err
	}

	if err := s.clients.Git.AddAll(); err != nil {
		return err
	}
	s.logger.Info("Commit", "status", "staged all changes")

	message := commitMessageFromBranch(branch)
	if err := s.clients.Git.Commit(message); err != nil {
		return err
	}
	s.logger.Info("Commit", "committed with message", message)

	return nil
}

func (s *CodeFlowManageService) Push() error {
	out, err := s.clients.Git.StatusWithPorcelain()
	if err != nil {
		return err
	}
	if out != "" {
		if err := s.Commit(); err != nil {
			return err
		}
	} else {
		s.logger.Info("Push", "no changes to commit", "skip commit")
	}

	branch, err := s.clients.Git.GetCurrentBranchName()
	if err != nil {
		return err
	}

	if err := s.clients.Git.Push(branch); err != nil {
		return err
	}
	s.logger.Info("Push", "pushed branch", branch)

	return nil
}

func (s *CodeFlowManageService) Squash(comparableBranch string, commitMessage string) error {
	status, err := s.clients.Git.StatusWithPorcelain()
	if err != nil {
		return fmt.Errorf("failed to get working tree status: %s", err)
	}
	if strings.TrimSpace(status) != "" {
		return fmt.Errorf("working tree is not clean: %s", status)
	}

	currentBranch, err := s.clients.Git.GetCurrentBranchName()
	if err != nil {
		return err
	}
	s.logger.Info("Squash", "current branch", currentBranch)

	if comparableBranch == currentBranch {
		return fmt.Errorf("comparable branch is the same as current branch")
	}

	diff, err := s.clients.Git.GetCommitsDiffCount(comparableBranch)
	if err != nil {
		return err
	}
	s.logger.Info("Squash", "diff count", diff, "between", currentBranch, "and", comparableBranch)

	if diff <= 1 {
		s.logger.Info("Squash", "diff <= 1", "nothing to squash")
		return nil
	}

	ts := s.clients.Git.GenerateTimestamp()
	fallbackBranch := fmt.Sprintf("%s-%s-%s", "kk-fallback", currentBranch, ts)
	err = s.clients.Git.NewBranch(fallbackBranch)
	if err != nil {
		return err
	}
	s.logger.Info("Squash", "fallback branch", fallbackBranch)

	err = s.clients.Git.ResetSoft(diff)
	if err != nil {
		return err
	}
	s.logger.Info("Squash", "commits reset", diff, "on branch", currentBranch)

	err = s.clients.Git.AddAll()
	if err != nil {
		return nil
	}
	s.logger.Info("Squash", "add all changes on branch", currentBranch)

	message := commitMessage
	if message == "" {
		message = commitMessageFromBranch(currentBranch)
	}
	err = s.clients.Git.Commit(message)
	if err != nil {
		return err
	}
	s.logger.Info("Squash", "squash committed as", message, "on branch", currentBranch)

	return nil
}

var branchRegex = regexp.MustCompile(`^([A-Za-z]+)[/-](\d+)-(.+)$`)

func commitMessageFromBranch(branch string) string {
	matches := branchRegex.FindStringSubmatch(branch)
	if matches == nil {
		return branch
	}
	prefix := matches[1]
	number := matches[2]
	description := strings.ReplaceAll(matches[3], "-", " ")
	return fmt.Sprintf("[%s-%s] %s", prefix, number, description)
}
