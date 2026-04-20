<!-- TRIAGED -->
# Workflow

I start by onboarding a github project onto Seraphine:

seraphine onboard <github_project>

The seraphine server acts on this reqeust by:

(a) trying to check out the project
    - If this doens't succeed we report an onboarding failure highlighting why
(b) creating a seraphine/ directory in the project and pushing the change
    - If this doesn't succeed we report on why we couldn't push the change.

The server then stores the project as ONBOARDING in it's central repo, alongside
the seraphine version number.

We then enqueue an action to ONBOARD the server. This is a non-agentic taskt that
(a) checks out the repo and (b) adds in the seraphine manager and configures the
system to use the seraphine webhook.

Once these are complete we push the change to the repo.

## Open Questions

Please provide your answers below:

1. **State Management ("Central Repo"):**
   The proposal mentions: *"The server then stores the project as ONBOARDING in its central repo"*. Is this "central repo" a backend database (like SQLite, Postgres), a local configuration file, or an actual central GitHub repository dedicated to tracking the server's global state?
   - **Answer:** 

2. **Authentication & Permissions:**
   How is the server authenticating to GitHub to check out, create a `seraphine/` directory, and push changes? Will it use a GitHub App installation token, a Personal Access Token (PAT), or SSH keys? 
   - **Answer:** 

3. **Defining the "Seraphine Manager":**
   The non-agentic task *"adds in the seraphine manager"*. What exactly is the Seraphine manager? Is it a specific GitHub Action workflow file `.github/workflows/...`, a GitHub App installation, or something else?
   - **Answer:** 

4. **Webhook Configuration:**
   It states it *"configures the system to use the seraphine webhook"*. To do this via the GitHub API requires admin privileges on the repo. Is the Seraphine server using a token with these privileges? Is the webhook endpoint URL stored in the server's environment?
   - **Answer:** 

5. **Client-Server Communication:**
   How does the CLI (`seraphine onboard <github_project>`) communicate this request to the Seraphine server? REST API, gRPC, or something else?
   - **Answer:** 

6. **Signifier File Consistency:**
   `README.md` mentions: *"Any repo... that has a `seraphine.md` file in the root is considered in scope"*. But `onboard.md` creates a `seraphine/` directory instead. Should the onboarding process also create the root `seraphine.md` signifier file, or are we just relying on the `seraphine/` directory now?
   - **Answer:** 