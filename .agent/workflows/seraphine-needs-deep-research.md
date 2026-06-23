# 🔍 The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the agent should execute a deep dive to explore the problem space and find technical solutions before moving on to the requirements phase.

## 🔄 Workflow Sequence

### 1. Search
* Begin by executing comprehensive web and codebase searches relevant to the issue description.
* Explore documentation, similar existing issues, relevant PRs, and architectural decisions.
* Gather context on available libraries, patterns, and alternative approaches.

### 2. Synthesize
* Synthesize the collected information into a cohesive understanding of the problem.
* Identify potential solutions, trade-offs, and technical considerations.
* Formulate clear recommendations.

### 3. Comment
* Post a detailed summary comment on the issue containing your synthesized findings.
* Include links to sources, code snippets, and structured recommendations.

### 4. Remove Label
* Once the comment is successfully posted, remove the `seraphine-needs-deep-research` label from the issue.

### 5. Create Sub-issue
* If the research points to actionable next steps, label the issue with `seraphine-needs-requirements` to transition to the requirements gathering phase, or create a specific sub-issue if appropriate based on the findings.
* (As per general workflow guidelines, the parent issue typically transitions to `seraphine-needs-requirements` after deep research).

---

## ⚠️ Error States & Edge Cases

### No Viable Options Found
* **Condition:** After exhaustive searching, no clear technical solution or relevant information is found.
* **Action:** 
  1. Comment on the issue detailing the exact search queries used, the areas explored, and a clear statement that no viable options were found.
  2. Ask clarifying questions to the user or maintainer to unblock the research.
  3. Remove the `seraphine-needs-deep-research` label.
  4. Do **not** transition to the requirements phase.

### Ambiguous Problem Space
* **Condition:** The issue description or initial research yields multiple, highly divergent paths, and the exact problem to solve is too poorly defined to make a recommendation.
* **Action:**
  1. Comment on the issue outlining the ambiguity, the different potential paths discovered, and why they conflict or require clarification.
  2. Request specific input from the stakeholders to narrow down the scope.
  3. Remove the `seraphine-needs-deep-research` label.
  4. Do **not** transition to the requirements phase until clarity is provided.
