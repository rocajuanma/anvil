/*
Copyright © 2022 Juanma Roca juanmaxroca@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package installer

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// MockOutputHandler implements interfaces.OutputHandler for testing
type MockOutputHandler struct {
	messages []string
}

func (m *MockOutputHandler) PrintHeader(message string) {
	m.messages = append(m.messages, fmt.Sprintf("HEADER: %s", message))
}

func (m *MockOutputHandler) PrintStage(message string) {
	m.messages = append(m.messages, fmt.Sprintf("STAGE: %s", message))
}

func (m *MockOutputHandler) PrintSuccess(message string) {
	m.messages = append(m.messages, fmt.Sprintf("SUCCESS: %s", message))
}

func (m *MockOutputHandler) PrintError(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf("ERROR: %s", fmt.Sprintf(format, args...)))
}

func (m *MockOutputHandler) PrintWarning(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf("WARNING: %s", fmt.Sprintf(format, args...)))
}

func (m *MockOutputHandler) PrintInfo(format string, args ...interface{}) {
	m.messages = append(m.messages, fmt.Sprintf("INFO: %s", fmt.Sprintf(format, args...)))
}

func (m *MockOutputHandler) PrintProgress(current, total int, message string) {
	m.messages = append(m.messages, fmt.Sprintf("PROGRESS: %d/%d %s", current, total, message))
}

func (m *MockOutputHandler) Confirm(message string) bool {
	return true
}

func (m *MockOutputHandler) IsSupported() bool {
	return true
}

// PrintWithLevel is removed from interface

func (m *MockOutputHandler) GetMessages() []string {
	return m.messages
}

func (m *MockOutputHandler) Clear() {
	m.messages = nil
}

func TestNewConcurrentInstaller(t *testing.T) {
	tests := []struct {
		name       string
		maxWorkers int
		expected   int
	}{
		{
			name:       "Default workers (NumCPU)",
			maxWorkers: 0,
			expected:   runtime.NumCPU(),
		},
		{
			name:       "Negative workers (NumCPU)",
			maxWorkers: -1,
			expected:   runtime.NumCPU(),
		},
		{
			name:       "Custom workers",
			maxWorkers: 4,
			expected:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOutput := &MockOutputHandler{}
			installer := NewConcurrentInstaller(tt.maxWorkers, mockOutput, false)

			if installer.maxWorkers != tt.expected {
				t.Errorf("Expected maxWorkers to be %d, got %d", tt.expected, installer.maxWorkers)
			}

			if installer.timeout != time.Minute*10 {
				t.Errorf("Expected timeout to be 10 minutes, got %v", installer.timeout)
			}

			if installer.retryAttempts != 2 {
				t.Errorf("Expected retryAttempts to be 2, got %d", installer.retryAttempts)
			}
		})
	}
}

func TestConcurrentInstaller_InstallToolsEmptyList(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, false)

	ctx := context.Background()
	stats, err := installer.InstallTools(ctx, []string{})

	if err == nil {
		t.Error("Expected error for empty tools list")
	}

	if stats != nil {
		t.Error("Expected stats to be nil for empty tools list")
	}
}

func TestConcurrentInstaller_InstallToolsDryRun(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, true)

	ctx := context.Background()
	tools := []string{"tool1", "tool2", "tool3"}
	stats, err := installer.InstallTools(ctx, tools)

	if err != nil {
		t.Errorf("Expected no error for dry run, got %v", err)
	}

	if stats == nil {
		t.Error("Expected stats to be non-nil")
	}

	if stats.TotalTools != 3 {
		t.Errorf("Expected 3 total tools, got %d", stats.TotalTools)
	}

	if stats.SuccessfulTools != 3 {
		t.Errorf("Expected 3 successful tools in dry run, got %d", stats.SuccessfulTools)
	}

	if stats.FailedTools != 0 {
		t.Errorf("Expected 0 failed tools in dry run, got %d", stats.FailedTools)
	}
}

func TestConcurrentInstaller_ContextCancellation(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(1, mockOutput, false)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	tools := []string{"tool1", "tool2"}
	stats, err := installer.InstallTools(ctx, tools)

	if err == nil {
		t.Error("Expected error for cancelled context")
	}

	if stats == nil {
		t.Error("Expected stats to be non-nil even with cancelled context")
	}

	// Should have some failed tools due to cancellation
	if stats.FailedTools == 0 {
		t.Error("Expected some failed tools due to cancellation")
	}
}

func TestConcurrentInstaller_Timeout(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(1, mockOutput, false)
	installer.SetTimeout(time.Millisecond * 10) // Very short timeout

	ctx := context.Background()
	tools := []string{"nonexistent-tool"}
	stats, err := installer.InstallTools(ctx, tools)

	// Should have errors due to timeout or installation failure
	if err == nil {
		t.Error("Expected error due to timeout or installation failure")
	}

	if stats == nil {
		t.Error("Expected stats to be non-nil")
	}

	if stats.FailedTools == 0 {
		t.Error("Expected failed tools due to timeout or installation failure")
	}
}

func TestConcurrentInstaller_SetTimeout(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, false)

	newTimeout := time.Minute * 5
	installer.SetTimeout(newTimeout)

	if installer.timeout != newTimeout {
		t.Errorf("Expected timeout to be %v, got %v", newTimeout, installer.timeout)
	}
}

func TestConcurrentInstaller_SetRetryAttempts(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, false)

	newRetryAttempts := 5
	installer.SetRetryAttempts(newRetryAttempts)

	if installer.retryAttempts != newRetryAttempts {
		t.Errorf("Expected retryAttempts to be %d, got %d", newRetryAttempts, installer.retryAttempts)
	}
}

func TestInstallationResult(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(time.Second * 5)

	result := InstallationResult{
		ToolName:  "test-tool",
		Success:   true,
		Error:     nil,
		Duration:  endTime.Sub(startTime),
		StartTime: startTime,
		EndTime:   endTime,
	}

	if result.ToolName != "test-tool" {
		t.Errorf("Expected tool name to be 'test-tool', got %s", result.ToolName)
	}

	if !result.Success {
		t.Error("Expected success to be true")
	}

	if result.Error != nil {
		t.Errorf("Expected error to be nil, got %v", result.Error)
	}

	expectedDuration := time.Second * 5
	if result.Duration != expectedDuration {
		t.Errorf("Expected duration to be %v, got %v", expectedDuration, result.Duration)
	}
}

func TestInstallationStats(t *testing.T) {
	stats := InstallationStats{
		TotalTools:      5,
		SuccessfulTools: 3,
		FailedTools:     2,
		TotalDuration:   time.Minute * 2,
		AverageDuration: time.Second * 30,
		MaxDuration:     time.Second * 45,
		MinDuration:     time.Second * 15,
		ConcurrentJobs:  4,
	}

	if stats.TotalTools != 5 {
		t.Errorf("Expected total tools to be 5, got %d", stats.TotalTools)
	}

	if stats.SuccessfulTools != 3 {
		t.Errorf("Expected successful tools to be 3, got %d", stats.SuccessfulTools)
	}

	if stats.FailedTools != 2 {
		t.Errorf("Expected failed tools to be 2, got %d", stats.FailedTools)
	}

	if stats.ConcurrentJobs != 4 {
		t.Errorf("Expected concurrent jobs to be 4, got %d", stats.ConcurrentJobs)
	}
}

func TestConcurrentInstaller_calculateStats(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, false)

	startTime := time.Now()
	results := []InstallationResult{
		{
			ToolName:  "tool1",
			Success:   true,
			Duration:  time.Second * 10,
			StartTime: startTime,
			EndTime:   startTime.Add(time.Second * 10),
		},
		{
			ToolName:  "tool2",
			Success:   false,
			Error:     fmt.Errorf("installation failed"),
			Duration:  time.Second * 5,
			StartTime: startTime,
			EndTime:   startTime.Add(time.Second * 5),
		},
		{
			ToolName:  "tool3",
			Success:   true,
			Duration:  time.Second * 15,
			StartTime: startTime,
			EndTime:   startTime.Add(time.Second * 15),
		},
	}

	stats := installer.calculateStats(results, startTime)

	if stats.TotalTools != 3 {
		t.Errorf("Expected total tools to be 3, got %d", stats.TotalTools)
	}

	if stats.SuccessfulTools != 2 {
		t.Errorf("Expected successful tools to be 2, got %d", stats.SuccessfulTools)
	}

	if stats.FailedTools != 1 {
		t.Errorf("Expected failed tools to be 1, got %d", stats.FailedTools)
	}

	expectedAverage := time.Second * 10 // (10 + 5 + 15) / 3
	if stats.AverageDuration != expectedAverage {
		t.Errorf("Expected average duration to be %v, got %v", expectedAverage, stats.AverageDuration)
	}

	if stats.MinDuration != time.Second*5 {
		t.Errorf("Expected min duration to be 5s, got %v", stats.MinDuration)
	}

	if stats.MaxDuration != time.Second*15 {
		t.Errorf("Expected max duration to be 15s, got %v", stats.MaxDuration)
	}
}

func TestConcurrentInstaller_printProgress(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, false)

	result := InstallationResult{
		ToolName: "test-tool",
		Success:  true,
		Duration: time.Second * 2,
	}

	installer.printProgress(result, 1, 3)

	messages := mockOutput.GetMessages()
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if messages[0] != "PROGRESS: 1/3 ✓ test-tool (2s)" {
		t.Errorf("Expected progress message, got: %s", messages[0])
	}
}

func TestConcurrentInstaller_printProgressFailed(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, false)

	result := InstallationResult{
		ToolName: "test-tool",
		Success:  false,
		Duration: time.Second * 2,
	}

	installer.printProgress(result, 1, 3)

	messages := mockOutput.GetMessages()
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}

	if messages[0] != "PROGRESS: 1/3 ✗ test-tool (2s)" {
		t.Errorf("Expected progress message with failure, got: %s", messages[0])
	}
}

func TestConcurrentInstaller_printSummary(t *testing.T) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(2, mockOutput, false)

	stats := &InstallationStats{
		TotalTools:      3,
		SuccessfulTools: 2,
		FailedTools:     1,
		TotalDuration:   time.Second * 30,
		AverageDuration: time.Second * 10,
		ConcurrentJobs:  2,
	}

	results := []InstallationResult{
		{ToolName: "tool1", Success: true},
		{ToolName: "tool2", Success: false, Error: fmt.Errorf("failed")},
		{ToolName: "tool3", Success: true},
	}

	installer.printSummary(stats, results)

	messages := mockOutput.GetMessages()

	// Should have header, success info, timing info, workers info, failed tools warning, and speedup estimate
	if len(messages) < 5 {
		t.Errorf("Expected at least 5 messages, got %d", len(messages))
	}

	// Check for header
	found := false
	for _, msg := range messages {
		if msg == "HEADER: Concurrent Installation Complete" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected header message not found")
	}

	// Check for failed tools warning
	found = false
	for _, msg := range messages {
		if msg == "WARNING: Failed installations:" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected failed tools warning not found")
	}
}

// Benchmark tests
func BenchmarkConcurrentInstaller_DryRun(b *testing.B) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(4, mockOutput, true)
	ctx := context.Background()
	tools := []string{"tool1", "tool2", "tool3", "tool4", "tool5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockOutput.Clear()
		installer.InstallTools(ctx, tools)
	}
}

func BenchmarkConcurrentInstaller_calculateStats(b *testing.B) {
	mockOutput := &MockOutputHandler{}
	installer := NewConcurrentInstaller(4, mockOutput, false)

	startTime := time.Now()
	results := make([]InstallationResult, 100)
	for i := 0; i < 100; i++ {
		results[i] = InstallationResult{
			ToolName:  fmt.Sprintf("tool%d", i),
			Success:   i%2 == 0,
			Duration:  time.Duration(i) * time.Millisecond,
			StartTime: startTime,
			EndTime:   startTime.Add(time.Duration(i) * time.Millisecond),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		installer.calculateStats(results, startTime)
	}
}
