# 🔍 The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the AI assistant is triggered to execute an extensive research and discovery phase.

## 🔄 Workflow Lifecycle

1. **Search:** Conduct comprehensive research using available tools (web search, codebase scanning, etc.) to gather context and possible solutions.
2. **Synthesize:** Analyze and synthesize the gathered information into a coherent summary, evaluating potential approaches, trade-offs, and design decisions.
3. **Comment:** Post the synthesized findings and proposed path forward as a comment on the issue.
4. **Remove Label:** Remove the `seraphine-needs-deep-research` label from the issue once the research is adequately documented.
5. **Create Sub-issue:** Create a new sub-issue or transition the current issue to the next appropriate state (e.g., `seraphine-needs-implementation-plan` or `seraphine-ready-to-implement`) based on the findings.

---

## ⚠️ Error States & Edge Cases

If the agent encounters issues during the research phase, it must document them clearly and halt progression:

### 1. No Viable Options Found
If the research does not yield any viable solutions, libraries, or architectural paths to achieve the goal:
* **Action:** Document the failed search vectors and the reasons why options are unviable.
* **Next Steps:** Leave a comment explaining the dead-end and ask for clarification or a pivot from the user/maintainers. Do **not** remove the label or progress the state.

### 2. Ambiguous Problem Space
If the initial issue description or the discovered research is too ambiguous or broad to form a concrete synthesis:
* **Action:** Document the specific areas of ambiguity and what questions need answering.
* **Next Steps:** Leave a comment detailing the missing constraints or confusing requirements. Wait for user input. Do **not** proceed to synthesis or state transition.

## 🚫 Completion
As per the critical general rules, the agent must stop once the label is removed or if an error state is reached requiring user feedback. Do not proceed to process the newly created sub-issue or next label in the same run.
