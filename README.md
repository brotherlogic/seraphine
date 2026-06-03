# Seraphine

Seraphine is a workflow and system for running code - it is planning and review centric with an eye on putting humans involvement at the start and end of code changes - i.e. in the feature development and review stages.

## Workflow

The workflow consists of three steps - we start from a broadly defined github issue which outlines a new features. We then refine the feature and tease out the requirements, plan out some work, breaking it down into managable features and then review the code that implements the feature.

## Commands

### `seraphine init`

Initializes a repository with the "blessed" Seraphine configuration.

### `seraphine upgrade`

Automatically reconciles the local repository files with the current configuration mandated by the Seraphine server.

- Ensures alignment of GitHub workflows and scripts.
- Overwrites any local modifications to managed files to strictly enforce organizational standards.
- Safely cleans up previously managed files that have been removed from the server's blessed list.
- Commits changes to an upgrade branch (`seraphine/upgrade-<version>`) and pushes them to GitHub.

## GitHub Repository Configuration

During the upgrade process (`seraphine upgrade`), the CLI automatically reconciles GitHub repository and workflow settings (such as workflow default permissions and merge behaviors) to align with desired states provided by the server. This requires the `gh` CLI to be installed and authenticated.

### GitHub Workflows

We use automated workflows configured in the `.github/workflows/` directory to manage our processes:
- **Issue Closer**: Automatically runs periodically (or manually) to close resolved issues based on the workflow lifecycle.

## Feature planning

To build out a new feature we either create an issue with a broadly defined statement or label an existing feature request under 'seraphine-feature'.

We can then ask seraphine which feature we should be working on - it'll give you the project and issue number:

seraphine feature

If available, it'll then call the devcontainer router for that project, and initiate the interview process to take the broadly defined feature and work it into something more usable. Ultimately after this process, the seraphine-feature tag will be removed and the we'll add the 'seraphine-implement' tag

## Implementation Planning

We ask seraphine to build an implementation plan for the feature:

$ seraphine implement

It will then take the requirements from the issue, and break it down into implementation steps, potentially creating sub-issues for each piece of work. For the detailed GitHub issue workflow transitions, including decoupled planning and breakdown phases, please refer to [ISSUES.md](file:///workspaces/seraphine/ISSUES.md).

## Coding

Each of the code tagged issues are now ready for actual coding. The breadcrumb trail of proposal plan and issue specifics can guide the agent when coding. We use standard TDD approaches to address each piece.We also request the agent to focus on the revewability of the code rather than trying to  create a large single change. A number of agentic reviews are undertaken when coding the piece, but eventually pull requests are created, and assigned to the project owner for review.

## Review

You can then run seraphine-review and work through any available reviews. You can either comment, and have the agent address those comments, or approve and have the code be merged.

Note that each of these steps can be run through either (a) your subscription, (b) Model API Keys or (c) local agents - this is configured in your seraphine config for the project.