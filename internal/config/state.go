package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type State struct {
	SyncedIDs []string `json:"synced_ids"`
}

var GlobalState State

func LoadState(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			GlobalState = State{SyncedIDs: []string{}}
			return nil
		}
		return fmt.Errorf("failed to read state file: %w", err)
	}

	if err := json.Unmarshal(data, &GlobalState); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return nil
}

func SaveState(path string) error {
	data, err := json.MarshalIndent(GlobalState, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func IsSynced(id string) bool {
	for _, syncedID := range GlobalState.SyncedIDs {
		if syncedID == id {
			return true
		}
	}
	return false
}

func AddSyncedID(id string) {
	if !IsSynced(id) {
		GlobalState.SyncedIDs = append(GlobalState.SyncedIDs, id)
	}
}
