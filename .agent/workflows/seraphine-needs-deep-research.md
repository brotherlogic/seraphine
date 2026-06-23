# 🕵️ The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the agent must execute a comprehensive research phase to understand the problem space before any requirements can be gathered.

## 🔄 Workflow Sequence

1. **Search**: Programmatically gather information from the codebase, existing documentation, and external resources (e.g., web searches) that are relevant to the issue.
2. **Synthesize**: Analyze the gathered information and formulate a comprehensive summary of the problem, the context, and potential paths forward.
3. **Comment**: Post the synthesized research findings as a detailed comment on the current issue.
4. **Remove Label**: Remove the `seraphine-needs-deep-research` label from the current issue to mark this phase as complete.
5. **Create Sub-issue**: Create a new child issue and label it with `seraphine-needs-requirements` to initiate the requirements gathering phase.

## 🚨 Error States & Handling

* **No Viable Options Found**:
  * If the research yields no practical solutions or feasible paths forward, document this finding clearly in the issue comment, explaining why no viable options exist.
  * Remove the `seraphine-needs-deep-research` label.
  * Tag the issue author requesting guidance or alternative directions, and do NOT create the next sub-issue.

* **Ambiguous Problem Space**:
  * If the problem remains too broad or lacks sufficient context to synthesize a concrete summary, detail the specific missing information in the issue comment.
  * Remove the `seraphine-needs-deep-research` label.
  * Tag the issue author to provide the necessary context, and do NOT create the next sub-issue.
