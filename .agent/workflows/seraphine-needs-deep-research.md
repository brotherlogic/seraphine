# 🔬 The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the agent performs a deep research phase to identify options before requirements gathering.

## 🔄 Workflow Lifecycle

1. Search
2. Synthesize
3. Comment
4. Remove Label
5. Create Sub-issue

---

## 📋 Phase Guidelines

### 1. Search
* **Action:** Activate the `deep-research` skill to begin evaluating the problem space.
* **Process:** Research the domain, read necessary documentation, and evaluate the problem space until at least the top 3 viable options are identified.

### 2. Synthesize
* **Action:** Synthesize an organic summary of the research.
* **Process:** Decide organically how best to present the pros and cons of the identified options. There is no strict formatting template required, but ensure the findings are clear and actionable.

### 3. Comment
* **Action:** Post the summary of your findings to the parent issue.
* **Process:** Add a comment to the GitHub issue detailing the research, the top options evaluated, and any recommendations.

### 4. Remove Label
* **Action:** Remove the `seraphine-needs-deep-research` label from the issue.

### 5. Create Sub-issue
* **Action:** Transition to the requirements gathering phase.
* **Process:** Create a new sub-issue labeled with `seraphine-needs-requirements` to initiate requirements gathering for the chosen path.

---

## ⚠️ Error States & Handling

### "No Viable Options Found"
If no viable options can be identified:
1. **Comment:** Document the exhaustion of research avenues on the issue, explaining why the options explored were not viable.
2. **Halt Process:** Do NOT remove the label or create a requirements sub-issue.
3. **Escalate:** Ask the user for clarification or additional context to unblock the research.

### "Ambiguous Problem Space"
If the problem space is too broad or ambiguous:
1. **Comment:** Post a comment detailing the ambiguity and what specific information is needed to narrow the scope.
2. **Halt Process:** Do NOT remove the label or create a requirements sub-issue.
3. **Escalate:** Wait for the user or stakeholder to provide the necessary clarification before resuming the research.
