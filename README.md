# Seraphine

Seraphine is a pure agentic coder - effectively you tell it to track a github project and it will monitor for feature proposals, write an implementation plan to meet these proposals, and then farm out the implementation work to other agents in order to make the feature live. It then documents the feature, updates an overall design doc and marks the feature as complete.

As such it is a coding agent for which all interaction happens through either (a) writing a design doc for a feature, and (b) reviewing pull requests. This releases you from having to write code and focuses on a holistic view of the system.

## Features

Everything in Seraphine happens through a feature proposal - the user adds a file to a tracked repo under the /features directory. This is a markdown file that describes the feature at a relatively high level. It should describe how the feature works and what libraries it should use and how it should fit into the existing project structure. It should have a section at the top describing the state of the feature (PROPOSED, PLANNED, IMPLEMENTING, COMPLETE).

A PROPOSED feature is picked up by the architecht agent - this looks at the proposal and defines an implementation plan and individual tasks that tie into this proposal. It updates the feature file with the implementation plan and marks the feature as PLANNED. It then farms out the tasks to the implementation agents.

An IMPLEMENTING feature is picked up by the implementation agents - these agents look at the tasks and implement them. They update the feature file with the results of their work and mark the feature as COMPLETE.

The tech writer agent then (a) writes a document outlining how the feature was implemented and (b) updates the overall system readme to adjust how the system appears as a whole.

## Changes

Individual code changes are pushed by one of the agents, and is signed by that agent (i.e. the reader knows which agent produced a given code / documentation change). The repo is configured in such a way that any change made to the system is reviewed by both a code reviewing agent and you. Once approved and tests pass the change is merged.

## Tracked Repos

Any repo under the given github user that has a seraphine.md file in the root is considered in scope for tracking. The full set of user repos is checked every hour and any recently changed repos are checked for the seraphine signifier file.

## Development

This project includes a [devcontainer](file:///var/home/simon/seraphine/.devcontainer/devcontainer.json) configuration. To use it:
1. Ensure you have VS Code and the "Dev Containers" extension installed.
2. Open the project in VS Code.
3. Click "Reopen in Container" when prompted.
The environment is based on Ubuntu and includes Go, Zsh, and Oh My Zsh.

Once a repo is in scope we configure the github repo to match the settings required by the changes section above:

  1. It should have a branch protection rule against merges to main
  1. It should not allow direct pushes to the main branch
  1. It should only pass a pull request that has a gemini review
  1. It should only pass a pull request where tests all pass
  1. It should have a commit that addresses the findings of the gemini review
