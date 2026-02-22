package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type GitClient struct{}

func NewGitClient() *GitClient {
	return &GitClient{}
}

func (g *GitClient) GetCurrentBranchName() (string, error) {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	branch := strings.TrimSpace(string(out))

	return branch, nil
}

func (p *GitClient) GenerateTimestamp() string {
	return time.Now().Format("2006-01-02-15-04-05")
}

func (g *GitClient) NewBranch(branchName string) error {
	if out, err := exec.Command("git", "branch", branchName).Output(); err != nil {
		return fmt.Errorf("failed to create branch %s: %s", branchName, out)
	}
	return nil
}

func (g *GitClient) SwitchToBranch(branchName string) error {
	if out, err := exec.Command("git", "switch", branchName).Output(); err != nil {
		return fmt.Errorf("failed to switch branch %s: %s", branchName, out)
	}
	return nil
}

func (g *GitClient) GetCommitsDiffCount(targetBranch string) (int, error) {
	out, err := exec.Command("git", "cherry", "-v", targetBranch).Output()
	if err != nil {
		return 0, fmt.Errorf("failed to count diff for target branch %s: %w", targetBranch, err)
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return 0, nil
	}

	lines := strings.Split(trimmed, "\n")
	return len(lines), nil
}

func (g *GitClient) Commit(message string) error {
	out, err := exec.Command("git", "commit", "-m", message).Output()
	if err != nil {
		return fmt.Errorf("failed to commit %s", out)
	}
	return err
}

func (g *GitClient) Push(branchName string) error {
	arg := []string{"push", "--set-upstream", "origin", branchName}
	out, err := exec.Command("git", arg...).Output()
	if err != nil {
		return fmt.Errorf("failed to push %s", out)
	}
	return err
}

func (g *GitClient) ResetSoft(commitsFromHead int) error {
	commitsToReset := fmt.Sprintf("HEAD~%d", commitsFromHead)
	out, err := exec.Command("git", "reset", "--soft", commitsToReset).Output()
	if err != nil {
		return fmt.Errorf("failed to reset softly %s", out)
	}
	return err
}

func (g *GitClient) AddAll() error {
	out, err := exec.Command("git", "add", ".").Output()
	if err != nil {
		return fmt.Errorf("failed to add all changes %s", out)
	}
	return err
}

func (g *GitClient) DeleteLocalBranch(branchName string) error {
	out, err := exec.Command("git", "branch", "-D", branchName).Output()
	if err != nil {
		return fmt.Errorf("failed to delete branch %s err %s", branchName, out)
	}
	return err
}

func (g *GitClient) StatusWithPorcelain() (string, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return "", fmt.Errorf("failed to check working tree status: %s", out)
	}
	return string(out), nil
}

func (g *GitClient) ListBranchesWithPrefix(prefix string) ([]string, error) {
	out, err := exec.Command("git", "branch", "--list", prefix+"*").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return nil, nil
	}

	var branches []string
	for _, line := range strings.Split(trimmed, "\n") {
		branch := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		if branch != "" {
			branches = append(branches, branch)
		}
	}
	return branches, nil
}
