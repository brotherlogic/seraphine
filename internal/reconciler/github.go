package reconciler

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"

	pb "github.com/brotherlogic/seraphine/proto"
)

var runCommand = func(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Output()
}

// ReconcileGithubSettings compares the desired settings with the repository's current state and applies updates using the gh CLI.
func ReconcileGithubSettings(ctx context.Context, ownerRepo string, settings []*pb.GithubSetting) error {
	if len(settings) == 0 {
		return nil
	}

	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI is not installed or not in PATH: %w", err)
	}

	// Registry of known setting definitions mapping proto keys to GitHub API endpoints and fields
	type settingDef struct {
		getPath    string // relative to repos/{ownerRepo}
		updatePath string // relative to repos/{ownerRepo}
		method     string // "PUT" or "PATCH"
		jsonKey    string
		valueType  string // "bool", "string"
	}

	knownSettings := map[string]settingDef{
		"default_workflow_permissions": {
			getPath:    "/actions/permissions/workflow",
			updatePath: "/actions/permissions/workflow",
			method:     "PUT",
			jsonKey:    "default_workflow_permissions",
			valueType:  "string",
		},
		"can_approve_pull_request_reviews": {
			getPath:    "/actions/permissions/workflow",
			updatePath: "/actions/permissions/workflow",
			method:     "PUT",
			jsonKey:    "can_approve_pull_request_reviews",
			valueType:  "bool",
		},
		"delete_branch_on_merge": {
			getPath:    "",
			updatePath: "",
			method:     "PATCH",
			jsonKey:    "delete_branch_on_merge",
			valueType:  "bool",
		},
		"allow_squash_merge": {
			getPath:    "",
			updatePath: "",
			method:     "PATCH",
			jsonKey:    "allow_squash_merge",
			valueType:  "bool",
		},
		"has_issues": {
			getPath:    "",
			updatePath: "",
			method:     "PATCH",
			jsonKey:    "has_issues",
			valueType:  "bool",
		},
	}

	// Cache map to store parsed responses per endpoint path
	cache := make(map[string]map[string]any)

	for _, setting := range settings {
		def, ok := knownSettings[setting.Key]
		if !ok {
			return fmt.Errorf("unknown GitHub setting: %s", setting.Key)
		}

		apiPath := fmt.Sprintf("repos/%s%s", ownerRepo, def.getPath)

		// Fetch current state of this endpoint if not already cached
		if _, cached := cache[apiPath]; !cached {
			out, err := runCommand(ctx, "gh", "api", apiPath)
			if err != nil {
				return fmt.Errorf("failed to fetch current GitHub settings from %s: %w, output: %s", apiPath, err, string(out))
			}

			var data map[string]any
			if err := json.Unmarshal(out, &data); err != nil {
				return fmt.Errorf("failed to parse GitHub API response from %s: %w", apiPath, err)
			}
			cache[apiPath] = data
		}

		currentData := cache[apiPath]
		currentVal, valExists := currentData[def.jsonKey]

		needsUpdate := false
		switch def.valueType {
		case "bool":
			desiredBool, err := strconv.ParseBool(setting.Value)
			if err != nil {
				return fmt.Errorf("invalid boolean value for setting %s: %s", setting.Key, setting.Value)
			}

			var currentBool bool
			if valExists {
				if cb, ok := currentVal.(bool); ok {
					currentBool = cb
				} else {
					return fmt.Errorf("current value for %s is not a boolean: %v", setting.Key, currentVal)
				}
			}
			if !valExists || currentBool != desiredBool {
				needsUpdate = true
			}

		case "string":
			var currentStr string
			if valExists {
				if cs, ok := currentVal.(string); ok {
					currentStr = cs
				} else {
					return fmt.Errorf("current value for %s is not a string: %v", setting.Key, currentVal)
				}
			}
			if !valExists || currentStr != setting.Value {
				needsUpdate = true
			}
		}

		if needsUpdate {
			updatePath := fmt.Sprintf("repos/%s%s", ownerRepo, def.updatePath)
			fieldArg := fmt.Sprintf("%s=%s", def.jsonKey, setting.Value)

			_, err := runCommand(ctx, "gh", "api", "-X", def.method, updatePath, "-F", fieldArg)
			if err != nil {
				return fmt.Errorf("failed to update GitHub setting %s to %s: %w", setting.Key, setting.Value, err)
			}

			// Update the cache with the new value
			var parsedDesired any
			if def.valueType == "bool" {
				parsedDesired, _ = strconv.ParseBool(setting.Value)
			} else {
				parsedDesired = setting.Value
			}
			currentData[def.jsonKey] = parsedDesired
		}
	}

	return nil
}
