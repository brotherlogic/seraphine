# 🕵️ The `seraphine-needs-deep-research` Label Workflow

When a parent issue is labeled with `seraphine-needs-deep-research`, the AI assistant is triggered to execute a deep research phase to evaluate the problem space and propose solutions.

## 🔄 Workflow Lifecycle

1. **Search**: The agent leverages the `deep-research` skill to explore the problem space and evaluate potential solutions. The research continues until at least the top 3 viable options are identified (this is the sufficiency threshold).
2. **Synthesize**: Once the options are identified, the agent synthesizes the findings, organically presenting the pros and cons of these options.
3. **Comment**: The agent posts the synthesized summary as a detailed comment on the parent issue.
4. **Remove Label**: The agent removes the `seraphine-needs-deep-research` label from the parent issue using the `gh` CLI.
5. **Create Sub-issue**: The agent successfully creates a child `[Requirements]` sub-issue via the GitHub API and labels it with `seraphine-needs-requirements` to initiate the requirements gathering phase.

---

## 📋 Phase Guidelines

### 1. Search
* Perform deep semantic searches across the entire repository.
* Execute targeted web queries for relevant open-source libraries, documentation, or best practices.
* Review any architectural decision records (ADRs) and product requirements documents (PRDs).
* Leverage the `deep-research` skill and stop searching once at least 3 viable options are found.

### 2. Synthesize
* Organize findings into distinct, viable approaches (at least 3 options).
* For each approach, document the pros, cons, implementation complexity, and risks organically.

### 3. Comment
* Format the synthesized research as a detailed GitHub comment.
* Ensure markdown elements like tables, code blocks, and bold text are used for readability.

### 4. Remove Label
* Once the research comment is posted successfully, remove the `seraphine-needs-deep-research` label using the `gh` CLI.

### 5. Create Sub-issue
* Create a child sub-issue using the `gh` CLI or GitHub API. Title it `[Requirements] <Title>`.
* Ensure the sub-issue describes that deep research is complete and it is time for requirements gathering.
* Label the new sub-issue with `seraphine-needs-requirements`.

---

## ⚠️ Error States & Handling

The agent must explicitly handle the following error states:

### No Viable Options Found
If the agent is unable to identify at least 3 viable options after an exhaustive search:
1. Document the findings and constraints preventing a solution.
2. Explain why viable options were not found.
3. Post this summary as a comment to the parent issue.
4. Tag the author/product owner for clarification.
5. Remove the `seraphine-needs-deep-research` label.
6. Do **not** create the requirements sub-issue; leave the parent issue in an open, unlabeled state.

### Ambiguous Problem Space
If the problem space is too ambiguous to conduct meaningful research:
1. Comment detailing what specific information is missing or unclear and ask for clarification.
2. Remove the `seraphine-needs-deep-research` label.
3. Halt execution and do **not** create the requirements sub-issue.
